package server

import (
	"backend/config"
	"backend/internal/domain/entities"
	"backend/internal/domain/usecase"
	"context"
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