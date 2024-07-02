package main

import (
	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/api/routers"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}
}

func main() {
	app := fiber.New()

	// Auth Group
	authRouter := app.Group("/auth")

	routers.GoogleRouter(authRouter)
	routers.NaverRouter(authRouter)

	if err := app.Listen(":12121"); err != nil {
		panic(err)
	}
}
