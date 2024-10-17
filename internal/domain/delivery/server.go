package server

import (
	"backend/config"
	"backend/internal/domain/entities"
	"backend/internal/domain/usecase"
	"context"
	"database/sql"
	"log"
	"strconv"

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

		s.app.Get("/", func(c *fiber.Ctx) error {
			err := c.SendString("And the API is UP!")
			return err
		})
		s.app.Get("/get/product/all_info", s.GetProductAllInfo)
		s.app.Post("/user/create", s.CreateUser)
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

func (s *Server) HelloWorld(c *fiber.Ctx) error {
	return c.SendString("products")
}

func (s *Server) GetProductAllInfo(FCtx *fiber.Ctx) error {
	adIDParam := FCtx.Query("ad_id")
	// Преобразуем ad_id из строки в число (если требуется)
	adID, err := strconv.Atoi(adIDParam)
	if err != nil {
		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ad_id parameter",
		})
	}
	product := &entities.Advertisment{
		ID: uint64(adID),
	}
	if err := s.Usecase.GetProductAllInfo(
		FCtx.Context(),
		product,
	); err != nil {
		return err
	}

	return FCtx.JSON(product)
}

func convertedByCreateUser(FCtx *fiber.Ctx, user *entities.User) error {
	return FCtx.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "User created successfully",
		"user": fiber.Map{
			"id":           user.ID,
			"path_ava":     user.PathAva,
			"username":     user.Username,
			"firstname":    user.Firstname,
			"lastname":     user.Lastname,
			"number_phone": user.NumberPhone,
		},
	})
}
func (s *Server) CreateUser(FCtx *fiber.Ctx) error {
	var user entities.User

	if err := FCtx.BodyParser(&user); err != nil {
		s.logger.Error("Failed to parse body", zap.Error(err))
		return FCtx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to parse body",
		})
	}
	sqlUser := entities.SqlUser{
		ID:          user.ID,
		PathAva:     sql.NullString{String: user.PathAva, Valid: user.PathAva != ""},
		Username:    sql.NullString{String: user.Username, Valid: user.Username != ""},
		Firstname:   sql.NullString{String: user.Firstname, Valid: user.Firstname != ""},
		Lastname:    sql.NullString{String: user.Lastname, Valid: user.Lastname != ""},
		NumberPhone: sql.NullString{String: user.NumberPhone, Valid: user.NumberPhone != ""},
	}
	log.Println(sqlUser)

	if err := s.Usecase.IsPhoneExist(
		FCtx.Context(),
		&sqlUser,
	); err != nil {
		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  1,
			"message": "number phone is exist",
		})
	}
	if err := s.Usecase.IsUsernameExist(
		FCtx.Context(),
		&sqlUser,
	); err != nil {
		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  2,
			"message": "username is exist",
		})
	}
	// Создание пользователя в БД
	if err := s.Usecase.CreateUser(
		FCtx.Context(),
		&sqlUser, // передаем user как указатель
	); err != nil {
		// Обработка ошибки
		s.logger.Error("Error creating user", zap.Error(err))
		return FCtx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Failed to create user",
		})
	}

	// Если все прошло успешно
	return convertedByCreateUser(FCtx, &user)
}
