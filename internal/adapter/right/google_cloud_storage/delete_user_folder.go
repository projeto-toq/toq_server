package gcsadapter

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) DeleteUserFolder(ctx context.Context, userID int64) error {
	if g.adminClient == nil {
		return status.Error(codes.FailedPrecondition, "gcs admin client not initialized")
	}

	bucketHandle := g.adminClient.Bucket(UsersBucketName)
	prefix := fmt.Sprintf("%d/", userID)

	slog.Info("starting efficient user folder deletion", "userID", userID, "bucket", UsersBucketName, "prefix", prefix)

	// 1. Coleta TODOS os objetos de uma vez (incluindo subpastas)
	allObjects, err := g.collectAllObjectsEfficiently(ctx, bucketHandle, prefix, userID)
	if err != nil {
		return err
	}

	slog.Info("collected all objects for deletion", "userID", userID, "totalCount", len(allObjects))

	// 2. Deleção em paralelo otimizada
	if err := g.deleteAllObjectsInParallel(ctx, bucketHandle, allObjects, userID); err != nil {
		return err
	}

	// 3. NOVA: Deleção explícita de marcadores de pasta
	if err := g.deleteExplicitFolderMarkers(ctx, bucketHandle, prefix, userID); err != nil {
		return err
	}

	// 4. Verificação final simplificada
	if err := g.verifyCompleteRemoval(ctx, bucketHandle, prefix, userID); err != nil {
		return err
	}

	slog.Info("user folder completely deleted", "userID", userID, "bucket", UsersBucketName)
	return nil
}

// collectAllObjectsEfficiently coleta todos os objetos com paginação otimizada
func (g *GCSAdapter) collectAllObjectsEfficiently(ctx context.Context, bucketHandle *storage.BucketHandle, prefix string, userID int64) ([]string, error) {
	var allObjects []string

	slog.Debug("starting comprehensive object collection", "userID", userID, "prefix", prefix)

	// Query otimizada para capturar TODOS os objetos recursivamente
	query := &storage.Query{
		Prefix: prefix,
		// Garante que capture objetos em subpastas também
	}

	it := bucketHandle.Objects(ctx, query)
	objectCount := 0

	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			slog.Error("failed to iterate over objects", "userID", userID, "error", err)
			return nil, status.Error(codes.Internal, "failed to collect objects")
		}

		allObjects = append(allObjects, objAttrs.Name)
		objectCount++

		// Log detalhado de cada objeto encontrado
		slog.Debug("object collected",
			"userID", userID,
			"object", objAttrs.Name,
			"size", objAttrs.Size,
			"count", objectCount)
	}

	// Log do resultado da coleta
	slog.Debug("object collection completed",
		"userID", userID,
		"totalObjects", len(allObjects),
		"objects", allObjects) // Lista completa para debugging

	return allObjects, nil
}

// deleteAllObjectsInParallel deleta todos os objetos usando worker pool
func (g *GCSAdapter) deleteAllObjectsInParallel(ctx context.Context, bucketHandle *storage.BucketHandle, objects []string, userID int64) error {
	if len(objects) == 0 {
		slog.Info("no objects to delete", "userID", userID)
		return nil
	}

	const maxWorkers = 10
	const batchSize = 50

	// Divide objetos em lotes
	batches := g.chunkObjects(objects, batchSize)
	slog.Debug("deletion batches created", "userID", userID, "batchCount", len(batches), "batchSize", batchSize)

	// Canal para erros
	errChan := make(chan error, len(batches))

	// Semáforo para controlar workers
	semaphore := make(chan struct{}, maxWorkers)

	// Dispatcher de workers
	var wg sync.WaitGroup
	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, objectBatch []string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			slog.Debug("starting batch deletion",
				"userID", userID,
				"batchIndex", batchIndex,
				"batchSize", len(objectBatch))

			// Deleta todos os objetos do lote
			for _, objectName := range objectBatch {
				if err := g.deleteObjectWithRetryOptimized(ctx, bucketHandle, objectName, userID); err != nil {
					errChan <- fmt.Errorf("batch %d failed to delete %s: %w", batchIndex, objectName, err)
					return
				}
			}

			slog.Debug("batch deletion completed",
				"userID", userID,
				"batchIndex", batchIndex)

			errChan <- nil
		}(i, batch)
	}

	// Aguarda todos os workers terminarem
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Verifica erros
	for err := range errChan {
		if err != nil {
			slog.Error("parallel deletion failed", "userID", userID, "error", err)
			return status.Error(codes.Internal, "parallel deletion failed")
		}
	}

	slog.Debug("parallel deletion completed successfully", "userID", userID)
	return nil
}

// chunkObjects divide slice de objetos em lotes
func (g *GCSAdapter) chunkObjects(objects []string, batchSize int) [][]string {
	var batches [][]string

	for i := 0; i < len(objects); i += batchSize {
		end := i + batchSize
		if end > len(objects) {
			end = len(objects)
		}
		batches = append(batches, objects[i:end])
	}

	return batches
}

