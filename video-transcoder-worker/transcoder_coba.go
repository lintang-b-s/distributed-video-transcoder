package main

/*

ini buat coba coba doang heheh ffmpeg & bento4 yang dipake di biz/util/ffmpeg_bento4.go
*/	

// import (
// 	"fmt"
// 	"io"
// 	"log"
// 	"os"
// 	"os/exec"
// 	"time"

// 	"github.com/minio/minio-go/v7"
// )

// func createHLSFromMinioObject(inputFile *minio.Object, fileName string) error {

// 	os.Mkdir(fileName, 0777)
// 	filePath := fmt.Sprintf(`%s/%s.mp4`, fileName, fileName)
// 	mylocalFile, err := os.Create(filePath) // membuat file ke lokal dari minio object
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer mylocalFile.Close()
// 	var stat minio.ObjectInfo
// 	stat, err = inputFile.Stat()
// 	if err != nil {
// 		// kalau video belum sepenuhnya ada di minio
// 		for {
// 			// fetch video dari minio terus menerus (setiap 300 milisecond) sampai object ada di minio
// 			fileInfo, err := getUserVideoURL(fileName)
// 			if stat, err = fileInfo.Stat(); err == nil {
// 				break
// 			}
// 			time.Sleep(300 * time.Millisecond)
// 		}
// 		log.Fatal(err)
// 	}

// 	if _, err = io.CopyN(mylocalFile, inputFile, stat.Size); err != nil {
// 		log.Fatalln(err)
// 	}

// 	ffmpeg1080 := exec.Command(
// 		"ffmpeg",
// 		"-i", filePath,
// 		"-c:v", "h264",
// 		"-crf", "22",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "5000k",
// 		"-bufsize", "10000k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "128k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/1080.mp4", fileName),
// 	)

// 	output, err := ffmpeg1080.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg720 := exec.Command(
// 		"ffmpeg",
// 		"-i", filePath,
// 		"-s", "1280x720",
// 		"-c:v", "h264",
// 		"-crf", "24",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "2500k",
// 		"-bufsize", "5000k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "128k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/720.mp4", fileName),
// 	)

// 	output, err = ffmpeg720.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg480 := exec.Command(
// 		"ffmpeg",
// 		"-i", filePath,
// 		"-s", "854x480",
// 		"-c:v", "h264",
// 		"-crf", "30",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "1250k",
// 		"-bufsize", "2500k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "96k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/480.mp4", fileName),
// 	)

// 	output, err = ffmpeg480.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg360 := exec.Command(
// 		"ffmpeg",
// 		"-i", filePath,
// 		"-s", "640x360",
// 		"-c:v", "h264",
// 		"-crf", "33",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "900k",
// 		"-bufsize", "1800k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "96k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/360.mp4", fileName),
// 	)

// 	output, err = ffmpeg360.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg240 := exec.Command(

// 		"ffmpeg",
// 		"-i", filePath,
// 		"-s", "320x240",
// 		"-c:v", "h264",
// 		"-crf", "36",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "625k",
// 		"-bufsize", "1250k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "22050",
// 		"-b:a", "64k",
// 		"-ac", "1",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/240.mp4", fileName),
// 	)

// 	output, err = ffmpeg240.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment1080 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/1080.mp4", fileName),
// 		fmt.Sprintf("%s/1080-f.mp4", fileName),
// 	)
// 	output, err = mp4fragment1080.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment720 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/720.mp4", fileName),
// 		fmt.Sprintf("%s/720-f.mp4", fileName),
// 	)
// 	output, err = mp4fragment720.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment480 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/480.mp4", fileName),
// 		fmt.Sprintf("%s/480-f.mp4", fileName),
// 	)
// 	output, err = mp4fragment480.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment360 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/360.mp4", fileName),
// 		fmt.Sprintf("%s/360-f.mp4", fileName),
// 	)
// 	output, err = mp4fragment360.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment240 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/240.mp4", fileName),
// 		fmt.Sprintf("%s/240-f.mp4", fileName),
// 	)
// 	output, err = mp4fragment240.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	// upload semua versi video bitrate ke minio
// 	bitrateVersionVideoUploader(fmt.Sprintf("/%s/", fileName), "240-f.mp4", fmt.Sprintf("%s/240-f.mp4", fileName))
// 	bitrateVersionVideoUploader(fmt.Sprintf("/%s/", fileName), "360-f.mp4", fmt.Sprintf("%s/360-f.mp4", fileName))
// 	bitrateVersionVideoUploader(fmt.Sprintf("/%s/", fileName), "480-f.mp4", fmt.Sprintf("%s/480-f.mp4", fileName))
// 	bitrateVersionVideoUploader(fmt.Sprintf("/%s/", fileName), "720-f.mp4", fmt.Sprintf("%s/720-f.mp4", fileName))
// 	bitrateVersionVideoUploader(fmt.Sprintf("/%s/", fileName), "1080-f.mp4", fmt.Sprintf("%s/1080-f.mp4", fileName))

