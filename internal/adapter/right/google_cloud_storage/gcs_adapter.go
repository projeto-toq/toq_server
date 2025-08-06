package gcsadapter

import (
	"context"
	"log/slog"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

type GCSAdapter struct {
	adminClient  *storage.Client
	readerClient *storage.Client
	writerClient *storage.Client
}

func NewGCSAdapter(ctx context.Context) (gcs *GCSAdapter, CloseFunc func() error) {
	adminClient, err := storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_admin.json"))
	if err != nil {
		slog.Error("failed to create admin storage client", "error", err)
		os.Exit(1)
	}

	writerClient, err := storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_writer.json"))
	if err != nil {
		slog.Error("failed to create writer storage client", "error", err)
		os.Exit(1)
	}

	readerClient, err := storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_reader.json"))
	if err != nil {
		slog.Error("failed to create reader storage client", "error", err)
		os.Exit(1)
	}

	gcs = &GCSAdapter{
		adminClient:  adminClient,
		readerClient: readerClient,
		writerClient: writerClient,
	}

	CloseFunc = func() error {
		if gcs.adminClient != nil {
			if err := gcs.adminClient.Close(); err != nil {
				return err
			}
		}
		if gcs.readerClient != nil {
			if err := gcs.readerClient.Close(); err != nil {
				return err
			}
		}
		if gcs.writerClient != nil {
			if err := gcs.writerClient.Close(); err != nil {
				return err
			}
		}
		return nil
	}
	return
}
