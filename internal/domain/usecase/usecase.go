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

func (uc *Usecase) GetAdvertismentAllInfo(ctx context.Context, advertisment *entities.Advertisment) error {
	if exist, err := uc.Repo.IsAdExist(ctx, advertisment); err != nil || !exist {
		uc.log.Error("advertisment does not exist", zap.Error(err))
		return errors.New("advertisment does not exist")
	}

	if err := uc.Repo.GetAdvertismentAllInfo(ctx, advertisment); err != nil{
		uc.log.Error("fail to get Advertisment", zap.Error(err))
		return err
	}
	if exist, err := uc.Repo.IsUserExist(ctx, &advertisment.User); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return errors.New("user does not exist")
	}
	if err := uc.Repo.GetUserInfo(ctx, &advertisment.User); err != nil{
		uc.log.Error("fail to get seller info by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetReviews(ctx, advertisment); err != nil{
		uc.log.Error("fail to get Reviews by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetPhotos(ctx, advertisment); err != nil{
		uc.log.Error("fail to get Photos by Advertisment ID", zap.Error(err))
		return err
	}

	return nil
}