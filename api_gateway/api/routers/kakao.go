package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/kakao"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

func KakaoRouter(router fiber.Router) {
	sql := db.NewDatabase()
	authRepository := auth.NewRepository(sql)

	jwtUtil := jwt.New()
	redisUtil := redis.New()

	kakaoService := kakao.NewService(authRepository, jwtUtil, redisUtil)

	handler := handlers.NewKakaoHandler(kakaoService)

	router.Get("/kakao/login", handler.Login)
	router.Get("/kakao/callback", handler.Callback)
}
