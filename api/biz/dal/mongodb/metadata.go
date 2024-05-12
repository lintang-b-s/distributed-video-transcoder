package mongodb

import (
	"context"
	"lintang/video-transcoder-api/biz/dal/domain"

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

	_, err := coll.InsertOne(context.TODO(), coll)
	if err != nil {
		zap.L().Error("coll.InsertOne (MetadataRepo)", zap.Error(err ))
		return domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}
	return nil
}

func (r MetadataRepo) GetAll(ctx context.Context) ([]domain.VideoMetadata, error) {
	coll := r.db.Cli.Database("lintang_video").Collection("metadata")

	cursor, err := coll.Find(context.TODO(), bson.M{})
	if err != nil {
		zap.L().Debug(" metadata not found", zap.Error(err))

		return nil, domain.WrapErrorf(err, domain.ErrNotFound, "metadata video belum ada coi")
	}

	results := []bson.M{}
	if err = cursor.All(context.TODO(), &results); err != nil {
		zap.L().Error(" cursor.All (MetadataRepo)", zap.Error(err))
		return nil, domain.WrapErrorf(err, domain.ErrInternalServerError, domain.MessageInternalServerError)
	}

	var metadatas []domain.VideoMetadata
	for _ ,res := range results {
		metadatas = append(metadatas, domain.VideoMetadata{
			ID: res["id"].(string),
			VideoURL: res["video_url"].(string),
			Thumbnail: res["thumbnail"].(string),
		})
	}
	return metadatas, nil 
}
