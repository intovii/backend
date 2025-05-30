package server

import (
	"backend/config"
	"backend/internal/domain/entities"
	"backend/internal/domain/usecase"
	"context"
	// "database/sql"
	// "log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"backend/common"
)

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	app 	*fiber.App
	Usecase *usecase.Usecase
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		app:	 fiber.New(),
		Usecase: uc,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	go func() {
		s.logger.Debug("fiber app started")
		s.initRouter()
		if err := s.app.Listen(s.cfg.Server.Host+":"+s.cfg.Server.Port); err != nil {
			s.logger.Error("failed to serve: " + err.Error())
		}
	}()
	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.logger.Debug("stop fiber app")
	s.app.Shutdown()
	return nil
}

func (s *Server) GetAdvertismentAllInfo(FCtx *fiber.Ctx) error {
	var adID int
	var err error
	adIDParam := FCtx.Query("ad_id")
    // Преобразуем ad_id из строки в число (если требуется)
    if adID, err = strconv.Atoi(adIDParam); err != nil {
		s.logger.Error("Invalid ad_id parameter", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusInvalidParams,
					"text": common.ErrInvalidParams,
				},
			},
		)
    }
	advertisment := &entities.Advertisment{
		ID: uint64(adID),
	}
	if err = s.Usecase.GetAdvertismentAllInfo(FCtx.Context(), advertisment); err != nil {
		s.logger.Error("Can not get all advertisment info", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusGetInfo,
					"text": common.ErrGetInfo,
				},
        })
	}
    return FCtx.JSON(advertisment)
}

// func convertedByCreateUser(FCtx *fiber.Ctx, user *entities.User) error {
// 	return FCtx.Status(fiber.StatusCreated).JSON(fiber.Map{
// 		"message": "User created successfully",
// 		"user": fiber.Map{
// 			"id":           user.ID,
// 			"path_ava":     user.PathAva,
// 			"username":     user.Username,
// 			"firstname":    user.Firstname,
// 			"lastname":     user.Lastname,
// 			"number_phone": user.NumberPhone,
// 		},
// 	})
// }
// func (s *Server) CreateUser(FCtx *fiber.Ctx) error {
// 	var user entities.User

// 	if err := FCtx.BodyParser(&user); err != nil {
// 		s.logger.Error("Failed to parse body", zap.Error(err))
// 		return FCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
// 			"message": "Failed to parse body",
// 		})
// 	}
// 	sqlUser := entities.SqlUser{
// 		ID:          user.ID,
// 		PathAva:     sql.NullString{String: user.PathAva, Valid: user.PathAva != ""},
// 		Username:    sql.NullString{String: user.Username, Valid: user.Username != ""},
// 		Firstname:   sql.NullString{String: user.Firstname, Valid: user.Firstname != ""},
// 		Lastname:    sql.NullString{String: user.Lastname, Valid: user.Lastname != ""},
// 		NumberPhone: sql.NullString{String: user.NumberPhone, Valid: user.NumberPhone != ""},
// 	}
// 	log.Println(sqlUser)

// 	if err := s.Usecase.IsPhoneExist(
// 		FCtx.Context(),
// 		&sqlUser,
// 	); err != nil {
// 		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  1,
// 			"message": "number phone is exist",
// 		})
// 	}
// 	if err := s.Usecase.IsUsernameExist(
// 		FCtx.Context(),
// 		&sqlUser,
// 	); err != nil {
// 		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"status":  2,
// 			"message": "username is exist",
// 		})
// 	}
// 	// Создание пользователя в БД
// 	if err := s.Usecase.CreateUser(
// 		FCtx.Context(),
// 		&sqlUser, // передаем user как указатель
// 	); err != nil {
// 		// Обработка ошибки
// 		s.logger.Error("Error creating user", zap.Error(err))
// 		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
// 			"message": "Failed to create user",
// 		})
// 	}

// 	// Если все прошло успешно
// 	return convertedByCreateUser(FCtx, &user)
// }

func (s *Server) GetProfileUserAllInfo(FCtx *fiber.Ctx) error {
	var uID int
	var err error
	uIDParam := FCtx.Query("user_id")
    if uID, err = strconv.Atoi(uIDParam); err != nil {
		s.logger.Error("Invalid user_id parameter", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusInvalidParams,
					"text": common.ErrInvalidParams,
				},
			},
		)
    }
	user := &entities.User{
		ID: uint64(uID),
	}
	if err = s.Usecase.GetProfileUserAllInfo(FCtx.Context(), user); err != nil {
		s.logger.Error("Can not get all user info", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusGetInfo,
					"text": common.ErrGetInfo,
				},
        })
	}
    return FCtx.JSON(user)
}

func (s *Server) GetProfileUserStatistics(FCtx *fiber.Ctx) error {
	var uID int
	var err error
	uIDParam := FCtx.Query("user_id")
    if uID, err = strconv.Atoi(uIDParam); err != nil {
		s.logger.Error("Invalid user_id parameter", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusInvalidParams,
					"text": common.ErrInvalidParams,
				},
			},
		)
    }
	var statisticAdsInfo *[]*entities.ProfileStatistic
	if statisticAdsInfo, err = s.Usecase.GetProfileUserStatistics(FCtx.Context(), uint64(uID)); err != nil {
		s.logger.Error("Can not get statistics info", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusGetInfo,
					"text": common.ErrGetInfo,
				},
			},
		)
	}
	return FCtx.JSON(statisticAdsInfo)
}

func (s *Server) GetProfileMyAdvertisments(FCtx *fiber.Ctx) error {
	var uID int
	var err error
	uIDParam := FCtx.Query("user_id")
    if uID, err = strconv.Atoi(uIDParam); err != nil {
		s.logger.Error("Invalid user_id parameter", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusInvalidParams,
					"text": common.ErrInvalidParams,
				},
			},
		)
    }
	var advertisements *[]*entities.MyAdvertisement
	if advertisements, err = s.Usecase.GetProfileMyAdvertisments(FCtx.Context(), uint64(uID)); err != nil {
		s.logger.Error("Can not get info for Profile My Ads", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusGetInfo,
					"text": common.ErrGetInfo,
				},
			},
		)
	}
	return FCtx.JSON(advertisements)
}

func (s *Server) GetProfileReviews(FCtx *fiber.Ctx) error {
	var uID int
	var err error
	uIDParam := FCtx.Query("user_id")
    if uID, err = strconv.Atoi(uIDParam); err != nil {
		s.logger.Error("Invalid user_id parameter", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusInvalidParams,
					"text": common.ErrInvalidParams,
				},
			},
		)
    }
	var reviews *[]*entities.ProfileReview
	if reviews, err = s.Usecase.GetProfileReviews(FCtx.Context(), uint64(uID)); err != nil {
		s.logger.Error("Can not get info for Profile Reviews", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(
			fiber.Map{
				"message": fiber.Map{
					"status": common.StatusGetInfo,
					"text": common.ErrGetInfo,
				},
			},
		)
	}
	return FCtx.JSON(reviews)
}