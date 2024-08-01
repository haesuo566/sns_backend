package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth/google"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/db"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
)

func GoogleRouter(router fiber.Router) {
	sql := db.NewDatabase()
	authRepository := auth.NewRepository(sql)

	jwtUtil := jwt.New()
	redisUtil := redis.New()

	googleService := google.NewService(authRepository, jwtUtil, redisUtil)

	handler := handlers.NewGoogleHandler(googleService)

	router.Get("/google/login", handler.Login)
	router.Get("/google/callback", handler.Callback)
}
