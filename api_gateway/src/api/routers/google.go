package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/user"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/user/google"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
)

func GoogleRouter(router fiber.Router) {
	jwtUtil := jwt.New()
	googleService := google.NewService()
	producer := producer.New()

	userService := user.NewService(googleService, jwtUtil, producer)

	handler := handlers.NewGoogleHandler(userService)

	router.Get("/google/login", handler.Login)
	router.Get("/google/callback", handler.Callback)
}
