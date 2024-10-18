package usecase

import (
	"backend/internal/domain/entities"
	"backend/internal/domain/repository/postgres"
	"context"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"log"
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

	if err := uc.Repo.GetAdvertismentAllInfo(ctx, advertisment); err != nil {
		uc.log.Error("fail to get Advertisment", zap.Error(err))
		return err
	}
	if exist, err := uc.Repo.IsUserExist(ctx, &advertisment.User); err != nil || !exist {
		uc.log.Error("user does not exist", zap.Error(err))
		return errors.New("user does not exist")
	}
	if err := uc.Repo.GetUserInfo(ctx, &advertisment.User); err != nil {
		uc.log.Error("fail to get seller info by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetReviews(ctx, advertisment); err != nil {
		uc.log.Error("fail to get Reviews by Advertisment ID", zap.Error(err))
		return err
	}
	if err := uc.Repo.GetPhotos(ctx, advertisment); err != nil {
		uc.log.Error("fail to get Photos by Advertisment ID", zap.Error(err))
		return err
	}

	return nil
}
func (uc *Usecase) IsUserExit(ctx context.Context, user *entities.User) error {
	if exist, err := uc.Repo.IsUserExist(ctx, user); exist || err != nil {
		return fmt.Errorf("the user_id is already there: %w", err)
	}
	return nil
}

func (uc *Usecase) IsPhoneExist(ctx context.Context, user *entities.User) error {
	if user.NumberPhone != "" {
		exist, err := uc.Repo.IsPhoneExist(ctx, user)
		if err != nil {
			// Если произошла ошибка при запросе, возвращаем её
			return fmt.Errorf("failed to check if phone number exists: %w", err)
		}

		if exist {
			// Если номер телефона уже существует, возвращаем ошибку
			return fmt.Errorf("the phone number %s is already in use", user.NumberPhone)
		}
	}

	// Если номер телефона не существует или он пустой, продолжаем выполнение без ошибок
	return nil
}

func (uc *Usecase) IsUsernameExist(ctx context.Context, user *entities.User) error {
	if user.Username != "" {
		exist, err := uc.Repo.IsUsernameExist(ctx, user)
		if err != nil {
			// Если произошла ошибка при запросе, возвращаем её
			return fmt.Errorf("failed to check if username exists: %w", err)
		}
		if exist {
			// Если номер телефона уже существует, возвращаем ошибку
			return fmt.Errorf("the username %s is already in use", user.Username)
		}
	}
	// Если номер телефона не существует или он пустой, продолжаем выполнение без ошибок
	return nil

}

func (uc *Usecase) CreateUser(ctx context.Context, user *entities.User) interface{} {
	log.Println(user)
	if err := uc.IsUserExit(
		ctx,
		user,
	); err != nil {
		uc.log.Error("CreateUser error: %v", zap.Error(err))
		return fiber.Map{
			"status":  0,
			"message": "user is exist",
		}
	}
	if err := uc.IsPhoneExist(
		ctx,
		user,
	); err != nil {
		uc.log.Error("CreateUser error: %v", zap.Error(err))
		return fiber.Map{
			"status":  1,
			"message": "number_phone is exist", // Возвращаем текст ошибки
		}
	}

	if err := uc.IsUsernameExist(
		ctx,
		user,
	); err != nil {
		uc.log.Error("CreateUser error: %v", zap.Error(err))
		return fiber.Map{
			"status":  2,
			"message": "username is exist",
		}
	}
	if err := uc.Repo.CreateUser(ctx, user); err != nil {
		uc.log.Error("CreateUser error: %v", zap.Error(err))
		return fiber.Map{
			"status":  3,
			"message": "Cannot create user",
		}
	}
	return fiber.Map{
		"status":  4,
		"message": "Successfully created user",
	}
}
