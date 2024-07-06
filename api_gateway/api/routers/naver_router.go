package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/api/auth/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/naver"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
)

func NaverRouter(router fiber.Router) {
	sql := db.NewDatabase()
	authRepository := auth.NewRepository(sql)

	googleService := naver.NewService(authRepository)
	jwtUtil := jwt.New()

	handler := handlers.NewNaverHandler(googleService, jwtUtil)

	router.Get("/naver/login", handler.Login)
	router.Get("/naver/callback", handler.Callback)
}
