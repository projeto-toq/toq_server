package creciadapter

import (
	"context"
	"log/slog"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"google.golang.org/api/option"
)

type CreciAdapter struct {
	client *vision.ImageAnnotatorClient
}

func NewCreciAdapter(ctx context.Context) *CreciAdapter {
	return &CreciAdapter{}
}

func (ca *CreciAdapter) Close() {
	ca.client.Close()
}

func (ca *CreciAdapter) Open(ctx context.Context) (err error) {
	ca.client, err = vision.NewImageAnnotatorClient(ctx,
		option.WithCredentialsFile("../configs/gcs_writer.json"))
	if err != nil {
		slog.Error("Failed to create vision client: ", "error:", err.Error())
		return
	}
	return
}