// 	// yang dibawah ini dikerjain oleh worker lain setelah worker 240p selesai upload 240-f ke minio , dg cara bikin cronjob dkron buat assign task ke transcoder-server
// 	allBitrateVideo := getAllBitrateVideoVersion(fmt.Sprintf("/%s", fileName)) // mendapatkan semua video bitrate yg diupload ke minio
// 	files := []string{"240-f.mp4", "360-f.mp4", "480-f.mp4", "720-f.mp4", "1080-f.mp4"}

// 	for i, _ := range files {
// 		os.Mkdir(fileName+"/minio", 0777)
// 		mylocalFile, err := os.Create(fileName + "/minio" + "/" + files[i])
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer mylocalFile.Close()

// 		stat, err := allBitrateVideo[i].Stat()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if _, err = io.CopyN(mylocalFile, allBitrateVideo[i], stat.Size); err != nil {
// 			log.Fatalln(err)
// 		}
// 		allBitrateVideo[i].Close()
// 	}

// 	mp4dash := exec.Command(
// 		"mp4dash",
// 		// fmt.Sprintf("%s/240-f.mp4 %s/360-f.mp4 %s/480-f.mp4 %s/720-f.mp4 %s/1080-f.mp4", fileName, fileName, fileName, fileName, fileName),
// 		fmt.Sprintf("%s/240-f.mp4", fileName+"/minio"),
// 		fmt.Sprintf("%s/360-f.mp4", fileName+"/minio"),
// 		fmt.Sprintf("%s/480-f.mp4", fileName+"/minio"),
// 		fmt.Sprintf("%s/720-f.mp4", fileName+"/minio"),
// 		fmt.Sprintf("%s/1080-f.mp4", fileName+"/minio"),
// 		"-o", fmt.Sprintf("%s/output", fileName),
// 	)

// 	output, err = mp4dash.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	return nil
// }

// func createHLSFromLocal(inputFile string, outputDir string) error {

// 	if err := os.MkdirAll(outputDir, 0755); err != nil {
// 		return fmt.Errorf("failed to create output directory: %v", err)
// 	}

// 	ffmpeg1080 := exec.Command(
// 		"ffmpeg",
// 		"-i", inputFile,
// 		"-c:v", "h264",
// 		"-crf", "22",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "5000k",
// 		"-bufsize", "10000k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "128k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/1080.mp4", outputDir),
// 	)

// 	output, err := ffmpeg1080.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg720 := exec.Command(
// 		"ffmpeg",
// 		"-i", inputFile,
// 		"-s", "1280x720",
// 		"-c:v", "h264",
// 		"-crf", "24",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "2500k",
// 		"-bufsize", "5000k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "128k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/720.mp4", outputDir),
// 	)

// 	output, err = ffmpeg720.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg480 := exec.Command(
// 		"ffmpeg",
// 		"-i", inputFile,
// 		"-s", "854x480",
// 		"-c:v", "h264",
// 		"-crf", "30",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "1250k",
// 		"-bufsize", "2500k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "96k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/480.mp4", outputDir),
// 	)

// 	output, err = ffmpeg480.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg360 := exec.Command(
// 		"ffmpeg",
// 		"-i", inputFile,
// 		"-s", "640x360",
// 		"-c:v", "h264",
// 		"-crf", "33",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "900k",
// 		"-bufsize", "1800k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "44100",
// 		"-b:a", "96k",
// 		"-ac", "2",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/360.mp4", outputDir),
// 	)

// 	output, err = ffmpeg360.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	ffmpeg240 := exec.Command(

