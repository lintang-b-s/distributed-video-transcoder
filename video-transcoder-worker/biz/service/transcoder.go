package service

import (
	"context"
	"fmt"
	"io"
	"lintang/video-processing-worker/biz/domain"
	"lintang/video-processing-worker/biz/util"
	"log"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"go.uber.org/zap"
)

type DkronCLIAPI interface {
	AddJobUploadPlaylistToMinio(ctx context.Context, filename string) error
}
type MinioAPI interface {
	GetUserVideoMP4URL(filename string) (*minio.Object, error)
	GetAllBitrateVideoVersion(objectFolderName string) ([]*minio.Object, error)
	BitrateVersionVideoUploader(objectFolderName string, fileName string, filePath string) error
	HlsPlaylistUploader(objectFolderName string, folderSource string) error
}
type TranscoderService struct {
	dkronAPI DkronCLIAPI
	minioAPI MinioAPI
}

func NewTranscoderService(d DkronCLIAPI, m MinioAPI) *TranscoderService {
	return &TranscoderService{
		d, m,
	}
}

func (s *TranscoderService) Transcode(ctx context.Context, filename string, resolution string) error {


	// pull video mp4 yang baru diupload user (kalau belum ada di minioo, ya pull lagi setiap 300 miliseconds)
	videoFile, err := s.minioAPI.GetUserVideoMP4URL(filename)
	if err != nil {
		zap.L().Error("createHLSFromMinioObject", zap.Error(err))
	}

	// membuat folder dengan nama filename
	os.Mkdir(filename, 0777)
	filePath := fmt.Sprintf(`%s/%s.mp4`, filename, filename)
	mylocalFile, err := os.Create(filePath) // membuat file filename/filename.mp4 ke lokal dari minio object
	if err != nil {
		zap.L().Error("os.Create (Transcode) (TranscoderService)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	defer mylocalFile.Close()
	var stat minio.ObjectInfo
	stat, err = videoFile.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			fileInfo, err := s.minioAPI.GetUserVideoMP4URL(filename)
			if stat, err = fileInfo.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
		zap.L().Debug("pull video mp4 user lagi...  GetUserVideoMP4URL (Transcode) (TranscoderService)", zap.String("filename", filename))
	}

	if _, err = io.CopyN(mylocalFile, videoFile, stat.Size); err != nil {
		// copy bytes dari minio object ke local file yang baru dibuat
		zap.L().Error("io.CopyN(mylocalFile, videoFile, stat.Size) (Transcode) (TranscoderService)", zap.String("filename", filename), zap.Error(err))
		return err 
	}

	if resolution == "240p" {
		
	}
	util.CreateHLSFromMinioObject()

	return nil
}
