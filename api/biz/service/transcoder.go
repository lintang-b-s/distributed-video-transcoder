package service

import (
	"context"
	"lintang/video-transcoder-api/biz/domain"
	"strings"
)

type DkronAPI interface {
	AddJobUploadPlaylistToMinio(ctx context.Context, filename string) error
}

type MinioAPI interface {
	CreatePresignedURLForUpload(ctx context.Context, filename string) (string, error)
}

type MetadataRepo interface {
	Insert(ctx context.Context, m domain.VideoMetadata) error
	GetAll(ctx context.Context) ([]domain.VideoMetadata, error)
}

type TranscoderService struct {
	dkronAPI DkronAPI
	minioAPI MinioAPI
	mongo    MetadataRepo
}

func NewTranscoderService(d DkronAPI, m MinioAPI, mongo MetadataRepo) *TranscoderService {
	return &TranscoderService{
		d, m, mongo,
	}
}

func (s *TranscoderService) CreatePresignedURLForUpload(ctx context.Context, filename string) (string, error) {
	err := s.dkronAPI.AddJobUploadPlaylistToMinio(ctx, filename)
	if err != nil {
		return "", err
	}
	presignedURL, err := s.minioAPI.CreatePresignedURLForUpload(ctx, filename)
	if err != nil {
		return "", err
	}
	presignedURL = strings.Replace(presignedURL, "minio:9000", "localhost:9091", 1)

	return presignedURL, nil
}

func (s *TranscoderService) GetAllVideosMetadata(ctx context.Context) ([]domain.VideoMetadata, error) {
	metadatas, err := s.mongo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return metadatas, nil
}
