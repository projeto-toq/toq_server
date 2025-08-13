package gcsadapter

import (
	"context"
	"encoding/json"
	"log/slog"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// serviceAccount is a helper struct to unmarshal the private key from the credentials JSON.
type serviceAccount struct {
	PrivateKey string `json:"private_key"`
}

type GCSAdapter struct {
	adminClient      *storage.Client
	readerClient     *storage.Client
	writerClient     *storage.Client
	projectID        string
	adminSAEmail     string
	writerSAEmail    string
	readerSAEmail    string
	writerPrivateKey []byte
	readerPrivateKey []byte
}

func NewGCSAdapter(ctx context.Context, projectID string, adminCreds, writerCreds, readerCreds []byte, adminSAEmail, writerSAEmail, readerSAEmail string) (gcs *GCSAdapter, CloseFunc func() error, err error) {
	adminClient, err := storage.NewClient(ctx, option.WithCredentialsJSON(adminCreds))
	if err != nil {
		slog.Error("failed to create admin storage client", "error", err)
		return nil, nil, err
	}

	writerClient, err := storage.NewClient(ctx, option.WithCredentialsJSON(writerCreds))
	if err != nil {
		slog.Error("failed to create writer storage client", "error", err)
		return nil, nil, err
	}

	readerClient, err := storage.NewClient(ctx, option.WithCredentialsJSON(readerCreds))
	if err != nil {
		slog.Error("failed to create reader storage client", "error", err)
		return nil, nil, err
	}

	// Parse credentials to extract private keys for signing URLs
	var writerSA, readerSA serviceAccount
	if err := json.Unmarshal(writerCreds, &writerSA); err != nil {
		slog.Error("failed to parse writer service account credentials", "error", err)
		return nil, nil, err
	}
	if err := json.Unmarshal(readerCreds, &readerSA); err != nil {
		slog.Error("failed to parse reader service account credentials", "error", err)
		return nil, nil, err
	}

	gcs = &GCSAdapter{
		adminClient:      adminClient,
		readerClient:     readerClient,
		writerClient:     writerClient,
		projectID:        projectID,
		adminSAEmail:     adminSAEmail,
		writerSAEmail:    writerSAEmail,
		readerSAEmail:    readerSAEmail,
		writerPrivateKey: []byte(writerSA.PrivateKey),
		readerPrivateKey: []byte(readerSA.PrivateKey),
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
