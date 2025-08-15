package creciadapter

import (
	"context"
	"fmt"
	"log/slog"

	vision "cloud.google.com/go/vision/v2/apiv1"
	"google.golang.org/api/option"
)

type CreciAdapter struct {
	client      *vision.ImageAnnotatorClient
	readerCreds []byte
}

func NewCreciAdapter(ctx context.Context, readerCreds []byte) (*CreciAdapter, func(), error) {
	client, err := vision.NewImageAnnotatorClient(ctx, option.WithCredentialsJSON(readerCreds))
	if err != nil {
		slog.Error("Failed to create vision client", "error", err)
		return nil, nil, err
	}

	adapter := &CreciAdapter{
		client: client,
	}

	closeFunc := func() {
		if adapter.client != nil {
			adapter.client.Close()
		}
	}

	return adapter, closeFunc, nil
}

// As funções Open e Close agora são gerenciadas pelo ciclo de vida do factory
// e podem ser removidas se não forem mais chamadas diretamente.
// Por segurança, vamos mantê-las vazias por enquanto.

func (ca *CreciAdapter) Close() {
	// O cliente agora é fechado pela closeFunc retornada pelo NewCreciAdapter
}

func (ca *CreciAdapter) Open(ctx context.Context) (err error) {
	// O cliente agora é aberto pelo NewCreciAdapter
	if ca.client == nil {
		return fmt.Errorf("vision client not initialized")
	}
	return nil
}
