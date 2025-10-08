package s3adapter

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (s *S3Adapter) CreateUserFolder(ctx context.Context, UserID int64) (err error) {
	if s.adminClient == nil {
		err = errors.New("s3 admin client is nil")
		return
	}

	ctx = utils.ContextWithLogger(ctx)
	logger := utils.LoggerFromContext(ctx)
	logger.Debug("adapter.s3.create_user_folder.start", "user_id", UserID, "bucket", s.bucketName)

	// Lista de placeholders para criar toda a estrutura de pastas
	placeholders := []string{
		fmt.Sprintf("%d/.placeholder", UserID),       // Pasta raiz do usuário
		fmt.Sprintf("%d/photo/.placeholder", UserID), // Pasta de fotos padronizada
	}

	// Criar cada placeholder para garantir que as "pastas" existam
	for _, placeholderPath := range placeholders {
		_, err := s.adminClient.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(placeholderPath),
			Body:   nil, // Objeto vazio
		})

		if err != nil {
			utils.SetSpanError(ctx, err)
			logger.Error("adapter.s3.create_user_folder.placeholder_error", "user_id", UserID, "path", placeholderPath, "error", err)
			return err
		}

		logger.Debug("adapter.s3.create_user_folder.placeholder_created", "user_id", UserID, "path", placeholderPath)
	}

	logger.Info("adapter.s3.create_user_folder.success", "user_id", UserID, "bucket", s.bucketName)
	return nil
}
