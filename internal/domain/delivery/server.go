package server

import (
	"backend/config"
	"backend/internal/domain/entities"
	"backend/internal/domain/usecase"
	"context"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
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
		
		s.app.Get("/", func(c *fiber.Ctx) error {
			err := c.SendString("And the API is UP!")
			return err
		})
		s.app.Get("/get/product/all_info", s.GetProductAllInfo)
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