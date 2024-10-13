package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/postgres"
	"context"
	"errors"
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

func (uc *Usecase) GetProductAllInfo(ctx context.Context, product *entities.Advertisment) error {
	if exist, err := uc.Repo.IsAdExist(ctx, product); err != nil || !exist {
		return errors.New("product does not exist")
	}

	if err := uc.Repo.GetProductAllInfo(ctx, product); err != nil{
		uc.log.Error("fail to get Advertisment", zap.Error(err))
		return err
	}
	if exist, err := uc.Repo.IsUserExist(ctx, &product.User); err != nil || !exist {
		return errors.New("user does not exist")
	}
	if err := uc.Repo.GetUserInfo(ctx, &product.User); err != nil{
		uc.log.Error("fail to get Seller by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetReviews(ctx, product); err != nil{
		uc.log.Error("fail to get Seller by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetPhotos(ctx, product); err != nil{
		uc.log.Error("fail to get Seller by Advertisment ID", zap.Error(err))
		return err
	}

	return nil
}