package common

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

type Service struct {
	jwtUtil   jwt.Util
	redisUtil redis.Util
}

var once sync.Once
var instance *Service

func NewService(jwtUtil jwt.Util, redisUtil redis.Util) *Service {
	once.Do(func() {
		instance = &Service{
			jwtUtil,
			redisUtil,
		}
	})
	return instance
}

func (s *Service) RefreshToken(emailHash string) (fiber.Map, error) {
	accessId := strings.ReplaceAll(uuid.NewString(), "-", "")
	accessToken, err := s.jwtUtil.GenerateAccessToken(accessId, emailHash)
	if err != nil {
		return nil, e.Wrap(err)
	}

	refreshId := strings.ReplaceAll(uuid.NewString(), "-", "")
	refreshToken, err := s.jwtUtil.GenerateRefreshToken(refreshId, emailHash, accessId)
	if err != nil {
		return nil, e.Wrap(err)
	}

	if err := s.redisUtil.Set(context.Background(), accessId, emailHash, time.Minute*30).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	if err := s.redisUtil.Set(context.Background(), refreshId, emailHash, time.Hour*24*7).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	return fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (c *Service) Logout(refreshId, accessId string) error {
	// delete refreshToken
	if err := c.redisUtil.Del(context.Background(), refreshId).Err(); err != nil {
		return e.Wrap(err)
	}

	// delete accessToken
	if err := c.redisUtil.Del(context.Background(), accessId).Err(); err != nil {
		return e.Wrap(err)
	}

	return nil
}
