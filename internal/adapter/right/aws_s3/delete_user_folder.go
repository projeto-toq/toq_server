package s3adapter

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/projeto-toq/toq_server/internal/core/utils"
)

func (s *S3Adapter) DeleteUserFolder(ctx context.Context, userID int64) error {
	if s.adminClient == nil {
		return errors.New("s3 admin client is nil")
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	prefix := fmt.Sprintf("%d/", userID)
	logger.Info("adapter.s3.delete_user_folder.start", "user_id", userID, "bucket", s.bucketName, "prefix", prefix)

	// 1. Listar todos os objetos do usuário
	allObjects, err := s.listAllObjectsWithPrefix(ctx, prefix, userID)
	if err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}

	logger.Info("adapter.s3.delete_user_folder.collected", "user_id", userID, "total_count", len(allObjects))

	// 2. Deletar em lotes (S3 permite até 1000 objetos por batch)
	if err := s.deleteObjectsInBatches(ctx, allObjects, userID); err != nil {
		utils.SetSpanError(ctx, err)
		return err
	}

	logger.Info("adapter.s3.delete_user_folder.success", "user_id", userID, "bucket", s.bucketName)
	return nil
}

// listAllObjectsWithPrefix lista todos os objetos com o prefixo especificado
func (s *S3Adapter) listAllObjectsWithPrefix(ctx context.Context, prefix string, userID int64) ([]string, error) {
	var allObjects []string

	logger := utils.LoggerFromContext(ctx)
	logger.Debug("adapter.s3.delete_user_folder.list.start", "user_id", userID, "prefix", prefix)

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucketName),
		Prefix: aws.String(prefix),
	}

	paginator := s3.NewListObjectsV2Paginator(s.adminClient, input)

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("adapter.s3.delete_user_folder.list.error", "user_id", userID, "error", err)
			return nil, err
		}

		for _, obj := range output.Contents {
			if obj.Key != nil {
				allObjects = append(allObjects, *obj.Key)
				logger.Debug("adapter.s3.delete_user_folder.object_collected", "user_id", userID, "object", *obj.Key, "size", obj.Size)
			}
		}
	}

	logger.Debug("adapter.s3.delete_user_folder.list.completed", "user_id", userID, "total_objects", len(allObjects))
	return allObjects, nil
}

// deleteObjectsInBatches deleta objetos em lotes de até 1000 (limite do S3)
func (s *S3Adapter) deleteObjectsInBatches(ctx context.Context, objects []string, userID int64) error {
	if len(objects) == 0 {
		logger := utils.LoggerFromContext(ctx)
		logger.Info("adapter.s3.delete_user_folder.no_objects", "user_id", userID)
		return nil
	}

	logger := utils.LoggerFromContext(ctx)
	const batchSize = 1000 // Limite máximo do S3 para delete em lote
	const maxWorkers = 5

	// Dividir em lotes
	batches := s.chunkObjects(objects, batchSize)
	logger.Debug("adapter.s3.delete_user_folder.batches_created", "user_id", userID, "batch_count", len(batches), "batch_size", batchSize)

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

			logger.Debug("adapter.s3.delete_user_folder.batch_start", "user_id", userID, "batch_index", batchIndex, "batch_size", len(objectBatch))

			if err := s.deleteBatch(ctx, objectBatch, userID, batchIndex); err != nil {
				errChan <- fmt.Errorf("batch %d failed: %w", batchIndex, err)
				return
			}

			logger.Debug("adapter.s3.delete_user_folder.batch_completed", "user_id", userID, "batch_index", batchIndex)
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
			logger.Error("adapter.s3.delete_user_folder.parallel_error", "user_id", userID, "error", err)
			return err
		}
	}

	logger.Debug("adapter.s3.delete_user_folder.parallel_success", "user_id", userID)
	return nil
}

// deleteBatch deleta um lote de objetos usando S3 DeleteObjects API
func (s *S3Adapter) deleteBatch(ctx context.Context, objects []string, userID int64, batchIndex int) error {
	if len(objects) == 0 {
		return nil
	}

	logger := utils.LoggerFromContext(ctx)
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
		utils.SetSpanError(ctx, err)
		logger.Error("adapter.s3.delete_user_folder.batch_error", "user_id", userID, "batch_index", batchIndex, "error", err)
		return fmt.Errorf("failed to delete batch: %w", err)
	}

	// Verificar se houve erros em objetos específicos
	if len(output.Errors) > 0 {
		for _, deleteError := range output.Errors {
			logger.Warn("adapter.s3.delete_user_folder.object_delete_error", "user_id", userID, "object", *deleteError.Key, "error", *deleteError.Message)
		}
		return fmt.Errorf("some objects failed to delete in batch %d", batchIndex)
	}

	logger.Debug("adapter.s3.delete_user_folder.batch_success", "user_id", userID, "batch_index", batchIndex, "object_count", len(objects))
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
