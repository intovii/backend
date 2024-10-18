package server

import "github.com/gofiber/fiber/v2"

func (s *Server) initRouter() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		err := c.SendString("And the API is UP!")
		return err
	})
	s.app.Get("/get/advertisment/all_info", s.GetAdvertismentAllInfo)
	s.app.Get("/get/profile/all_info", s.GetProfileUserAllInfo)
	s.app.Get("/get/profile/statistics", s.GetProfileUserStatistics)
}