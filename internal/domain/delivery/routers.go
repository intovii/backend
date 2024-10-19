package server

import "github.com/gofiber/fiber/v2"

func (s *Server) initRouter() {
	s.app.Get("/", func(c *fiber.Ctx) error {
		err := c.SendString("And the API is UP!")
		return err
	})
	s.app.Get("/get/advertisment/all_info", s.GetAdvertismentAllInfo)
	s.app.Post("/user/create", s.CreateUser)
	s.app.Get("/user/get", s.GetUser) // Маршрут для получения пользователя по ID
	//s.app.Get("/user/username", s.GetUserByUsername) // Маршрут для получения пользователя по username

}
