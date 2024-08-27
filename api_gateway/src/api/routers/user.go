package routers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/handlers"
)

func UserRouter(app fiber.Router) {
	handler := handlers.NewUserHandler()

	app.Post("/change-profile-image", handler.ChangeUserProfileImage)
	app.Get("/get-profile", handler.GetUserProfile)
}
