package webapi

import (
	"context"
	"fmt"
	"lintang/video-processing-worker/biz/domain"
	"lintang/video-processing-worker/config"
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
	return &MinioAPI{
		BaseURL:         cfg.Minio.BaseURL,
		AccessKeyID:     cfg.Minio.AccessKeyID,
		SecretAccessKey: cfg.Minio.SecretAccessKey,
	}
}

// pulling video dari minio
func (m *MinioAPI) GetUserVideoMP4URL(filename string) (*minio.Object, error) {
	ctx := context.Background()
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
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	bucketName := "oti-be-bucket"

	userVideo, err := minioClient.GetObject(ctx, bucketName, fmt.Sprintf("%s/%s.mp4", filename, filename), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return userVideo, nil
}

func (m *MinioAPI) GetAllBitrateVideoVersion(objectFolderName string) ([]*minio.Object, error) {
	ctx := context.Background()
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
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	bucketName := "oti-be-bucket"

	// mendapatkan semua versi bitraate video
	var file240 *minio.Object
	file240, err = minioClient.GetObject(ctx, bucketName, objectFolderName+"/240-f.mp4", minio.GetObjectOptions{})
	if err != nil {
		zap.L().Error("GetObject (GetAllBitrateVideoVersion)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	_, err = file240.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			file240, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/240-f.mp4", minio.GetObjectOptions{})
			if _, err = file240.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}

	}

	file360, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/360-f.mp4", minio.GetObjectOptions{})
	if err != nil {
		zap.L().Error("GetObject (GetAllBitrateVideoVersion)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	_, err = file360.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			file360, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/360-f.mp4", minio.GetObjectOptions{})
			if _, err = file360.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}

	file480, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/480-f.mp4", minio.GetObjectOptions{})
	if err != nil {
		zap.L().Error("GetObject (GetAllBitrateVideoVersion)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	_, err = file480.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			file480, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/480-f.mp4", minio.GetObjectOptions{})
			if _, err = file480.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}

	file720, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/720-f.mp4", minio.GetObjectOptions{})
	if err != nil {
		zap.L().Error("GetObject (GetAllBitrateVideoVersion)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	_, err = file720.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			file720, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/720-f.mp4", minio.GetObjectOptions{})
			if _, err = file720.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}

	file1080, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/1080-f.mp4", minio.GetObjectOptions{})
	if err != nil {
		zap.L().Error("GetObject (GetAllBitrateVideoVersion)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	_, err = file1080.Stat()
	if err != nil {
		// kalau video belum sepenuhnya ada di minio
		for {
			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
			file1080, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/1080-f.mp4", minio.GetObjectOptions{})
			if _, err = file1080.Stat(); err == nil {
				break
			}
			time.Sleep(300 * time.Millisecond)
		}
	}

	var objects []*minio.Object
	objects = append(objects, file240, file360, file480, file720, file1080)
	return objects, nil

}

func (m *MinioAPI) BitrateVersionVideoUploader(objectFolderName string, fileName string, filePath string) error {
	ctx := context.Background()
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
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	bucketName := "oti-be-bucket"

	contentType := "application/octet-stream"

	// upload single bitrate version ke minio
	info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		zap.L().Error("minioClient.FPutObject BitrateVersionVideoUploader)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	zap.L().Info(fmt.Sprintf("Successfully uploaded %s of size %d\n", objectFolderName, info.Size))

	return nil
}

// buat upload playlist dash ke minio (object: filename/output)
func (m *MinioAPI) HlsPlaylistUploader(objectFolderName string, folderSource string) error {
	ctx := context.Background()
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
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	bucketName := "oti-be-bucket"

	contentType := "application/octet-stream"

	isiOutputDir, err := os.ReadDir(folderSource)
	if err != nil {
		zap.L().Error("os.ReadDir (folderSource) (HlsPlaylistUploader)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	// upload master dash playlist
	info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/stream.mpd", folderSource+"/stream.mpd", minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		zap.L().Error("minioClient.FPutObject (HlsPlaylistUploader)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	zap.L().Info(fmt.Sprintf("Successfully uploaded %s of size %d\n", objectFolderName, info.Size))

	// upload all video
	for _, isiOutput := range isiOutputDir {

		if isiOutput.Name() == "video" {
			isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s", folderSource, isiOutput.Name()))
			if err != nil {
				zap.L().Error("os.ReadDir isiOutputDir (HlsPlaylistUploader)", zap.Error(err))
				return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
			}
			avcDirName := isiVideoDir[0].Name()
			isiAvcDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s", folderSource, isiOutput.Name(), avcDirName))

			for _, isiAvc := range isiAvcDir {
				dirNumber := isiAvc.Name()
				isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s", folderSource, isiOutput.Name(), avcDirName, dirNumber))
				if err != nil {
					zap.L().Error("os.ReadDir isiAvcDir (HlsPlaylistUploader)", zap.Error(err))
					return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
				}
				for _, fileInsideDirNumber := range isiVideoDir {
					fileInsideDirNumberName := fileInsideDirNumber.Name()
					info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/"+isiOutput.Name()+"/"+avcDirName+"/"+dirNumber+"/"+fileInsideDirNumberName, folderSource+"/"+isiOutput.Name()+"/"+avcDirName+"/"+dirNumber+"/"+fileInsideDirNumberName, minio.PutObjectOptions{ContentType: contentType})
					if err != nil {
						zap.L().Error("minioClient.FPutObject (HlsPlaylistUploader)", zap.Error(err))
						return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
					}
					zap.L().Info(fmt.Sprintf("Successfully uploaded %s of size %d\n", info))
				}
			}
		}
	}

	// upload all audio
	for _, isiOutput := range isiOutputDir {

		if isiOutput.Name() == "audio" {
			isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s", folderSource, isiOutput.Name()))
			if err != nil {
				zap.L().Error("os.ReadDir audio isiOutputDir (HlsPlaylistUploader)", zap.Error(err))
				return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
			}
			endDirName := isiVideoDir[0].Name()
			mpa4Dir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s", folderSource, isiOutput.Name(), endDirName))
			if err != nil {
				zap.L().Error("os.ReadDir audio isiVideoDir  (HlsPlaylistUploader)", zap.Error(err))
				return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
			}

			isiMpa4Dir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s", folderSource, isiOutput.Name(), endDirName, mpa4Dir[0].Name()))
			if err != nil {
				zap.L().Error("os.ReadDir isiMpa4Dir  (HlsPlaylistUploader)", zap.Error(err))
				return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
			}
			for _, isiMpa4Dir := range isiMpa4Dir {
				dirNumber := isiMpa4Dir.Name()
				isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s/%s", folderSource, isiOutput.Name(), endDirName, mpa4Dir[0].Name(), dirNumber))
				if err != nil {
					zap.L().Error("os.ReadDir  isiMpa4Dir  (HlsPlaylistUploader)", zap.Error(err))
					return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
				}
				for _, fileInsideDirNumber := range isiVideoDir {
					fileInsideDirNumberName := fileInsideDirNumber.Name()
					info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/"+isiOutput.Name()+"/"+endDirName+"/"+mpa4Dir[0].Name()+"/"+dirNumber+"/"+fileInsideDirNumberName, folderSource+"/"+isiOutput.Name()+"/"+endDirName+"/"+mpa4Dir[0].Name()+"/"+dirNumber+"/"+fileInsideDirNumberName, minio.PutObjectOptions{ContentType: contentType})
					if err != nil {
						zap.L().Error("minioClient.FPutObject isiVideoDir  (HlsPlaylistUploader)", zap.Error(err))
						return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
					}
					zap.L().Info(fmt.Sprintf("Successfully uploaded %s of size %d\n", info))
				}
			}
		}
	}

	return nil
}
