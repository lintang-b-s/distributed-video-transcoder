// go:build wireinject
//go:build wireinject
// +build wireinject

package di

import (
	"lintang/video-transcoder-api/biz/dal/mongodb"
	"lintang/video-transcoder-api/biz/service"
	"lintang/video-transcoder-api/biz/webapi"
	"lintang/video-transcoder-api/config"

	"github.com/google/wire"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	service.NewTranscoderService,
	webapi.NewDkronAPI,
	webapi.NewMinioAPI,
	mongodb.NewMetadataRepo,

	wire.Bind(new(service.DkronAPI), new(*webapi.DkronAPI)),
	wire.Bind(new(service.MinioAPI), new(*webapi.MinioAPI)),
	wire.Bind(new(service.MetadataRepo), new(*mongodb.MetadataRepo)),
)


func InitTranscoderService(cfg *config.Config, mongo *mongodb.MongoDB) *service.TranscoderService {
	wire.Build(
		ProviderSet,
	)
	return nil
}