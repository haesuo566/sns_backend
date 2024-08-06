package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/user"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/user/kakao"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
)

func KakaoRouter(router fiber.Router) {
	jwtUtil := jwt.New()
	kakaoService := kakao.NewService()
	producer := producer.New()

	userService := user.NewService(kakaoService, jwtUtil, producer)

	handler := handlers.NewKakaoHandler(userService)

	router.Get("/kakao/login", handler.Login)
	router.Get("/kakao/callback", handler.Callback)
}
