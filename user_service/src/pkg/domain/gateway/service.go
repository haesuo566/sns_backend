package gateway

import (
	"context"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/dto"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/errx"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/redis"
)

type Service struct {
	gatewayRepository *Repository
	redisUtil         *redis.Client
	producer          *kafka.Producer
}

var serviceSyncInit sync.Once
var serviceInstance *Service

func NewService() *Service {
	serviceSyncInit.Do(func() {
		serviceInstance = &Service{
			gatewayRepository: NewRepository(),
			redisUtil:         redis.New(),
			producer:          producer.New(),
		}
	})
	return serviceInstance
}

func (s *Service) SaveUser(jwtTokenInfo dto.JwtTokenInfo) error {
	pipe := s.redisUtil.TxPipeline()

	ctx := context.Background()
	// 로그아웃 확인을 위해 accessToken을 redis에 저장
	if err := pipe.Set(ctx, jwtTokenInfo.AccessId, jwtTokenInfo.User.Email, time.Minute*15).Err(); err != nil {
		return errx.Trace(err)
	}

	// Refresh Token을 도난 당했을때를 대비해 refresh토큰을 rotation해서 저장한 값과 비교함
	if err := pipe.Set(ctx, jwtTokenInfo.RefreshId, jwtTokenInfo.User.Email, time.Hour*24*7).Err(); err != nil {
		return errx.Trace(err)
	}

	_, err := s.gatewayRepository.SaveUser(ctx, jwtTokenInfo.User)
	if err != nil {
		return errx.Trace(err)
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return err
	}

	return nil
}

// 유저 이름 변경
// 길이 체크??
func (s *Service) ChangeUserName(user *entities.User) error {
	ctx := context.Background()
	if err := s.gatewayRepository.UpdateName(ctx, user); err != nil {
		return errx.Trace(err)
	}

	return nil
}

// 유저 아이디(태그) 변경
func (s *Service) ChangeUserTag(user *entities.User) error {
	ctx := context.Background()
	if err := s.gatewayRepository.UpdateTag(ctx, user); err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (s *Service) GetUserProfile() error {
	return nil
}
