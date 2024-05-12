package domain

import "go.mongodb.org/mongo-driver/bson/primitive"


type VideoMetadata struct {
	Id primitive.ObjectID `json:"_id"`
	VideoURL string `json:"video_url"`
	Thumbnail string `json:"thumbnail"`
}