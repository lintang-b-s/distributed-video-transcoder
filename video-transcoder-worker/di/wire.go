// go:build wireinject
//go:build wireinject
// +build wireinject

package di

import (
	"lintang/video-processing-worker/biz/service"
	"lintang/video-processing-worker/biz/webapi"
	"lintang/video-processing-worker/config"

	"github.com/google/wire"
)




var ProviderSet wire.ProviderSet = wire.NewSet(
	service.NewTranscoderService,
	webapi.NewMinioAPI,
	webapi.NewDkronAPI,

	wire.Bind(new(service.DkronCLIAPI), new(*webapi.DkronAPI)),
	wire.Bind(new(service.MinioAPI), new(*webapi.MinioAPI)),
)

func InitTranscoderService(cfg *config.Config) *service.TranscoderService {
	wire.Build(
		ProviderSet,
	)
	return nil 
}


