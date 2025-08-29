package s3adapter

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/giulio-alfieri/toq_server/internal/core/utils"
)

func (s *S3Adapter) CreateUserFolder(ctx context.Context, UserID int64) (err error) {
	if s.adminClient == nil {
		err = utils.ErrInternalServer
		return
	}

	slog.Debug("Creating user folder structure in S3", "userID", UserID, "bucket", s.bucketName)

	// Lista de placeholders para criar toda a estrutura de pastas
	placeholders := []string{
		fmt.Sprintf("%d/.placeholder", UserID),            // Pasta raiz do usuário
		fmt.Sprintf("%d/thumbnails/.placeholder", UserID), // Pasta thumbnails
	}

	// Criar cada placeholder para garantir que as "pastas" existam
	for _, placeholderPath := range placeholders {
		_, err := s.adminClient.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s.bucketName),
			Key:    aws.String(placeholderPath),
			Body:   nil, // Objeto vazio
		})

		if err != nil {
			slog.Error("failed to create placeholder in S3", "userID", UserID, "path", placeholderPath, "error", err)
			return utils.ErrInternalServer
		}

		slog.Debug("placeholder created in S3", "userID", UserID, "path", placeholderPath)
	}

	slog.Info("user folder structure created successfully in S3", "userID", UserID, "bucket", s.bucketName)
	return nil
}
