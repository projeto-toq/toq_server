package s3adapter

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *S3Adapter) DeleteUserFolder(ctx context.Context, userID int64) error {
	if s.adminClient == nil {
		return status.Error(codes.FailedPrecondition, "s3 admin client not initialized")
	}

	prefix := fmt.Sprintf("%d/", userID)
	slog.Info("starting efficient user folder deletion in S3", "userID", userID, "bucket", s.bucketName, "prefix", prefix)

	// 1. Listar todos os objetos do usuário
	allObjects, err := s.listAllObjectsWithPrefix(ctx, prefix, userID)
	if err != nil {
		return err
	}

	slog.Info("collected all objects for deletion", "userID", userID, "totalCount", len(allObjects))

	// 2. Deletar em lotes (S3 permite até 1000 objetos por batch)
	if err := s.deleteObjectsInBatches(ctx, allObjects, userID); err != nil {
		return err
	}

	slog.Info("user folder completely deleted from S3", "userID", userID, "bucket", s.bucketName)
	return nil
}

// listAllObjectsWithPrefix lista todos os objetos com o prefixo especificado
func (s *S3Adapter) listAllObjectsWithPrefix(ctx context.Context, prefix string, userID int64) ([]string, error) {
	var allObjects []string

	slog.Debug("starting comprehensive object collection in S3", "userID", userID, "prefix", prefix)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(s.adminClient, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			slog.Error("failed to list objects in S3", "userID", userID, "error", err)
			return nil, status.Error(codes.Internal, "failed to list objects")
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				allObjects = append(allObjects, *obj.Key)
				slog.Debug("object collected", "userID", userID, "object", *obj.Key, "size", obj.Size)
			}
		}
	}

	slog.Debug("object collection completed", "userID", userID, "totalObjects", len(allObjects))
	return allObjects, nil
}

// deleteObjectsInBatches deleta objetos em lotes de até 1000 (limite do S3)
func (s *S3Adapter) deleteObjectsInBatches(ctx context.Context, objects []string, userID int64) error {
	if len(objects) == 0 {
		slog.Info("no objects to delete", "userID", userID)
		return nil
	}

	const batchSize = 1000 // Limite máximo do S3 para delete em lote
	const maxWorkers = 5

	// Dividir em lotes
	batches := s.chunkObjects(objects, batchSize)
	slog.Debug("deletion batches created", "userID", userID, "batchCount", len(batches), "batchSize", batchSize)

	// Canal para erros
	errChan := make(chan error, len(batches))

	// Semáforo para controlar workers
	semaphore := make(chan struct{}, maxWorkers)

	// Worker pool para processamento paralelo
	var wg sync.WaitGroup
	for i, batch := range batches {
		wg.Add(1)
		go func(batchIndex int, objectBatch []string) {
			defer wg.Done()

			// Acquire semaphore
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			slog.Debug("starting batch deletion", "userID", userID, "batchIndex", batchIndex, "batchSize", len(objectBatch))

			if err := s.deleteBatch(ctx, objectBatch, userID, batchIndex); err != nil {
				errChan <- fmt.Errorf("batch %d failed: %w", batchIndex, err)
				return
			}

			slog.Debug("batch deletion completed", "userID", userID, "batchIndex", batchIndex)
			errChan <- nil
		}(i, batch)
	}

	// Aguardar todos os workers
	go func() {
		wg.Wait()
		close(errChan)
	}()

	// Verificar erros
	for err := range errChan {
		if err != nil {
			slog.Error("parallel deletion failed", "userID", userID, "error", err)
			return status.Error(codes.Internal, "parallel deletion failed")
		}
	}

	slog.Debug("parallel deletion completed successfully", "userID", userID)
	return nil
}

// deleteBatch deleta um lote de objetos usando S3 DeleteObjects API
func (s *S3Adapter) deleteBatch(ctx context.Context, objects []string, userID int64, batchIndex int) error {
	if len(objects) == 0 {
		return nil
	}

	// Converter para formato S3
	var objectIdentifiers []types.ObjectIdentifier
	for _, objKey := range objects {
		objectIdentifiers = append(objectIdentifiers, types.ObjectIdentifier{
			Key: aws.String(objKey),
		})
	}

	// Executar delete em lote
	input := &s3.DeleteObjectsInput{
		Bucket: aws.String(s.bucketName),
		Delete: &types.Delete{
			Objects: objectIdentifiers,
			Quiet:   aws.Bool(true), // Não retornar objetos deletados com sucesso
		},
	}

	output, err := s.adminClient.DeleteObjects(ctx, input)
	if err != nil {
		slog.Error("failed to delete batch in S3", "userID", userID, "batchIndex", batchIndex, "error", err)
		return fmt.Errorf("failed to delete batch: %w", err)
	}

	// Verificar se houve erros em objetos específicos
	if len(output.Errors) > 0 {
		for _, deleteError := range output.Errors {
			slog.Warn("failed to delete specific object", "userID", userID, "object", *deleteError.Key, "error", *deleteError.Message)
		}
		return fmt.Errorf("some objects failed to delete in batch %d", batchIndex)
	}

	slog.Debug("batch deleted successfully", "userID", userID, "batchIndex", batchIndex, "objectCount", len(objects))
	return nil
}

// chunkObjects divide slice de objetos em lotes
func (s *S3Adapter) chunkObjects(objects []string, batchSize int) [][]string {
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
