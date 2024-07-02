package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/api/handlers"
)

func GoogleRouter(router fiber.Router) {
	handler := handlers.NewGoogleHandler()

	router.Post("/google/login", handler.Login)

}
