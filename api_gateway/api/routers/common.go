package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

func CommonRouter(group fiber.Router) {
	jwtUtil := jwt.New()
	redisUtil := redis.New()

	handler := handlers.NewCommonHandler(jwtUtil, redisUtil)

	jwtConfig := jwt.GetJwtConfig(redisUtil)

	group.Get("/refresh-token", jwtConfig, handler.RefreshToken)
}
