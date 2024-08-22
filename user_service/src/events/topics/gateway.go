package topics

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/haesuo566/sns_backend/user_service/src/events/impls"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/domain/gateway"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/dto"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/producer"
)

type gatewayTopic struct {
	gatewayService *gateway.Service
	producer       *kafka.Producer
}

var gatewayInit sync.Once
var gatewayInstance impls.Topic

func NewGateWayTopic() impls.Topic {
	gatewayInit.Do(func() {
		gatewayInstance = &gatewayTopic{
			gatewayService: gateway.NewService(),
			producer:       producer.New(),
		}
	})
	return gatewayInstance
}

func (g *gatewayTopic) ExecuteEvent(correlationId, key string, value []byte) error {
	var err error
	var data interface{}

	switch key {
	case "SaveUser":
		if data, err = impls.Marshal[dto.JwtTokenInfo](value); err == nil {
			err = g.gatewayService.SaveUser(data.(dto.JwtTokenInfo))
		}
	case "ChangeUserName":
		if data, err = impls.Marshal[*entities.User](value); err == nil {
			err = g.gatewayService.ChangeUserName(data.(*entities.User))
		}
	case "ChangeUserTag":
		if data, err = impls.Marshal[*entities.User](value); err == nil {
			err = g.gatewayService.ChangeUserTag(data.(*entities.User))
		}
	case "ChangeProfile":
	default:
		// return fmt.Errorf("")
	}

	// error 대해서 log만 찍도록 해야할 듯
	// producer로 error 날리는것 까지 해야하나?
	// 원래는 transaction outbox pattern으로 처리해야함
	if err != nil {
		err := g.producer.Produce(&kafka.Message{}, nil)
		return err
	}

	return nil
}
