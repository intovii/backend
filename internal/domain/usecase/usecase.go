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

// func (uc *Usecase) IsPhoneExist(ctx context.Context, user *entities.SqlUser) error {
// 	if user.NumberPhone.Valid == true {
// 		if exist, err := uc.Repo.IsPhoneExist(ctx, user); err != nil || !exist {
// 			return fmt.Errorf("the phone number is already there: %w", err)
// 		}
// 	}
// 	return nil
// }

// func (uc *Usecase) IsUsernameExist(ctx context.Context, user *entities.SqlUser) error {
// 	if user.Username.Valid == true {
// 		if exist, err := uc.Repo.IsUsernameExist(ctx, user); err != nil || !exist {
// 			return fmt.Errorf("the phone number is already there: %w", err)
// 		}
// 	}
// 	return nil
// }
// func (uc *Usecase) CreateUser(ctx context.Context, user *entities.SqlUser) error {
// 	//if exist, err := uc.Repo.IsUserExist(ctx, user); err != nil {
// 	//	return fmt.Errorf("failed to check if user exists: %w", err)
// 	//} else if exist {
// 	//	return errors.New("the user already exists")
// 	//}

// 	if err := uc.Repo.CreateUser(ctx, user); err != nil {
// 		uc.log.Error("fail to create User", zap.Error(err))
// 		return err
// 	}
// 	return nil
// }
