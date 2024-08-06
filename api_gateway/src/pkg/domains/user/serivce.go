package user

import (
	"encoding/json"
	"strings"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/google/uuid"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/dto"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"golang.org/x/oauth2"
)

type TemplateService interface {
	GetOauthUser(*oauth2.Token) (*entities.User, error)
}

type Service struct {
	TemplateService
	jwtUtil  jwt.Util
	producer *kafka.Producer
}

var mutex sync.Mutex
var instances map[TemplateService]*Service = make(map[TemplateService]*Service)

func NewService(service TemplateService, jwtUtil jwt.Util, producder *kafka.Producer) *Service {
	instance, exist := instances[service]
	if !exist {
		mutex.Lock()
		instance = &Service{
			TemplateService: service,
			jwtUtil:         jwtUtil,
			producer:        producder,
		}
		instances[service] = instance
		mutex.Unlock()
	}

	return instance
}

func (t *Service) GetJwtToken(token *oauth2.Token) (*jwt.AllToken, error) {
	user, err := t.TemplateService.GetOauthUser(token)
	if err != nil {
		return nil, e.Wrap(err)
	}

	return t.SaveUser(user)
}

// template method pattern 주체 메서드
func (t *Service) SaveUser(user *entities.User) (*jwt.AllToken, error) {
	accessId := strings.ReplaceAll(uuid.NewString(), "-", "")
	accessToken, err := t.jwtUtil.GenerateAccessToken(accessId, user.Email)
	if err != nil {
		return nil, e.Wrap(err)
	}

	refreshId := strings.ReplaceAll(uuid.NewString(), "-", "")
	refreshToken, err := t.jwtUtil.GenerateRefreshToken(refreshId, user.Email, accessId)
	if err != nil {
		return nil, e.Wrap(err)
	}

	jwtTokenInfo := dto.JwtTokenInfo{
		User:      user,
		AccessId:  accessId,
		RefreshId: refreshId,
	}

	data, err := json.Marshal(jwtTokenInfo)
	if err != nil {
		return nil, e.Wrap(err)
	}

	correlationId := strings.ReplaceAll(uuid.NewString(), "-", "")
	topic := "gateway"

	// async
	err = t.producer.Produce(&kafka.Message{
		Key:   []byte("SaveUser"),
		Value: data,
		Headers: []kafka.Header{
			{Key: "CorrelationId", Value: []byte(correlationId)},
		},
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
	}, nil)

	if err != nil {
		return nil, e.Wrap(err)
	}

	return &jwt.AllToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
