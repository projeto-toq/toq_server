package gcsadapter

import (
	"context"
	"fmt"
	"log/slog"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (g *GCSAdapter) CreateUserFolder(ctx context.Context, UserID int64) (err error) {
	if g.writerClient == nil {
		err = status.Error(codes.FailedPrecondition, "gcs writer client not initialized")
		return
	}

	bucketHandle := g.writerClient.Bucket(UsersBucketName)

	// Lista de placeholders para criar toda a estrutura de pastas
	placeholders := []string{
		fmt.Sprintf("%d/.placeholder", UserID),            // Pasta raiz do usuário
		fmt.Sprintf("%d/thumbnails/.placeholder", UserID), // Pasta thumbnails
	}

	// Criar cada placeholder para garantir que as "pastas" existam
	for _, placeholderPath := range placeholders {
		writer := bucketHandle.Object(placeholderPath).NewWriter(ctx)
		_, writeErr := writer.Write([]byte(""))
		closeErr := writer.Close()

		if writeErr != nil {
			slog.Error("failed to write placeholder", "userID", UserID, "path", placeholderPath, "error", writeErr)
			err = status.Error(codes.Internal, "failed to create user folder structure")
			return
		}

		if closeErr != nil {
			slog.Error("failed to close placeholder writer", "userID", UserID, "path", placeholderPath, "error", closeErr)
			err = status.Error(codes.Internal, "failed to create user folder structure")
			return
		}

		slog.Debug("placeholder created", "userID", UserID, "path", placeholderPath)
	}

	slog.Info("user folder structure created successfully", "userID", UserID, "bucket", UsersBucketName)
	return
}
