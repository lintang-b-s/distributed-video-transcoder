package mongodb

import (
	"context"
	"lintang/video-transcoder-api/biz/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.uber.org/zap"
)

type MetadataRepo struct {
	db *MongoDB
}

func NewMetadataRepo(db *MongoDB) *MetadataRepo {
	return &MetadataRepo{db}
}

func (r MetadataRepo) Insert(ctx context.Context, m domain.VideoMetadata) error {
	coll := r.db.Cli.Database("lintang_video").Collection("metadata")

	_, err := coll.InsertOne(context.TODO(), m)
	if err != nil {
		zap.L().Error("coll.InsertOne (MetadataRepo)", zap.Error(err))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return nil
}

func (r MetadataRepo) GetAll(ctx context.Context) ([]domain.VideoMetadata, error) {
	coll := r.db.Cli.Database("lintang_video").Collection("metadata")

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		zap.L().Debug(" metadata not found", zap.Error(err))

		return nil, domain.WrapErrorf(err, domain.ErrNotFound, "metadata video belum ada coi")
	}

	// if err = cursor.All(context.TODO(), &results); err != nil {
	// 	zap.L().Error(" cursor.All (MetadataRepo)", zap.Error(err))
	// 	return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	// }
	var metadatas = make([]domain.VideoMetadata, 0)

	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		metadata := domain.VideoMetadata{}
		err = cursor.Decode(&metadata)
		if err != nil {
			zap.L().Error("cursor.Decode (GetAll)", zap.Error(err))
			return nil, err
		}
		metadatas = append(metadatas, metadata)

	}

	return metadatas, nil
}
