package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
)

func TokenRouter(group fiber.Router) {
	jwtUtil := jwt.New()
	handler := handlers.NewTokenHandler(jwtUtil)

	jwtConfig := jwt.GetJwtConfig()

	group.Get("/refresh-token", jwtConfig, handler.RefreshToken)
}
