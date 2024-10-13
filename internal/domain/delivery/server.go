package server

import (
	"backend/config"
	"backend/internal/domain/usecase"
	"context"

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
	// _, err := net.Listen("tcp", s.cfg.Server.Host+":"+s.cfg.Server.Port)
	// if err != nil {
	// 	s.logger.Error("failed to listen: !!! ", zap.Error(err))
	// 	return fmt.Errorf("failed to listen:  %w", err)
	// }
	// protos.RegisterContentServer(s.RPC, s)
	// reflection.Register(s.RPC) //по сети теперь видно все методы сети
	go func() {
		s.logger.Debug("fiber app started")
		
		s.app.Get("/", func(c *fiber.Ctx) error {
			err := c.SendString("And the API is UP!")
			return err
		})
		s.app.Get("/get/user", s.HelloWorld)
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
