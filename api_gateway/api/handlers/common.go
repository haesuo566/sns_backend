package handlers

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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
	id := ctx.Locals("id").(string)
	userId := ctx.Locals("user_id").(string)

	result, err := c.redisUtil.Del(context.Background(), id).Result()
	if err != nil {
		return e.Wrap(err)
	}

	// refreshToken이 redis에 없는 경우 ->
	// 1. 훔친 예전 refreshToken을 사용하는 경우
	// 2. refreshToken의 유요기간이 끝나서 redis에 없는 경우
	if result == 0 {
		return ctx.Status(fiber.StatusUnauthorized).SendString("testest")
	}

	accessId := strings.ReplaceAll(uuid.NewString(), "-", "")
	accessToken, err := c.jwtUtil.GenerateAccessToken(accessId, userId)
	if err != nil {
		return e.Wrap(err)
	}

	refreshId := strings.ReplaceAll(uuid.NewString(), "-", "")
	refreshToken, err := c.jwtUtil.GenerateRefreshToken(refreshId, userId)
	if err != nil {
		return e.Wrap(err)
	}

	if err := c.redisUtil.Set(context.Background(), accessId, userId, time.Minute*30).Err(); err != nil {
		return e.Wrap(err)
	}

	if err := c.redisUtil.Set(context.Background(), refreshId, userId, time.Hour*24*7).Err(); err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (c *commonHandler) Logout(ctx *fiber.Ctx) error {
	// token := ctx.Locals("user").(*jwt.Token).Raw

	// if err := c.redisUtil.Del(context.Background(), token).Err(); err != nil {
	// 	return e.Wrap(err)
	// }

	return nil
}
