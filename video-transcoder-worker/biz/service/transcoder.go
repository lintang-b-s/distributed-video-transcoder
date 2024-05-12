package service

import (
	"context"
	"fmt"
	"io"
	"lintang/video-processing-worker/biz/domain"
	"lintang/video-processing-worker/biz/router"
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

func (s *TranscoderService) Transcode(ctx context.Context, filename string, resolution router.Resolution) error {

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

	// bikin multi-bit-rate version video & segmented video, lalu upload ke minio semua segmented mp4nya
	switch resolution {
	case router.Res240p:
		err = util.CreateBitrate240pVideo(filePath, filename)
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "240-f.mp4", fmt.Sprintf("%s/240-f.mp4", filename))
		break
	case router.Res360p:
		err = util.CreateBitrate360pVideo(filePath, filename)
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "360-f.mp4", fmt.Sprintf("%s/360-f.mp4", filename))
	case router.Res480p:
		err = util.CreateBitrate480pVideo(filePath, filename)
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "480-f.mp4", fmt.Sprintf("%s/480-f.mp4", filename))
	case router.Res720p:
		err = util.CreateBitrate720pVideo(filePath, filename)
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "720-f.mp4", fmt.Sprintf("%s/720-f.mp4", filename))
	case router.Res1080p:
		err = util.CreateBitrate1080pVideo(filePath, filename)
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "1080-f.mp4", fmt.Sprintf("%s/1080-f.mp4", filename))

	}
	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	return nil
}

func (s *TranscoderService) GenerateDASHPlaylist(ctx context.Context, filename string) error {
	allBitrateVideo, err := s.minioAPI.GetAllBitrateVideoVersion(fmt.Sprintf("/%s", filename))
	if err != nil {
		return err
	}

	// pull segmented mp4 dari minio
	files := []string{"240-f.mp4", "360-f.mp4", "480-f.mp4", "720-f.mp4", "1080-f.mp4"}

	for i, _ := range files {
		os.Mkdir(filename+"/minio", 0777)
		mylocalFile, err := os.Create(filename + "/minio" + "/" + files[i])
		if err != nil {
			zap.L().Error("Os Create (GenerateDASHPlaylist)", zap.Error(err))
			return err
		}
		defer mylocalFile.Close()

		stat, err := allBitrateVideo[i].Stat()
		if err != nil {
			zap.L().Error("allBitrateVideo[i].Stat() (GenerateDASHPlaylist)", zap.Error(err))
			return err
		}

		if _, err = io.CopyN(mylocalFile, allBitrateVideo[i], stat.Size); err != nil {
			log.Fatalln(err)
		}
		allBitrateVideo[i].Close()
	}
	
	err = util.GenerateDASHPlaylistBento(filename)
	if err != nil {
		return err
	}
	// upload hls playlist ke minio
	err = s.minioAPI.HlsPlaylistUploader(fmt.Sprintf("/%s/output", filename), fmt.Sprintf("./%s/output", filename))
	if err != nil {
		return err 
	}

	err = os.Remove(filename) //remove directory <filename>
	if err != nil {
		zap.L().Error("os.Remove(filename) (GenerateDASHPlaylist)", zap.Error(err))
		return err 
	}
	return nil 
}
