package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
)

func CommonRouter(group fiber.Router) {
	handler := handlers.NewCommonHandler()

	jwtConfig := jwt.GetJwtConfig()

	group.Post("/refresh-token", jwtConfig, handler.RefreshToken)
}
