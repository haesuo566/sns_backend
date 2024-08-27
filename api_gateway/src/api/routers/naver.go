package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth/naver"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
)

func NaverRouter(router fiber.Router) {
	jwtUtil := jwt.New()
	naverService := naver.NewService()
	producer := producer.New()

	userService := auth.NewService(naverService, jwtUtil, producer)

	handler := handlers.NewNaverHandler(userService)

	router.Get("/naver/login", handler.Login)
	router.Get("/naver/callback", handler.Callback)
}
