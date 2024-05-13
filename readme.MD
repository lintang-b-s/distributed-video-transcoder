
## Description
video transcoding menggunakan ffmpeg & bento4 (buat convert ko protocol mpeg-dash) memakan resource cpu & memory yang sangat besar dan memakan waktu sangat lama, alangkah baiknya transcoding dikerjakan beberapa workers sekaligus. Di project ini proses untuk produce multi-bitrate segmented (ffmpeg) video dikerjakan oleh beberapa worker (kalau deploy pake kubernetes)

## Architecture
![distributed video transcoder](https://res.cloudinary.com/dex4u3rw4/image/upload/v1715497709/distributed_video_transcoer/distributed_video_transcoder_6_ymqm3w.png)

## Quick Start (Minikube/k8s, transcoding worker ada 3) (yg kujelasin works di minikube)
- setup mongodb operator, rabbitmq operator, minio, dkron, nginx-ingress (eksekusi semua commadn dibawah sebelum deploy app)
```

0. pastikan udah install, helm, minikube, docker
0b. minikube start --cpus max --memory=9000mb (kalau sebelumnya pernah init cluster minikube, delete dulu `minikube delete`)
0c. helm repo add mongodb https://mongodb.github.io/helm-charts
0d. helm install community-operator mongodb/community-operator --namespace mongodb  --create-namespace

1. kubectl apply -f k8s-deployment/mongodb/mongo_cr.yaml --namespace mongodb
1b. kubectl describe pod example-mongodb-0 -n mongodb (tunggu sampai pod started)
1c. kubectl get pod -n mongodb (tunggu sampai example-mongodb-0, example-mongodb-1, example-mongodb-2 started)

1d.dapetin connection string mongodb:  kubectl get secret example-mongodb-admin-admin    -n mongodb \
-o json | jq -r '.data | with_entries(.value |= @base64d)'

1e. contoh hasil dari command diatas: 
{
  "connectionString.standard": "mongodb://admin:lintang@example-mongodb-0.example-mongodb-svc.mongodb.svc.cluster.local:27017,example-mongodb-1.example-mongodb-svc.mongodb.svc.cluster.local:27017,example-mongodb-2.example-mongodb-svc.mongodb.svc.cluster.local:27017/admin?replicaSet=example-mongodb&ssl=false",
  "connectionString.standardSrv": "mongodb+srv://admin:lintang@example-mongodb-svc.mongodb.svc.cluster.local/admin?replicaSet=example-mongodb&ssl=false",
  "password": "lintang",
  "username": "admin"
}


1f. copy connectionString.standarrd ke ennvironement k8s-deployment/app/api-deployment.yaml (MONGO_URL)



1j. minikube addons enable ingress

--- rabbitmq----
2. kubectl apply -f https://github.com/rabbitmq/cluster-operator/releases/latest/download/cluster-operator.yml 
2b. tunggu sampai operator running  (kubectl get pod -n rabbitmq-system)
2c. dapetin user rabbitmq : kubectl -n default  get secret rabbitmq-default-user -o jsonpath="{.data.username}" | base64 --decode

2d. dapetin password rabbitmq:  kubectl -n default get secret rabbitmq-default-user -o jsonpath="{.data.password}" | base64 --decode

2f. kubectl get service rabbitmq -o jsonpath='{.spec.clusterIP}'

2g. kubectl apply -f k8s-deployment/rabbitmq/rmq_cr.yaml
conURL := amqp://<yang_didapet_dari_2c>:<yang_didapet_dari2d>@<yang_didapet_dari_2g>:5672/
contoh:
connURL= amqp://default_user_Z4KRpZEzc-7wictHAsl:0vpV52fDOzbx2UtHFMRDotjw27pvzB1w@10.102.74.165:5672/

2i. copy conn url ke ennvironement k8s-deployment/app/api-deployment.yaml dan k8s-deployment/app/transcoder-worker.yaml (RABBIT_URL)

--- minio---

3a.  kubectl create ns minio-dev
3b. kubectl apply -f k8s-deployment/minio/minio.yaml
3c. connURL = minio-service.minio-dev.svc.cluster.local:9000
3d. copy connURL ke environment api-deployment.yaml dan k8s-deployment/app/transcoder-worker.yaml  (BASE_URL_MINIO)

3f. minikube ip && buka <ip_minikube>:30023 di browser & login username=lintang, password=lintang123

3g. buat access keys dg expiry date > 5bulan dari sekarang. copy access key & secret key ke environment  api-deployment.yaml dan k8s-deployment/app/transcoder-worker.yaml  (ACC_KEY_MINIO,SECRET_KEY_MINIO )

3h. buat bucket dengan nama "oti-be-bucket" di minio ( <ip_minikube>:30023) & harus public  typenya


---- dkron ----
4a. kubectl  create configmap dkroncurl  --from-file ./dkron_curl.sh
4b. kubectl apply -f k8s-deployment/dkron/dkron-deployment.yaml
4c. kubectl get pod 

4d. kubectl exec -it <nama_pod_dkron>  -- bash -c "cp curl/* bisa/ && chmod 777 bisa/dkron_curl.sh && bisa/dkron_curl.sh"

4e. connURL = http://dkron-svc:8080/v1/jobs
4f. copy conURl dkron ke environment  api-deployment.yaml dan k8s-deployment/app/transcoder-worker.yaml (DKRON_URL) 


---nginx ingress---
5a. minikube addons enable ingress
5b. minikube ip (copy ip nya)



---reverse proxy minio---- (biar presignedurl minio bisa diakses di laptop)
6a. kubectl create configmap reverse-proxy-minio-config --from-file ./k8s-deployment/app/nginx.conf
```

-- deploy app nya
```
0. paste minikube ip ke environent Minikube_IP k8s-deployment/app/api-deployment.yaml & k8s-deployment/app/transcoder-worker.yaml 
1. kubectl apply -f k8s-deployment/app
2. dapetin url transcoder-api: minikube service lintang-transcoder-api --url
3. import postman collection  (distributed_video_transcoder Kubernetes.postman_collection)
4. request ke endpoint getPresignedURl minio
5. upload video ke endpoint yang ada di response endpoint getPresignedURl (method PUT && body typenya binary di postman/yaak)
6. kubectl get pod  . terus `kubectl logs <nama_pod_transcder_worker>` buat liat log process transcodingnya
7. setelah selesai upload semua resolusi (240p-1080p),  request ke endpoint get metadatas, buat dapetin video_url minio 
8. cd videoplayer/my-app && npm i && npm run dev, buka localhost:3000, masukin video url yg di response getMetaadata tadi, video player bakal mulai video nya dg protokol mpeg-dash
```



## Quick Start (Docker Compose, Transcoding worker cuma 1)
```
0.  docker network create lintang-video-transcoder 
1. docker compose up -d 
1b. sh ./setup_dkron.sh
2. buat access key & secretkey di localhost:9001 (minio) dengan EXPIRY DATE > 5 bulan dari sekarang dan copy & paste  access key & secretkey  di docker-compose-app.yaml (ACC_KEY_MINIO=<access_key_mu> &&
SECRET_KEY_MINIO=<secret_key_mu>) ke environement container transcoder-worker & api

2b. docker compose -f  docker-compose-app.yaml up -d   
(kalau transcoder worker dan api exited , docker compose  -f  docker-compose-app.yaml  up -d  lagi, karena rabbitmq belum startup)
3. buat bucket dengan nama "oti-be-bucket" di minio (localhost:9001) & harus public  typenya

4. import postman collection
5. get presigned url minio di endpoint postman
6. upload video pake url method PUT yg ada di response endpoint diatas  (pake binary type di postman/yaak), usahakan ukuran video kurang dari 15mb,  semakin besar ukuran file semakin lama transcodingnya
7. docker logs transcoder-worker -f, buat liat proses transcodingnya. 
8. setelah selesai transcoding, get video URL nya lewat URL http://192.168.49.2:31203/api/v1/tenflixdi postman
9. cd videoplayer/my-app && npm i && npm run dev, buka localhost:3000, masukin video url yg di response getMetaadata tadi, video player bakal mulai video nya dg protokol mpeg-dash

```

## DEMO
- demo video youtube: https://youtu.be/4y4NWCoMpAg
- demo video cloudinary: https://res.cloudinary.com/dex4u3rw4/video/upload/v1715572201/distributed_video_transcoer/distributed-video-transcoder_FL7gmrAm_tiu8su.webm

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




