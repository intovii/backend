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
	if err := uc.Repo.GetAdvertismentReviews(ctx, advertisment); err != nil{
		uc.log.Error("fail to get Reviews by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetAdvertismentPhotos(ctx, advertisment); err != nil{
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


func (uc *Usecase) GetProfileUserAllInfo(ctx context.Context, user *entities.User) error {
	if exist, err := uc.Repo.IsUserExist(ctx, user); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return errors.New("user does not exist")
	}
	if err := uc.Repo.GetUserInfo(ctx, user); err != nil{
		uc.log.Error("fail to get user profile info", zap.Error(err))
		return err
	}

	return nil
}

func (uc *Usecase) GetProfileUserStatistics(ctx context.Context, uID uint64) (*[]*entities.ProfileStatistic, error) {
	var stats []*entities.ProfileStatistic
	if exist, err := uc.Repo.IsUserExist(ctx, &entities.User{ID:uID}); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return nil, errors.New("user does not exist")
	}

	if err := uc.Repo.GetProfileUserStatistics(ctx, uID, &stats); err != nil {
		uc.log.Error("fail to ads by buyer id", zap.Error(err))
		return nil, err
	}
	for _, stat := range stats{
		var err error
		stat.AdPhotoPath, err = uc.Repo.GetMainAdPhotoByAdID(ctx, stat.AdID)
		if err != nil{
			uc.log.Error("fail to get Photos by Advertisment ID", zap.Error(err))
		}
		if exist, err := uc.Repo.IsReviewExistByDealID(ctx, stat.DealID); err != nil || !exist {
			uc.log.Error("review does not exist", zap.Error(err))
		}
		if err := uc.Repo.GetStatisticAdReviewMark(ctx, stat); err != nil{
			uc.log.Error("fail to get Review by Deal ID", zap.Error(err))
		}

	}
	
	return &stats, nil
}

func (uc *Usecase) GetProfileMyAdvertisments(ctx context.Context, uID uint64) (*[]*entities.MyAdvertisement, error) {
	var advertisements []*entities.MyAdvertisement
	if exist, err := uc.Repo.IsUserExist(ctx, &entities.User{ID:uID}); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return nil, errors.New("user does not exist")
	}

	if err := uc.Repo.GetProfileMyAdvertisments(ctx, uID, &advertisements); err != nil {
		uc.log.Error("fail to get ads by user id", zap.Error(err))
		return nil, err
	}
	
	for _, ad := range advertisements{
		var err error
		ad.AdPhotoPath, err = uc.Repo.GetMainAdPhotoByAdID(ctx, ad.AdID)
		if err != nil {
			uc.log.Error("fail to get Photos by Advertisment ID", zap.Error(err))
		}
	}
	
	return &advertisements, nil
}

func (uc *Usecase) GetProfileReviews(ctx context.Context, uID uint64) (*[]*entities.ProfileReview, error) {
	var reviews []*entities.ProfileReview
	if exist, err := uc.Repo.IsUserExist(ctx, &entities.User{ID:uID}); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return nil, errors.New("user does not exist")
	}

	if err := uc.Repo.GetProfileReviews(ctx, uID, &reviews); err != nil {
		uc.log.Error("fail to get reviews by user id", zap.Error(err))
		return nil, err
	}
	
	// for _, ad := range reviews{
	// 	var err error
	// 	ad.AdPhotoPath, err = uc.Repo.GetMainAdPhotoByAdID(ctx, ad.AdID)
	// 	if err != nil {
	// 		uc.log.Error("fail to get Photos by Advertisment ID", zap.Error(err))
	// 	}
	// }
	
	return &reviews, nil
}