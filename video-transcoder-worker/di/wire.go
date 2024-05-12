// go:build wireinject
//go:build wireinject
// +build wireinject

package di

import (
	"lintang/video-processing-worker/biz/dal/rabbitmq"
	"lintang/video-processing-worker/biz/service"
	"lintang/video-processing-worker/biz/webapi"
	"lintang/video-processing-worker/config"

	"github.com/google/wire"
)

var ProviderSet wire.ProviderSet = wire.NewSet(
	service.NewTranscoderService,
	webapi.NewMinioAPI,
	webapi.NewDkronAPI,
	rabbitmq.NewMetadataMQ,

	wire.Bind(new(service.DkronCLIAPI), new(*webapi.DkronAPI)),
	wire.Bind(new(service.MinioAPI), new(*webapi.MinioAPI)),
	wire.Bind(new(service.MetadataMQ), new(*rabbitmq.MetadataMQ)),
)

func InitTranscoderService(cfg *config.Config, rmq *rabbitmq.RabbitMQ) *service.TranscoderService {
	wire.Build(
		ProviderSet,
	)
	return nil
}
