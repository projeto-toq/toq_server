package s3adapter

import (
	"context"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	globalmodel "github.com/giulio-alfieri/toq_server/internal/core/model/global_model"
)

type S3Adapter struct {
	adminClient  *s3.Client
	readerClient *s3.Client
	uploader     *manager.Uploader
	downloader   *manager.Downloader
	bucketName   string
	region       string
}

func NewS3Adapter(ctx context.Context, env *globalmodel.Environment) (s3Adapter *S3Adapter, CloseFunc func() error, err error) {
	slog.Info("Creating S3 adapter", "region", env.S3.Region, "bucket", env.S3.BucketName)

	// Configuração básica da AWS
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(env.S3.Region),
	)
	if err != nil {
		slog.Error("failed to load AWS config", "error", err)
		return nil, nil, err
	}

	// Cliente Admin com credenciais específicas
	adminCfg := cfg.Copy()
	if env.S3.AdminRole.AccessKeyID != "" && env.S3.AdminRole.SecretAccessKey != "" {
		adminCfg.Credentials = credentials.NewStaticCredentialsProvider(
			env.S3.AdminRole.AccessKeyID,
			env.S3.AdminRole.SecretAccessKey,
			"", // session token
		)
	}
	adminClient := s3.NewFromConfig(adminCfg)

	// Cliente Reader com credenciais específicas
	readerCfg := cfg.Copy()
	if env.S3.ReaderRole.AccessKeyID != "" && env.S3.ReaderRole.SecretAccessKey != "" {
		readerCfg.Credentials = credentials.NewStaticCredentialsProvider(
			env.S3.ReaderRole.AccessKeyID,
			env.S3.ReaderRole.SecretAccessKey,
			"", // session token
		)
	}
	readerClient := s3.NewFromConfig(readerCfg)

	// Uploader e Downloader usando admin/reader clients
	uploader := manager.NewUploader(adminClient)
	downloader := manager.NewDownloader(readerClient)

	s3Adapter = &S3Adapter{
		adminClient:  adminClient,
		readerClient: readerClient,
		uploader:     uploader,
		downloader:   downloader,
		bucketName:   env.S3.BucketName,
		region:       env.S3.Region,
	}

	// CloseFunc (S3 clients não precisam de Close explícito, mas mantemos para compatibilidade)
	CloseFunc = func() error {
		slog.Debug("S3 adapter cleanup completed")
		return nil
	}

	slog.Info("S3 adapter created successfully", "bucket", env.S3.BucketName, "region", env.S3.Region)
	return s3Adapter, CloseFunc, nil
}
