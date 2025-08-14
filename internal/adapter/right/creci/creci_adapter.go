package creciadapter

import (
	"context"
	"log/slog"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"google.golang.org/api/option"
)

type CreciAdapter struct {
	client      *vision.ImageAnnotatorClient
	readerCreds []byte
}

func NewCreciAdapter(ctx context.Context, readerCreds []byte) *CreciAdapter {
	return &CreciAdapter{
		readerCreds: readerCreds,
	}
}

func (ca *CreciAdapter) Close() {
	if ca.client != nil {
		ca.client.Close()
	}
}

func (ca *CreciAdapter) Open(ctx context.Context) (err error) {
	ca.client, err = vision.NewImageAnnotatorClient(ctx,
		option.WithCredentialsJSON(ca.readerCreds))
	if err != nil {
		slog.Error("Failed to create vision client", "error", err)
		return
	}
	return
}
