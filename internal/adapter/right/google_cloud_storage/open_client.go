package gcsadapter

import (
	"context"
	"log/slog"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func (g *GCSAdapter) OpenAdmin(ctx context.Context) (err error) {
	g.adminClient, err = storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_admin.json"))
	if err != nil {
		slog.Error("failed to create admin storage client", "error", err)
	}
	return
}

func (g *GCSAdapter) OpenReader(ctx context.Context) (err error) {
	g.readerClient, err = storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_reader.json"))
	if err != nil {
		slog.Error("failed to create reader storage client", "error", err)
	}
	return
}

func (g *GCSAdapter) OpenWriter(ctx context.Context) (err error) {
	g.writerClient, err = storage.NewClient(ctx, option.WithCredentialsFile("../configs/gcs_writer.json"))
	if err != nil {
		slog.Error("failed to create writer storage client", "error", err)
	}
	return
}

func (g *GCSAdapter) Close() (err error) {
	if g.adminClient != nil {
		err = g.adminClient.Close()
		if err != nil {
			slog.Error("failed to close reader storage client", "error", err)
			return
		}
	}
	if g.readerClient != nil {
		err = g.readerClient.Close()
		if err != nil {
			slog.Error("failed to close reader storage client", "error", err)
			return
		}
	}

	if g.writerClient != nil {
		err = g.writerClient.Close()
		if err != nil {
			slog.Error("failed to close writer storage client", "error", err)
			return
		}
	}
	return
}
