package webapi

import (
	"context"
	"fmt"
	"lintang/video-transcoder-api/biz/domain"
	"lintang/video-transcoder-api/config"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

type MinioAPI struct {
	BaseURL         string
	AccessKeyID     string
	SecretAccessKey string
}

func NewMinioAPI(cfg *config.Config) *MinioAPI {
	zap.L().Info("minio")
	zap.L().Info(fmt.Sprintf("acc_key minio %s", os.Getenv("ACC_KEY_MINIO")))
	return &MinioAPI{
		BaseURL:         cfg.Minio.BaseURL,
		AccessKeyID:     os.Getenv("ACC_KEY_MINIO"),
		SecretAccessKey: os.Getenv("SECRET_KEY_MINIO"),
	}
}

func (m *MinioAPI) CreatePresignedURLForUpload(ctx context.Context, filename string) (string, error) {
	endpoint := m.BaseURL
	accessKeyID := m.AccessKeyID
	secretAccessKey := m.SecretAccessKey
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		zap.L().Error("new minio", zap.Error(err))
		return "", domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	bucketName := "oti-be-bucket"

	presignedURL, err := minioClient.PresignedPutObject(ctx, bucketName, fmt.Sprintf("%s/%s.mp4", filename, filename), 48*time.Hour)
	if err != nil {
		zap.L().Error("minioClient.PresignedPutObject (CreatePresignedURLForUpload)", zap.Error(err))
		return "", err
	}
	return presignedURL.String(), nil
}
