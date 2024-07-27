package handlers

import (
	"context"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/common"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

type CommonHandler struct {
	commonService *common.Service
	redisUtil     redis.Util
}

var commonOnce sync.Once
var commonInstance *CommonHandler

func NewCommonHandler(commonService *common.Service, redisUtil redis.Util) *CommonHandler {
	commonOnce.Do(func() {
		commonInstance = &CommonHandler{
			commonService,
			redisUtil,
		}
	})
	return commonInstance
}

func (c *CommonHandler) RefreshToken(ctx *fiber.Ctx) error {
	id := ctx.Locals("id").(string)
	emailHash := ctx.Locals("email_hash").(string)

	result, err := c.redisUtil.Del(context.Background(), id).Result()
	if err != nil {
		return e.Wrap(err)
	}

	// refreshToken이 redis에 없는 경우 ->
	// 1. 탈취당한 예전 refreshToken을 사용하는 경우
	// 2. refreshToken의 유요기간이 끝나서 redis에 없는 경우
	if result == 0 {
		return ctx.Status(fiber.StatusUnauthorized).SendString("invalid token") // content-type : plaintext/utf8
	}

	tokenMap, err := c.commonService.RefreshToken(emailHash)
	if err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(tokenMap)
}

func (c *CommonHandler) Logout(ctx *fiber.Ctx) error {
	refreshId := ctx.Locals("id").(string)
	accessId := ctx.Locals("access_id").(string)

	// refreshtoken으로 로그아웃 -> 이걸로 한 이유는 크게 차이가 나지는 않지만 accesstoken으로 보통 많이 ㅆ는데 이때 헤더 크기떄문에
	// redis에서 uuid 삭제
	// access_token, refresh_token 삭제 해야함
	return c.commonService.Logout(refreshId, accessId)
}