// deleteObjectWithRetryOptimized versão otimizada do delete com retry
func (g *GCSAdapter) deleteObjectWithRetryOptimized(ctx context.Context, bucketHandle *storage.BucketHandle, objectName string, userID int64) error {
	const maxRetries = 3
	const baseDelay = 100 * time.Millisecond

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if err := bucketHandle.Object(objectName).Delete(ctx); err != nil {
			if attempt == maxRetries {
				slog.Error("failed to delete object after retries",
					"userID", userID,
					"object", objectName,
					"error", err,
					"attempts", maxRetries)
				return fmt.Errorf("failed to delete %s after %d attempts: %w", objectName, maxRetries, err)
			}

			// Exponential backoff
			delay := time.Duration(attempt) * baseDelay
			slog.Debug("retrying object deletion",
				"userID", userID,
				"object", objectName,
				"attempt", attempt,
				"delay", delay,
				"error", err)

			time.Sleep(delay)
			continue
		}

		slog.Debug("object deleted successfully",
			"userID", userID,
			"object", objectName,
			"attempt", attempt)
		break
	}

	return nil
}

// verifyCompleteRemoval verificação final robusta
func (g *GCSAdapter) verifyCompleteRemoval(ctx context.Context, bucketHandle *storage.BucketHandle, prefix string, userID int64) error {
	// Delay maior para eventual consistency
	time.Sleep(1 * time.Second)

	slog.Debug("starting final verification", "userID", userID, "prefix", prefix)

	// Tenta múltiplas verificações se necessário
	const maxVerificationAttempts = 3

	for attempt := 1; attempt <= maxVerificationAttempts; attempt++ {
		it := bucketHandle.Objects(ctx, &storage.Query{Prefix: prefix})

		remainingObjects := []string{}
		for {
			objAttrs, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				slog.Warn("verification error", "userID", userID, "attempt", attempt, "error", err)
				if attempt == maxVerificationAttempts {
					return status.Error(codes.Internal, "verification failed due to iteration error")
				}
				time.Sleep(500 * time.Millisecond)
				break // Retry outer loop
			}

			remainingObjects = append(remainingObjects, objAttrs.Name)
		}

		if len(remainingObjects) == 0 {
			slog.Debug("verification passed - folder completely deleted",
				"userID", userID,
				"attempt", attempt)
			return nil
		}

		// Objetos ainda existem
		if attempt == maxVerificationAttempts {
			slog.Error("verification failed - objects still exist",
				"userID", userID,
				"remainingCount", len(remainingObjects),
				"remainingObjects", remainingObjects)
			return status.Error(codes.Internal, "folder deletion incomplete - objects still exist")
		}

		slog.Debug("verification found remaining objects, retrying",
			"userID", userID,
			"attempt", attempt,
			"remainingCount", len(remainingObjects))

		time.Sleep(1 * time.Second) // Wait before retry
	}

	return status.Error(codes.Internal, "verification failed after maximum attempts")
}

// deleteExplicitFolderMarkers deleta explicitamente marcadores de pasta conhecidos
func (g *GCSAdapter) deleteExplicitFolderMarkers(ctx context.Context, bucketHandle *storage.BucketHandle, prefix string, userID int64) error {
	slog.Debug("starting explicit folder marker deletion", "userID", userID, "prefix", prefix)

	// Lista de marcadores de pasta conhecidos que precisam ser deletados
	potentialFolders := []string{
		fmt.Sprintf("%d/", userID),            // Pasta raiz do usuário: "39/"
		fmt.Sprintf("%d/thumbnails/", userID), // Subpasta thumbnails: "39/thumbnails/"
		fmt.Sprintf("%d/documents/", userID),  // Subpasta documents: "39/documents/"
		fmt.Sprintf("%d/images/", userID),     // Subpasta images: "39/images/"
		fmt.Sprintf("%d/videos/", userID),     // Subpasta videos: "39/videos/"
		fmt.Sprintf("%d/uploads/", userID),    // Subpasta uploads: "39/uploads/"
	}

	deletedCount := 0
	for _, folder := range potentialFolders {
		slog.Debug("attempting to delete folder marker", "userID", userID, "folder", folder)

		obj := bucketHandle.Object(folder)

		// Tenta deletar o marcador de pasta
		if err := obj.Delete(ctx); err != nil {
			if err == storage.ErrObjectNotExist {
				slog.Debug("folder marker does not exist", "userID", userID, "folder", folder)
				continue
			}
			slog.Warn("failed to delete folder marker", "userID", userID, "folder", folder, "error", err)
			// Não retorna erro - continua tentando outras pastas
			continue
		}

		deletedCount++
		slog.Debug("folder marker deleted successfully", "userID", userID, "folder", folder)
	}

	slog.Debug("explicit folder marker deletion completed",
		"userID", userID,
		"totalAttempted", len(potentialFolders),
		"deletedCount", deletedCount)

	return nil
}
