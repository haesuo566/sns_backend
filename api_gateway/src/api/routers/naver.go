package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth/naver"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/db"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
)

func NaverRouter(router fiber.Router) {
	sql := db.NewDatabase()
	authRepository := auth.NewRepository(sql)

	jwtUtil := jwt.New()
	redisUtil := redis.New()

	googleService := naver.NewService(authRepository, jwtUtil, redisUtil)

	handler := handlers.NewNaverHandler(googleService)

	router.Get("/naver/login", handler.Login)
	router.Get("/naver/callback", handler.Callback)
}
