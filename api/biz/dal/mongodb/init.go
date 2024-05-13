package mongodb

import (
	"context"
	"fmt"
	"lintang/video-transcoder-api/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type MongoDB struct {
	Cli *mongo.Client
}

func NewMongoDB(cfg *config.Config) *MongoDB {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongo.MongoURL))
	if err != nil {
		zap.L().Fatal("mongo.Connect", zap.Error(err))
	}
	zap.L().Info(fmt.Sprintf("connected to mongodb ! %s", cfg.Mongo.MongoURL))
	return &MongoDB{Cli: client}
}

func (m *MongoDB) Close(ctx context.Context) {
	m.Cli.Disconnect(ctx)
}
