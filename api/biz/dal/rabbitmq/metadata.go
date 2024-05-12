package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"lintang/video-transcoder-api/biz/dal/domain"
	"lintang/video-transcoder-api/biz/dal/mongodb"

	"go.uber.org/zap"
)

type MetadatListenner struct {
	rmq *RabbitMQ
	metadataRepo *mongodb.MetadataRepo
	done chan struct{}
}

func NewMetadataListener(rmq *RabbitMQ, m *mongodb.MetadataRepo, d chan struct{}) *MetadatListenner {
	return &MetadatListenner{rmq, m, d }
}

func (l MetadatListenner) ListenAndServe() error {
	queue, err := l.rmq.Channel.QueueDeclare(
		"metadata",
		false,
		false,
		true,
		false,
		nil,
	)

	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	err = l.rmq.Channel.QueueBind(
		queue.Name,
		"metadata.new",
		"metadata",
		false,
		nil,
	)

	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	msgs, err := l.rmq.Channel.Consume(
		queue.Name,
		"video-transcoder-api",
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	var nack bool 
	go func() {
		for msg := range msgs {
			zap.L().Info(fmt.Sprintf(`received message: %s`, msg.RoutingKey))

			switch msg.RoutingKey {
			case "metadata.new":
				metadata, err := decodeMetadataMessage(msg.Body)
				if err != nil {
					zap.L().Error("decodeMetadataMessage (ListenAndServe)", zap.Error(err ))
				}

				err = l.metadataRepo.Insert(context.Background(), domain.VideoMetadata{VideoURL: metadata.VideoURL, Thumbnail: metadata.Thumbnail})
				if err != nil {
					nack = true
				}
			}	

			if nack {
				zap.L().Info("nack")
				_ = msg.Nack(false, nack)
			}else {
				zap.L().Info("ack")
				_ = msg.Ack(false)
			}
		}

		l.done <- struct{}{}
	}()
	return nil 
}



type VideoMetadataMessage struct {
	VideoURL string `json:"video_url"`
	Thumbnail string `json:"thumbnail"`
}


func decodeMetadataMessage(b []byte) (VideoMetadataMessage, error) {
	var res VideoMetadataMessage
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		zap.L().Error("NewDecoder (decodeMetadataMessage) (MetadataListener)", zap.Error(err))
		return VideoMetadataMessage{}, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return res, nil 
}