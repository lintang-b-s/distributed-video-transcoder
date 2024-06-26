// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"github.com/google/wire"
	"lintang/video-transcoder-api/biz/dal/mongodb"
	"lintang/video-transcoder-api/biz/service"
	"lintang/video-transcoder-api/biz/webapi"
	"lintang/video-transcoder-api/config"
)

// Injectors from wire.go:

func InitTranscoderService(cfg *config.Config, mongo *mongodb.MongoDB) *service.TranscoderService {
	dkronAPI := webapi.NewDkronAPI(cfg)
	minioAPI := webapi.NewMinioAPI(cfg)
	metadataRepo := mongodb.NewMetadataRepo(mongo)
	transcoderService := service.NewTranscoderService(dkronAPI, minioAPI, metadataRepo)
	return transcoderService
}

// wire.go:

var ProviderSet wire.ProviderSet = wire.NewSet(service.NewTranscoderService, webapi.NewDkronAPI, webapi.NewMinioAPI, mongodb.NewMetadataRepo, wire.Bind(new(service.DkronAPI), new(*webapi.DkronAPI)), wire.Bind(new(service.MinioAPI), new(*webapi.MinioAPI)), wire.Bind(new(service.MetadataRepo), new(*mongodb.MetadataRepo)))
