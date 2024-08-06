package events

import (
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/haesuo566/sns_backend/user_service/src/events/impls"
)

type Event struct {
	consumer     *kafka.Consumer
	producer     *kafka.Producer
	gatewayTopic impls.Topic
}

const (
	gateway string = "gateway"
)

var once sync.Once
var instance *Event

func New(consumer *kafka.Consumer, producer *kafka.Producer, gatewayTopic impls.Topic) *Event {
	once.Do(func() {
		instance = &Event{
			consumer,
			producer,
			gatewayTopic,
		}
	})
	return instance
}

func (e *Event) Execute() error {
	for {
		msg, err := e.consumer.ReadMessage(time.Second)
		if err != nil || msg == nil {
			continue
		}

		key := string(msg.Key)             // 어떤 메서드를 동작시킬지 구분
		topic := *msg.TopicPartition.Topic // 토픽으로 서비스 구분
		headers := msg.Headers             // jwt 토큰 인증
		value := msg.Value

		var correlationId string // requestId

		// 추후에 header나 metadata같은걸로 caching하면 괜찮을 듯 한데???
		for i := 0; i < len(headers); i++ {
			header := headers[i]
			key := header.Key
			value := header.Value

			switch key {
			case "CorrelationId":
				correlationId = string(value)
			default:
				continue
			}
		}

		// 만약 response가 필요없는 동작이면 그냥 goroutine으로 실행시켜도 될 듯
		switch topic {
		case gateway: // gateway service
			e.gatewayTopic.ExecuteEvent(correlationId, key, value)
		default:
			// log + error handling
		}
	}
}
