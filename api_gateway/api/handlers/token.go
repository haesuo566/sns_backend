package handlers

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
)

type TokenHandler interface {
	RefreshToken(*fiber.Ctx) error
}

type tokenHandler struct {
	jwtUtil jwt.Util
}

var tokenOnce sync.Once
var tokenInstance TokenHandler

func NewTokenHandler(jwtUtil jwt.Util) TokenHandler {
	tokenOnce.Do(func() {
		tokenInstance = &tokenHandler{
			jwtUtil,
		}
	})
	return tokenInstance
}

func (t *tokenHandler) RefreshToken(ctx *fiber.Ctx) error {

	return nil
}
