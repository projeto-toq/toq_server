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

	// Criar objeto placeholder para garantir que a "pasta" do usuário existe
	placeholderPath := fmt.Sprintf("%d/.placeholder", UserID)

	writer := bucketHandle.Object(placeholderPath).NewWriter(ctx)
	defer writer.Close()

	_, err = writer.Write([]byte(""))
	if err != nil {
		slog.Error("failed to create user folder", "userID", UserID, "error", err)
		err = status.Error(codes.Internal, "failed to create user folder")
		return
	}

	slog.Info("user folder created successfully", "userID", UserID, "bucket", UsersBucketName)
	return
}
