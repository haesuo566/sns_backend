package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth/common"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
)

func CommonRouter(group fiber.Router) {
	jwtUtil := jwt.New()
	redisUtil := redis.New()

	commonService := common.NewService(jwtUtil, redisUtil)

	handler := handlers.NewCommonHandler(commonService, redisUtil)

	jwtConfig := jwt.GetJwtConfig(redisUtil)

	group.Post("/refresh-token", jwtConfig, handler.RefreshToken)
}
