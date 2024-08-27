package common

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/errx"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
)

type Service struct {
	jwt   jwt.Util
	redis redis.Util
}

var once sync.Once
var instance *Service

func NewService() *Service {
	once.Do(func() {
		instance = &Service{
			jwt:   jwt.New(),
			redis: redis.New(),
		}
	})
	return instance
}

func (s *Service) RefreshToken(emailHash string) (fiber.Map, error) {
	accessId := strings.ReplaceAll(uuid.NewString(), "-", "")
	accessToken, err := s.jwt.GenerateAccessToken(accessId, emailHash)
	if err != nil {
		return nil, errx.Trace(err)
	}

	refreshId := strings.ReplaceAll(uuid.NewString(), "-", "")
	refreshToken, err := s.jwt.GenerateRefreshToken(refreshId, emailHash, accessId)
	if err != nil {
		return nil, errx.Trace(err)
	}

	if err := s.redis.Set(context.Background(), accessId, emailHash, time.Minute*30).Err(); err != nil {
		return nil, errx.Trace(err)
	}

	if err := s.redis.Set(context.Background(), refreshId, emailHash, time.Hour*24*7).Err(); err != nil {
		return nil, errx.Trace(err)
	}

	return fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil
}

func (c *Service) Logout(refreshId, accessId string) error {
	// delete refreshToken
	if err := c.redis.Del(context.Background(), refreshId).Err(); err != nil {
		return errx.Trace(err)
	}

	// delete accessToken
	if err := c.redis.Del(context.Background(), accessId).Err(); err != nil {
		return errx.Trace(err)
	}

	return nil
}
