
## Description
video transcoding menggunakan ffmpeg & bento4 memakan resource cpu & memory yang sangat besar dan memakan waktu sangat lama, alangkah baiknya transcoding dikerjakan beberapa workers sekaligus. Di project ini proses untuk produce multi-bitrate segmented (ffmpeg) video dikerjakan oleh beberapa worker (kalau deploy pake kubernetes)

## Architecture
![distributed video transcoder](https://res.cloudinary.com/dex4u3rw4/image/upload/v1715497709/distributed_video_transcoer/distributed_video_transcoder_6_ymqm3w.png)


## Quick Start
```
0.  docker network create lintang-video-transcoder 
1. docker compose up -d 
1b. sh ./setup_dkron.sh
2. buat access key & secretkey di localhost:9001 (minio) dengan EXPIRY DATE > 5 bulan dari sekarang dan copy & paste  access key & secretkey  di docker-compose-app.yaml (ACC_KEY_MINIO=<access_key_mu> &&
SECRET_KEY_MINIO=<secret_key_mu>) ke environement container transcoder-worker & api

2b. docker compose -f  docker-compose-app.yaml up -d   
(kalau transcoder worker dan api exited , docker compose  -f  docker-compose-app.yaml  up -d  lagi, karena rabbitmq belum startup)
3. buat bucket dengan nama "oti-be-bucket" di minio (localhost:9001) & harus public 

4. import postman collection
5. get presigned url minio di endpoint postman
6. upload video pake url method PUT yg ada di response endpoint diatas  (pake binary type di postman/yaak), usahakan ukuran video kurang dari 15mb,  semakin besar ukuran file semakin lama transcodingnya
7. docker logs transcoder-worker -f, buat liat proses transcodingnya. 
8. setelah selesai transcoding, get video URL nya lewat URL http://localhost:9996/api/v1/tenflix di postman
```

## Tugas transcoding workers
### ffmpeg produce multi-bit-rate version
```


setiap worker akan menerima request dari dkron untuk pull video  <input_file_name>.mp4 yang abis diupload user di object folder <input_file_name> minio, terus setiap worker akan menjalannkan salah satu dari command dibawah

ffmpeg -i <input_file_name>.mp4 -c:v h264 -crf 22 -tune film -profile:v main -level:v 4.0 -maxrate 5000k -bufsize 10000k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 44100 -b:a 128k -ac 2 -pix_fmt yuv420p -movflags +faststart 1080.mp4 -s 1280x720 && mp4fragment 1080.mp4 1080-f.mp4


ffmpeg -i <input_file_name>.mp4 -c:v h264 -crf 24 -tune film -profile:v main -level:v 4.0 -maxrate 2500k -bufsize 5000k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 44100 -b:a 128k -ac 2 -pix_fmt yuv420p -movflags +faststart 720.mp4 -s 854x480 && mp4fragment 720.mp4 720-f.mp4

ffmpeg -i <input_file_name>.mp4 -c:v h264 -crf 30 -tune film -profile:v main -level:v 4.0 -maxrate 1250k -bufsize 2500k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 44100 -b:a 96k -ac 2 -pix_fmt yuv420p -movflags +faststart 480.mp4 -s 640x360  && mp4fragment 480.mp4 480-f-.mp4

ffmpeg -i <input_file_name>.mp4 -c:v h264 -crf 33 -tune film -profile:v main -level:v 4.0 -maxrate 900k -bufsize 1800k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 44100 -b:a 96k -ac 2 -pix_fmt yuv420p -movflags +faststart 360.mp4 -s 320x24  && mp4fragment 360.mp4 360-f.mp4
 
ffmpeg -i <input_file_name>.mp4 -c:v h264 -crf 36 -tune film -profile:v main -level:v 4.0 -maxrate 625k -bufsize 1250k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 22050 -b:a 64k -ac 1 -pix_fmt yuv420p -movflags +faststart 240.mp4  && mp4fragment 240.mp4 240-f.mp4


nah command diatas menghasilkan video 1080p-f.mp4, 720p-f.mp4, 480p-f.mp4, 360p-f.mp4, 240p-f.mp4 yang maana akan dikirim ke minio buecket ke folder <input_file_name>, setelah itu worker akan create cron job ke dkron (buat dikirim ke transcoding worker & worker jalannin task bento4 generate dash playlist dari 5 file yang dipull dari minio)

```


### bento4 generate dash playlist dari segmented mp4s
```
salah satu worker akan dihit oleh dkron untuk pull file 1080p-f.mp4, 720p-f.mp4, 480p-f.mp4, 360p-f.mp4, 240p-f.mp4 dari object folder <input_file_name> di minio, lalu jalanin command dibawah 

`
mp4dash 240-f.mp4 360-f.mp4 480-f.mp4 720-f.mp4 1080-f.mp4
`
keempat command diatas masing masing menghasilkan folder output (ada 2 folder audio dan video)

nah, hasil  generated folder outputnya di upload lagi ke object folder <input_file_name> di minio, dan workernya ngirim transcoding completion message ke rabbitmq 

```


## Ref
- command command ffmpeg && bento4: https://dev.to/siumhossain/producing-dash-streming-using-ffmpeg-and-bento4-316n
-  https://bytebytego.com/courses/system-design-interview/design-youtube






ffmpeg -i lintang-joker.mp4 -c:v h264 -crf 22 -tune film -profile:v main -level:v 4.0 -maxrate 5000k -bufsize 10000k -r 25 -keyint_min 25 -g 50 -sc_threshold 0 -c:a aac -ar 44100 -b:a 128k -ac 2 -pix_fmt yuv420p -movflags +faststart 1080.mp4 -s 1280x720 && mp4fragment 1080.mp4 1080-f.mp4
