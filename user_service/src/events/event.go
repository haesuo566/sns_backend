package events

import (
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/haesuo566/sns_backend/user_service/src/events/impls"
	"github.com/haesuo566/sns_backend/user_service/src/events/topics"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/consumer"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/worker"
)

type Event struct {
	consumer     *kafka.Consumer
	producer     *kafka.Producer
	gatewayTopic impls.Topic
}

var once sync.Once
var instance *Event

func New() *Event {
	once.Do(func() {
		instance = &Event{
			consumer:     consumer.New(),
			producer:     producer.New(),
			gatewayTopic: topics.NewGateWayTopic(),
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
		// var accessToken string

		// 추후에 header나 metadata같은걸로 caching하면 괜찮을 듯 한데???
		for i := 0; i < len(headers); i++ {
			header := headers[i]
			key := header.Key
			value := header.Value

			switch key {
			case "CorrelationId":
				correlationId = string(value)
			case "AccessToken":
				// accessToken = string(value)
			default:
				continue
			}
		}

		// work job
		w := worker.Job{
			CorrelationId: correlationId,
			Key:           key,
			Value:         value,
		}

		switch topic {
		case "gateway": // gateway service
			w.Task = e.gatewayTopic.ExecuteEvent
		case "":
		default:
			// log + error handling
			// continue
		}

		worker.AddJob(w)
	}
}
