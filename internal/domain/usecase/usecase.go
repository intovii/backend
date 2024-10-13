package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/postgres"
	"context"

	"go.uber.org/zap"
)

type Usecase struct {
	log  *zap.Logger
	Repo *postgres.Repository
}

func NewUsecase(logger *zap.Logger, Repo *postgres.Repository) (*Usecase, error) {
	return &Usecase{
		log:  logger,
		Repo: Repo,
	}, nil
}

func (u *Usecase) GetClientWorkouts(ctx context.Context, request *entities.SchedulerGetter) ([]entities.Day, error) {
	days, err := u.Repo.GetClientWorkouts(ctx, request)
	if err != nil {
		u.log.Error("fail to get clients schedule", zap.Error(err))
		return nil, err
	}
	return days, nil
}
