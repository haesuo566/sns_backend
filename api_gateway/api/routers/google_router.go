package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/api/auth/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/google"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
)

func GoogleRouter(router fiber.Router) {
	sql := db.NewDatabase()
	authRepository := auth.NewRepository(sql)

	googleService := google.NewService(authRepository)
	jwtUtil := jwt.New()

	handler := handlers.NewGoogleHandler(googleService, jwtUtil)

	router.Get("/google/login", handler.Login)
	router.Get("/google/callback", handler.Callback)
}
