# Step 1: Modules caching
FROM golang:1.22.2-alpine as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM golang:1.22.2-alpine as builder
COPY --from=modules /go/pkg /go/pkg
COPY . /app
WORKDIR /app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /bin/app .

# Step 3: Final
FROM ubuntu:23.10
# install ffmpeg & bento4 buat video transcoding
RUN apt-get update
RUN apt-get update
RUN apt-get install ffmpeg -y
RUN apt-get install wget -y 
RUN wget https://www.bok.net/Bento4/binaries/Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip
RUN apt-get install unzip -y
RUN unzip Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip
ENV PATH="${PATH}:/Bento4-SDK-1-6-0-641.x86_64-unknown-linux/bin"
COPY --from=builder /app/config /config
COPY --from=builder /app/k8s.env .
COPY --from=builder /app/.env .
COPY --from=builder /bin/app /bin/app
CMD ["/bin/app"]


### gakbisa yg dibawah

# RUN apt-get install git -y
# RUN git clone  https://github.com/axiomatic-systems/Bento4.git
# # RUN wget https://github.com/axiomatic-systems/Bento4/archive/refs/heads/master.zip # gakbisa ini
# RUN apt-get install  build-essential -y
# RUN   apt-get install cmake -y
# RUN cd Bento4 && mkdir cmakebuild &&  cd cmakebuild  && cmake -DCMAKE_BUILD_TYPE=Release .. && make

# alpine gak bisa mp4dash , 
# FROM alpine:3.19.1
# # install ffmpeg & bento4 buat video transcoding
# RUN apk update
# RUN apk upgrade
# # RUN apk add --no-cache ffmpeg
# RUN apk add ffmpeg

# RUN  apk add --no-cache make
# RUN apk add cmake
# RUN apk add wget 
# RUN wget https://www.bok.net/Bento4/binaries/Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip
# RUN apt-get install unzip -y
# RUN unzip Bento4-SDK-1-6-0-641.x86_64-unknown-linux.zip
# ENV PATH="${PATH}:/Bento4-SDK-1-6-0-641.x86_64-unknown-linux/bin"

# # RUN apk add bento4
# COPY --from=builder /app/config /config
# COPY --from=builder /app/.env .
# COPY --from=builder /bin/app /bin/app
# CMD ["/bin/app"]