// 		"ffmpeg",
// 		"-i", inputFile,
// 		"-s", "320x240",
// 		"-c:v", "h264",
// 		"-crf", "36",
// 		"-tune", "film",
// 		"-profile:v", "main",
// 		"-level:v", "4.0",
// 		"-maxrate", "625k",
// 		"-bufsize", "1250k",
// 		"-r", "25",
// 		"-keyint_min", "25",
// 		"-g", "50",
// 		"-sc_threshold", "0",
// 		"-c:a", "aac",
// 		"-ar", "22050",
// 		"-b:a", "64k",
// 		"-ac", "1",
// 		"-pix_fmt", "yuv420p",
// 		"-movflags", "+faststart",
// 		fmt.Sprintf("%s/240.mp4", outputDir),
// 	)

// 	output, err = ffmpeg240.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment1080 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/1080.mp4", outputDir),
// 		fmt.Sprintf("%s/1080-f.mp4", outputDir),
// 	)
// 	output, err = mp4fragment1080.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment720 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/720.mp4", outputDir),
// 		fmt.Sprintf("%s/720-f.mp4", outputDir),
// 	)
// 	output, err = mp4fragment720.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment480 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/480.mp4", outputDir),
// 		fmt.Sprintf("%s/480-f.mp4", outputDir),
// 	)
// 	output, err = mp4fragment480.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment360 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/360.mp4", outputDir),
// 		fmt.Sprintf("%s/360-f.mp4", outputDir),
// 	)
// 	output, err = mp4fragment360.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	mp4fragment240 := exec.Command(
// 		"mp4fragment",
// 		fmt.Sprintf("%s/240.mp4", outputDir),
// 		fmt.Sprintf("%s/240-f.mp4", outputDir),
// 	)
// 	output, err = mp4fragment240.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	bitrateVersionVideoUploader("/video/", "240-f.mp4", fmt.Sprintf("%s/240-f.mp4", outputDir))
// 	bitrateVersionVideoUploader("/video/", "360-f.mp4", fmt.Sprintf("%s/360-f.mp4", outputDir))
// 	bitrateVersionVideoUploader("/video/", "480-f.mp4", fmt.Sprintf("%s/480-f.mp4", outputDir))
// 	bitrateVersionVideoUploader("/video/", "720-f.mp4", fmt.Sprintf("%s/720-f.mp4", outputDir))
// 	bitrateVersionVideoUploader("/video/", "1080-f.mp4", fmt.Sprintf("%s/1080-f.mp4", outputDir))

// 	allBitrateVideo := getAllBitrateVideoVersion("/video")
// 	files := []string{"240-f.mp4", "360-f.mp4", "480-f.mp4", "720-f.mp4", "1080-f.mp4"}

// 	for i, _ := range files {
// 		os.Mkdir(outputDir+"/minio", 0777)
// 		mylocalFile, err := os.Create(outputDir + "/minio" + "/" + files[i])
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer mylocalFile.Close()

// 		stat, err := allBitrateVideo[i].Stat()
// 		if err != nil {
// 			log.Fatal(err)
// 		}

// 		if _, err = io.CopyN(mylocalFile, allBitrateVideo[i], stat.Size); err != nil {
// 			log.Fatalln(err)
// 		}
// 		allBitrateVideo[i].Close()
// 	}

// 	mp4dash := exec.Command(
// 		"mp4dash",
// 		// fmt.Sprintf("%s/240-f.mp4 %s/360-f.mp4 %s/480-f.mp4 %s/720-f.mp4 %s/1080-f.mp4", outputDir, outputDir, outputDir, outputDir, outputDir),
// 		fmt.Sprintf("%s/240-f.mp4", outputDir+"/minio"),
// 		fmt.Sprintf("%s/360-f.mp4", outputDir+"/minio"),
// 		fmt.Sprintf("%s/480-f.mp4", outputDir+"/minio"),
// 		fmt.Sprintf("%s/720-f.mp4", outputDir+"/minio"),
// 		fmt.Sprintf("%s/1080-f.mp4", outputDir+"/minio"),
// 		"-o", fmt.Sprintf("%s/output", outputDir),
// 	)

// 	output, err = mp4dash.CombinedOutput()
// 	if err != nil {
// 		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
// 	}

// 	return nil
// }