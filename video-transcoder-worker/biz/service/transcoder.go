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
	"strings"
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
	UploadThumbnail(objectFolderName string, fileName string, filePath string) error
	GetThumbnail(filename string) (string, error)
	GetTranscodedVideoURL(filename string) (string, error)
}

type MetadataMQ interface {
	PublishNewMetadata(ctx context.Context, d domain.VideoMetadataMessage) error
}

type TranscoderService struct {
	dkronAPI   DkronCLIAPI
	minioAPI   MinioAPI
	metadataMQ MetadataMQ
}

func NewTranscoderService(d DkronCLIAPI, m MinioAPI, mtMQ MetadataMQ) *TranscoderService {
	return &TranscoderService{
		d, m, mtMQ,
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
			videoFile, _ = s.minioAPI.GetUserVideoMP4URL(filename)
			zap.L().Debug("pull video user ....", zap.String("filename", fmt.Sprintf("%s", filename)))
			if stat, err = videoFile.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}
	zap.L().Info(fmt.Sprintf("berhasil pull video %s.mp4", filename))
	if _, err = io.CopyN(mylocalFile, videoFile, stat.Size); err != nil {
		// copy bytes dari minio object ke local file yang baru dibuat
		zap.L().Error("io.CopyN(mylocalFile, videoFile, stat.Size) (Transcode) (TranscoderService)", zap.String("filename", filename), zap.Error(err))
		return err
	}

	// bikin multi-bit-rate version video & segmented video, lalu upload ke minio semua segmented mp4nya
	switch resolution {
	case router.Res240p:
		err = util.CreateBitrate240pVideo(filePath, filename)
		if err != nil {
			zap.L().Error("util.CreateBitrate240pVideo (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "240-f.mp4", fmt.Sprintf("./%s/240-f.mp4", filename))

		break
	case router.Res360p:
		err = util.CreateBitrate360pVideo(filePath, filename)
		if err != nil {
			zap.L().Error("util.CreateBitrateRes360pVideo (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)

		}
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "360-f.mp4", fmt.Sprintf("./%s/360-f.mp4", filename))

	case router.Res480p:
		err = util.CreateBitrate480pVideo(filePath, filename)
		if err != nil {
			zap.L().Error("util.CreateBitrate480ppVideo (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)

		}
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "480-f.mp4", fmt.Sprintf("./%s/480-f.mp4", filename))

	case router.Res720p:
		err = util.CreateBitrate720pVideo(filePath, filename)
		if err != nil {
			zap.L().Error("util.CreateBitrate720pVideo (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "720-f.mp4", fmt.Sprintf("./%s/720-f.mp4", filename))
	case router.Res1080p:
		err = util.CreateBitrate1080pVideo(filePath, filename)
		if err != nil {
			zap.L().Error("util.CreateBitrate1080pVideo (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}
		err = s.minioAPI.BitrateVersionVideoUploader(fmt.Sprintf("/%s/", filename), "1080-f.mp4", fmt.Sprintf("./%s/1080-f.mp4", filename))
		if err != nil {
			zap.L().Error("util.BitrateVersionVideoUploader Res1080p (Transcode)", zap.Error(err))
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}
		// bikin thumbnailURl
		err = util.CreateVideoThumbnaill(filePath, filename)
		if err != nil {
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}
		// upload thumbnail ke minio
		err = s.minioAPI.UploadThumbnail(fmt.Sprintf("/%s/", filename), "thumbnail.png", fmt.Sprintf("./%s/thumbnail.png", filename))
		if err != nil {
			return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
		}

		// buat cron job dkron buat bikin playlist dassh nya
		err = s.dkronAPI.AddJobUploadPlaylistToMinio(ctx, filename)
		if err != nil {
			zap.L().Error(" s.dkronAPI.AddJobUploadPlaylistToMinio(ctx, filename) (Transcode)", zap.Error(err))
		}
	}
	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	return nil
}

func (s *TranscoderService) GenerateDASHPlaylist(ctx context.Context, filename string) error {
	zap.L().Info(fmt.Sprintf("membuat dash playlist untuk file %s ....", filename))
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

	// get thumbnail video
	thumbnailURL, err := s.minioAPI.GetThumbnail(filename)
	if err != nil {
		return err
	}

	// get video format dash url
	transcodedVideoURL, err := s.minioAPI.GetTranscodedVideoURL(filename)
	if err != nil {
		return err
	}

	// ubah URL jadi localhost

	if os.Getenv("APP_ENV") == "k8s" {
		thumbnailURL = strings.Replace(thumbnailURL, "minio.minio-dev.svc.cluster.local:9000", fmt.Sprintf("%s:30009", os.Getenv("Minikube_IP")), 1)
		transcodedVideoURL = strings.Replace(transcodedVideoURL, "minio.minio-dev.svc.cluster.local:9000", fmt.Sprintf("%s:30009", os.Getenv("Minikube_IP")), 1)
	} else {
		thumbnailURL = strings.Replace(thumbnailURL, "minio:9000", "localhost:9091", 1)

		transcodedVideoURL = strings.Replace(transcodedVideoURL, "minio:9000", "localhost:9091", 1)
	}

	// publish ke rabbit mq biar diconsume sama api
	err = s.metadataMQ.PublishNewMetadata(ctx, domain.VideoMetadataMessage{
		VideoURL:  transcodedVideoURL,
		Thumbnail: thumbnailURL,
	})
	if err != nil {
		zap.L().Error("s.metadataMQ.PublishNewMetadata (GenerateDASHPlaylist)", zap.Error(err))
		return err
	}

	err = os.RemoveAll(filename) //remove directory <filename>
	if err != nil {
		zap.L().Error("os.Remove(filename) (GenerateDASHPlaylist)", zap.Error(err))
		return err
	}

	zap.L().Info(fmt.Sprintf("file %s selesai dibuatkan dash playlistnya dan dash playlist sudah diupload ke minio", filename))
	return nil
}
