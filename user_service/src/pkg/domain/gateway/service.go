package gateway

import (
	"context"
	"sync"
	"time"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/dto"
	e "github.com/haesuo566/sns_backend/user_service/src/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/redis"
)

type Service struct {
	gatewayRepository *Repository
	redisUtil         redis.Util
}

var serviceSyncInit sync.Once
var serviceInstance *Service

func NewService(gatewayRepository *Repository, redisUtil redis.Util) *Service {
	serviceSyncInit.Do(func() {
		serviceInstance = &Service{
			gatewayRepository,
			redisUtil,
		}
	})
	return serviceInstance
}

// 유저 저장
func (s *Service) SaveUser(jwtTokenInfo dto.JwtTokenInfo) error {
	_, err := s.gatewayRepository.Save(jwtTokenInfo.User)
	if err != nil {
		return e.Wrap(err)
	}

	// 로그아웃 확인을 위해 accessToken을 redis에 저장
	if err := s.redisUtil.Set(context.Background(), jwtTokenInfo.AccessId, jwtTokenInfo.User.Email, time.Minute*15).Err(); err != nil {
		return e.Wrap(err)
	}

	// Refresh Token을 도난 당했을때를 대비해 refresh토큰을 rotation해서 저장한 값과 비교함
	if err := s.redisUtil.Set(context.Background(), jwtTokenInfo.RefreshId, jwtTokenInfo.User.Email, time.Hour*24*7).Err(); err != nil {
		return e.Wrap(err)
	}

	return nil
}

// 유저 이름 변경
func (s *Service) ChangeUserName(value string) error {
	return nil
}

// 유저 아이디(태그) 변경
func (s *Service) ChangeUserTag(value string) error {
	return nil
}
