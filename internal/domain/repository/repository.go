package repository

import (
	"go.uber.org/fx"
	"backend/internal/domain/repository/postgres"
)

func New() fx.Option {
	return fx.Module("repository",
		fx.Provide(
			postgres.NewRepository,
		),
		fx.Invoke(
			func(lc fx.Lifecycle, a *postgres.Repository) {
				lc.Append(fx.Hook{
					OnStart: a.OnStart,
					OnStop:  a.OnStop,
				})
			},
		),
	)
}
