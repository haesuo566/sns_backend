package authHandler

import (
	"crypto/rand"
	"encoding/base64"
	"time"

	"github.com/gofiber/fiber/v3"
)

type Handler interface {
	Login(ctx fiber.Ctx) error
	Callback(ctx fiber.Ctx) error
}

func GenerateToken(ctx fiber.Ctx) string {
	data := make([]byte, 16)
	rand.Read(data)
	csrfToken := base64.URLEncoding.EncodeToString(data)

	ctx.Cookie(&fiber.Cookie{
		Name:    "state",
		Value:   csrfToken,
		Expires: time.Now().Add(time.Hour * 24),
	})
	return csrfToken
}
