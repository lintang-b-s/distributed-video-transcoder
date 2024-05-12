package rabbitmq

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"lintang/video-processing-worker/biz/domain"
	"time"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

type MetadataMQ struct {
	ch *amqp.Channel
}

func NewMetadataMQ(rmq *RabbitMQ) *MetadataMQ {
	return &MetadataMQ{
		rmq.Channel,
	}
}

func (m *MetadataMQ) PublishNewMetadata(ctx context.Context, d domain.VideoMetadataMessage) error {
	return m.publish(ctx, "metadata.new", d)
}

func (m *MetadataMQ) publish(ctx context.Context, routingKey string, event interface{}) error {

	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(event); err != nil {
		zap.L().Error("gob.NewEncoder(&b).Encode(event)", zap.Error(err))
		return err
	}

	err := m.ch.Publish(
		"metadata", // exchange
		routingKey,        // routing key
		false,
		false,
		amqp.Publishing{
			AppId:       "metadata-transcoding-worker",
			ContentType: "application/x-encoding-gob",
			Body:        b.Bytes(),
			Timestamp:   time.Now(),
		})
	if err != nil {
		zap.L().Error("m.ch.Publish: ", zap.Error(err))
		return err
	}
	zap.L().Info(fmt.Sprintf("sukses mengirimkan message ke rabbit mq routing key %s", routingKey))
	return nil
}
