package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// pulling video dari minio
func getUserVideoURL(filename string) (*minio.Object, error) {
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "kIiDoM7RUmofIPTyQlQQ"
	secretAccessKey := "I2ZJHJuGwYojd5yG47gmsdkNyTpnfrlhFALNhSeG"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}


	if err != nil {
		log.Fatalln(err)
	}
	// Make a new bucket called testbucket.
	bucketName := "oti-be-bucket"


	userVideo, err := minioClient.GetObject(ctx, bucketName, fmt.Sprintf("%s/%s.mp4", filename, filename), minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	return userVideo, nil 
}

func getAllBitrateVideoVersion(objectFolderName string) []*minio.Object {
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "kIiDoM7RUmofIPTyQlQQ"
	secretAccessKey := "I2ZJHJuGwYojd5yG47gmsdkNyTpnfrlhFALNhSeG"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}


	if err != nil {
		log.Fatalln(err)
	}
	// Make a new bucket called testbucket.
	bucketName := "oti-be-bucket"

	// mendapatkan semua versi bitraate video
	// ctxWithTimeout, cancel := context.WithTimeout(ctx, 40*time.Second)
	// defer cancel()
	file240, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/240-f.mp4", minio.GetObjectOptions{})
	file360, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/360-f.mp4", minio.GetObjectOptions{})
	file480, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/480-f.mp4", minio.GetObjectOptions{})
	file720, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/720-f.mp4", minio.GetObjectOptions{})
	file1080, err := minioClient.GetObject(ctx, bucketName, objectFolderName+"/1080-f.mp4", minio.GetObjectOptions{})

	if err != nil {
		log.Fatalln(err)
	}
	var objects []*minio.Object
	objects = append(objects, file240, file360, file480, file720, file1080)
	return objects

}

func bitrateVersionVideoUploader(objectFolderName string, fileName string, filePath string) {
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "kIiDoM7RUmofIPTyQlQQ"
	secretAccessKey := "I2ZJHJuGwYojd5yG47gmsdkNyTpnfrlhFALNhSeG"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}


	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called testbucket.
	bucketName := "oti-be-bucket"
	location := "us-east-1" // harus us-east-1

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	contentType := "application/octet-stream"

	// upload single bitrate version ke minio
	info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+fileName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectFolderName, info.Size)

}

// buat upload playlist dash ke minio (object: filename/output)
func hlsPlaylistUploader(objectFolderName string, folderSource string) {
	ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "kIiDoM7RUmofIPTyQlQQ"
	secretAccessKey := "I2ZJHJuGwYojd5yG47gmsdkNyTpnfrlhFALNhSeG"
	useSSL := false

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}


	if err != nil {
		log.Fatalln(err)
	}

	// Make a new bucket called testbucket.
	bucketName := "oti-be-bucket"
	location := "us-east-1" // harus us-east-1

	err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
		if errBucketExists == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}
	contentType := "application/octet-stream"

	isiOutputDir, err := os.ReadDir(folderSource)
	if err != nil {
		log.Fatal(err)
	}
	// upload master dash playlist
	info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/stream.mpd", folderSource+"/stream.mpd", minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Successfully uploaded %s of size %d\n", objectFolderName, info.Size)

	// upload all video
	for _, isiOutput := range isiOutputDir {

		if isiOutput.Name() == "video" {
			isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s", folderSource, isiOutput.Name()))
			if err != nil {
				log.Fatalln(err)
			}
			avcDirName := isiVideoDir[0].Name()
			isiAvcDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s", folderSource, isiOutput.Name(), avcDirName))

			for _, isiAvc := range isiAvcDir {
				dirNumber := isiAvc.Name()
				isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s", folderSource, isiOutput.Name(), avcDirName, dirNumber))
				if err != nil {
					log.Fatalln(err)
				}
				for _, fileInsideDirNumber := range isiVideoDir {
					fileInsideDirNumberName := fileInsideDirNumber.Name()
					info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/"+isiOutput.Name()+"/"+avcDirName+"/"+dirNumber+"/"+fileInsideDirNumberName, folderSource+"/"+isiOutput.Name()+"/"+avcDirName+"/"+dirNumber+"/"+fileInsideDirNumberName, minio.PutObjectOptions{ContentType: contentType})
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Sprintf("Successfully uploaded %s of size %d\n", info)
				}
			}
		}
	}

	// upload all audio
	for _, isiOutput := range isiOutputDir {

		if isiOutput.Name() == "audio" {
			isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s", folderSource, isiOutput.Name()))
			if err != nil {
				log.Fatalln(err)
			}
			endDirName := isiVideoDir[0].Name()
			mpa4Dir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s", folderSource, isiOutput.Name(), endDirName))
			if err != nil {
				log.Fatalln(err)
			}

			isiMpa4Dir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s", folderSource, isiOutput.Name(), endDirName, mpa4Dir[0].Name()))
			if err != nil {
				log.Fatalln(err)
			}
			for _, isiMpa4Dir := range isiMpa4Dir {
				dirNumber := isiMpa4Dir.Name()
				isiVideoDir, err := os.ReadDir(fmt.Sprintf("%s/%s/%s/%s/%s", folderSource, isiOutput.Name(), endDirName, mpa4Dir[0].Name(), dirNumber))
				if err != nil {
					log.Fatalln(err)
				}
				for _, fileInsideDirNumber := range isiVideoDir {
					fileInsideDirNumberName := fileInsideDirNumber.Name()
					info, err := minioClient.FPutObject(ctx, bucketName, objectFolderName+"/"+isiOutput.Name()+"/"+endDirName+"/"+mpa4Dir[0].Name()+"/"+dirNumber+"/"+fileInsideDirNumberName, folderSource+"/"+isiOutput.Name()+"/"+endDirName+"/"+mpa4Dir[0].Name()+"/"+dirNumber+"/"+fileInsideDirNumberName, minio.PutObjectOptions{ContentType: contentType})
					if err != nil {
						log.Fatalln(err)
					}
					fmt.Sprintf("Successfully uploaded %s of size %d\n", info)
				}
			}
		}
	}

	// info, err := minioClient.PutObject(ctx, bucketName, objectFolderName, dir, dirInfo.Size(), minio.PutObjectOptions{ContentType: contentType})
	// if err != nil {
	// 	log.Fatalln(err)
	// }
	// log.Printf("Successfully uploaded %s of size %d\n", objectFolderName, info.Size)

}
