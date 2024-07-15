package handlers

import (
	"context"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

type CommonHandler interface {
	RefreshToken(*fiber.Ctx) error
	Logout(*fiber.Ctx) error
}

type commonHandler struct {
	jwtUtil   jwt.Util
	redisUtil redis.Util
}

var commonOnce sync.Once
var commonInstance CommonHandler

func NewCommonHandler(jwtUtil jwt.Util, redisUtil redis.Util) CommonHandler {
	commonOnce.Do(func() {
		commonInstance = &commonHandler{
			jwtUtil,
			redisUtil,
		}
	})
	return commonInstance
}

func (c *commonHandler) RefreshToken(ctx *fiber.Ctx) error {
	accessToken, err := c.jwtUtil.GenerateAccessToken()
	if err != nil {
		return e.Wrap(err)
	}

	if err := c.redisUtil.Set(context.Background(), accessToken, nil, time.Minute*30).Err(); err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(fiber.Map{
		"access_token": accessToken,
	})
}

func (c *commonHandler) Logout(ctx *fiber.Ctx) error {
	token := ctx.Locals("user").(*jwt.Token).Raw

	if err := c.redisUtil.Del(context.Background(), token).Err(); err != nil {
		return e.Wrap(err)
	}

	return nil
}
