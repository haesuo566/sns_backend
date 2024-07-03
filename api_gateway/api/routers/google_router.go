package routers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/api/handlers"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/google"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
)

func GoogleRouter(router fiber.Router) {
	handler := handlers.NewGoogleHandler(
		google.NewGoogleService(
			google.NewGoogleRepository(
				db.NewDatabase(),
			),
		),
	)

	router.Post("/google/login", handler.Login)
}
