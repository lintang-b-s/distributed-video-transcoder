package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func main() {
	// ctx := context.Background()
	endpoint := "127.0.0.1:9000"
	accessKeyID := "kIiDoM7RUmofIPTyQlQQ"
	secretAccessKey := "I2ZJHJuGwYojd5yG47gmsdkNyTpnfrlhFALNhSeG"
	useSSL := false

	// Initialize minio client object.
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", minioClient) // minioClient is now set up
	log.Printf("Berhasill")

	if err != nil {
		log.Fatalln(err)
	}

	// // Make a new bucket called testbucket.
	// bucketName := "testbucket"
	// location := "us-east-1" // harus us-east-1
	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"stream.mpd\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err := minioClient.PresignedGetObject(context.Background(), "oti-be-bucket", "dash/output/stream.mpd", time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	// get object
	object, err := minioClient.GetObject(context.Background(), "oti-be-bucket", "dash/output/stream.mpd", minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return
	}
	defer object.Close()

	localFile, err := os.Create("./tes.mpd")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer localFile.Close()

	if _, err = io.Copy(localFile, object); err != nil {
		fmt.Println(err)
		return
	}

	reqParams = make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"stream.mpd\"")

	// Generates a presigned url which expires in a day.
	presignedURL, err = minioClient.PresignedGetObject(context.Background(), "oti-be-bucket", "video/output/stream.mpd", time.Second*24*60*60, reqParams)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

}
