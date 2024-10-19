package server

import (
	"backend/config"
	"backend/internal/domain/entities"
	"backend/internal/domain/usecase"
	"context"
	"log"

	// "database/sql"
	// "log"
	"strconv"

	"backend/common"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type Server struct {
	logger  *zap.Logger
	cfg     *config.ConfigModel
	app     *fiber.App
	Usecase *usecase.Usecase
}

func NewServer(logger *zap.Logger, cfg *config.ConfigModel, uc *usecase.Usecase) (*Server, error) {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		app:     fiber.New(),
		Usecase: uc,
	}, nil
}

func (s *Server) OnStart(_ context.Context) error {
	go func() {
		s.logger.Debug("fiber app started")
		s.initRouter()
		if err := s.app.Listen(s.cfg.Server.Host + ":" + s.cfg.Server.Port); err != nil {
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
					"text":   common.ErrInvalidParams,
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
					"text":   common.ErrGetInfo,
				},
			})
	}
	return FCtx.JSON(advertisment)
}

func (s *Server) CreateUser(FCtx *fiber.Ctx) error {
	var user entities.User
	if err := FCtx.BodyParser(&user); err != nil {
		s.logger.Error("Failed to parse body", zap.Error(err))
		return FCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse body",
		})
	}
	// Создание пользователя в БД
	resultInterface := s.Usecase.CreateUser(
		FCtx.Context(),
		&user,
	)

	// Приведение результата к типу fiber.Map
	result, ok := resultInterface.(fiber.Map)
	if !ok {
		// Обработка ошибки, если приведение не удалось
		log.Println("Ошибка приведения результата")
		return nil
	}

	// Доступ к полям status и message
	status := result["status"].(int)
	if status != 4 {
		return FCtx.Status(fiber.StatusBadRequest).JSON(resultInterface)
	} else {
		return FCtx.Status(fiber.StatusCreated).JSON(resultInterface)
	}
}

func (s *Server) ConvertedUser(user *entities.User) fiber.Map {
	return fiber.Map{
		"id":           user.ID,
		"path_ava":     user.PathAva,
		"username":     user.Username,
		"firstname":    user.Firstname,
		"lastname":     user.Lastname,
		"number_phone": user.NumberPhone,
		"role": fiber.Map{
			"id":   user.Role.ID,
			"Name": user.Role.Name,
		},
	}
}
func (s *Server) GetUser(FCtx *fiber.Ctx) error {
	var userID int
	var err error
	userIDParam := FCtx.Query("id")
	userNameParam := FCtx.Query("username")
	if userIDParam != "" {
		if userID, err = strconv.Atoi(userIDParam); err != nil {
			s.logger.Error("Invalid user_id parameter", zap.Error(err))
			return FCtx.Status(fiber.StatusBadRequest).JSON(
				fiber.Map{
					"status": common.StatusInvalidParams,
					"text":   common.ErrInvalidParams,
				},
			)
		}
		result := s.GetUserByID(FCtx, userID)
		return result

	}
	if userNameParam != "" {
		result := s.GetUserByUsername(FCtx, userNameParam)
		return result
	}
	return FCtx.Status(fiber.StatusBadRequest).JSON(
		fiber.Map{
			"status": common.StatusInvalidParams,
			"text":   common.ErrInvalidParams,
		},
	)
}
func (s *Server) GetUserByID(FCtx *fiber.Ctx, userID int) error {
	// Получаем userID из параметров
	var err error

	user := &entities.User{
		ID: uint64(userID),
	}
	// Получаем пользователя по ID через Usecase
	resultInterface := s.Usecase.GetUserByID(FCtx.Context(), user)
	result, ok := resultInterface.(fiber.Map)
	if !ok {
		// Обработка ошибки, если приведение не удалось
		s.logger.Error("result conversion error")
		return FCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  common.StatusInvalidParams,
			"message": "result conversion error",
		})
	}
	status := result["status"].(int)
	message := result["message"].(string)
	if status != 2 {
		s.logger.Error("Failed to get user by ID", zap.Error(err))
		return FCtx.Status(fiber.StatusNotFound).JSON(resultInterface)
	}

	// Возвращаем информацию о пользователе
	return FCtx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status":  status,
			"message": message,
			"user":    s.ConvertedUser(user),
		},
	)
}

// Метод для получения пользователя по username
func (s *Server) GetUserByUsername(FCtx *fiber.Ctx, username string) error {
	// Получаем username из параметров
	var err error
	user := &entities.User{
		Username: username,
	}
	// Получаем пользователя по username через Usecase
	resultInterface := s.Usecase.GetUserByUsername(FCtx.Context(), user)
	result, ok := resultInterface.(fiber.Map)
	if !ok {
		// Обработка ошибки, если приведение не удалось
		s.logger.Error("result conversion error")
		return FCtx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  common.StatusInvalidParams,
			"message": "result conversion error",
		})
	}
	status := result["status"].(int)
	message := result["message"].(string)
	if status != 2 {
		s.logger.Error("Failed to get user by username", zap.Error(err))
		return FCtx.Status(fiber.StatusNotFound).JSON(resultInterface)
	}

	// Возвращаем информацию о пользователе
	return FCtx.Status(fiber.StatusOK).JSON(
		fiber.Map{
			"status":  status,
			"message": message,
			"user":    s.ConvertedUser(user),
		},
	)
}
