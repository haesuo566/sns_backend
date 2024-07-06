package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/haesuo566/sns_backend/api_gateway/api/routers"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New()

	// middlewares
	app.Use(logger.New())

	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// oauth2 group
	authRouter := app.Group("/oauth2")

	routers.GoogleRouter(authRouter)
	routers.NaverRouter(authRouter)

	if err := app.Listen(":12121"); err != nil {
		panic(err)
	}
}

// https://accounts.google.com/o/oauth2/auth?client_id=1014459614066-945esdhcqevf8u9une9i0b7bvofsihld.apps.googleusercontent.com&redirect_uri=http%3A%2F%2Flocalhost%3A12121%2Fauth%2Fgoogle%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=1G-hoVtnjnsrYH4LMW-Yhg%3D%3D
// https://accounts.google.com/o/oauth2/auth?client_id=&redirect_uri=http%3A%2F%2Flocalhost%3A12121%2Fauth%2Fgoogle%2Fcallback&response_type=code&scope=https%3A%2F%2Fwww.googleapis.com%2Fauth%2Fuserinfo.email&state=1O08_gnpv3AhVDGpkXKprg%3D%3D
