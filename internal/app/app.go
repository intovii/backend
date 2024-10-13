package app

import (
	"backend/config"
	"backend/internal/domain/delivery"
	"backend/internal/domain/repository"
	"backend/internal/domain/usecase"
	"context"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func New() *fx.App {
	return fx.New(
		fx.Options(
			repository.New(),
			usecase.New(),
			server.New(),
		),
		fx.Provide(
			context.Background,
			config.NewConfig,
			zap.NewProduction,
		),
		fx.WithLogger(
			func(log *zap.Logger) fxevent.Logger {
				return &fxevent.ZapLogger{Logger: log}
			},
		),
	)
}