package events

import (
	"strings"
	"sync"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/haesuo566/sns_backend/user_service/events/topics"
	"github.com/haesuo566/sns_backend/user_service/pkg/utils/jwt"
)

type Event struct {
	consumer     *kafka.Consumer
	producer     *kafka.Producer
	jwtUtil      *jwt.Util
	gatewayTopic *topics.GatewayTopic
}

const (
	gateway string = "gateway"
)

var once sync.Once
var instance *Event

func New(
	consumer *kafka.Consumer,
	producer *kafka.Producer,
	jwtUtil *jwt.Util,
	gatewayTopic *topics.GatewayTopic,
) *Event {
	once.Do(func() {
		instance = &Event{
			consumer,
			producer,
			jwtUtil,
			gatewayTopic,
		}
	})
	return instance
}

func (e *Event) Test() error {
	topicSlice := []string{} // service 분류
	if err := e.consumer.SubscribeTopics(topicSlice, nil); err != nil {
		return err
	}

	for {
		msg, err := e.consumer.ReadMessage(time.Second)
		if err != nil {
			// log
		}

		value := string(msg.Value)         // 실제 데이터
		key := string(msg.Key)             // 어떤 메서드를 동작시킬지 구분
		topic := *msg.TopicPartition.Topic // 토픽으로 서비스 구분
		headers := msg.Headers             // jwt 토큰 인증

		isAuthorized := false
		// var jwtClaims jwt.MapClaims

		for i := 0; i < len(headers); i++ {
			header := headers[i]
			key := header.Key
			value := header.Value

			if !strings.EqualFold(key, "Bearer") {
				continue
			}

			if _, err := e.jwtUtil.Validation(string(value)); err != nil {
				// log
			} else {
				// jwtClaims = claims
			}
		}

		if !isAuthorized {
			// 여기서 response를 줘야할거 같은데?? producer로 pulish 해줘야 할 듯
			err := e.producer.Produce(&kafka.Message{
				Value: []byte(""), // 여기는 뭐 쓸게 없네
				Key:   []byte(""), // 여기에 error라는 표시를 주면 될 듯?
				TopicPartition: kafka.TopicPartition{
					Topic:     &topic,
					Partition: kafka.PartitionAny,
				},
			}, nil)

			if err != nil {

			}

			continue
		}

		// 만약 response가 필요없는 동작이면 그냥 goroutine으로 실행시켜도 될 듯
		switch topic {
		case gateway: // gateway service
			e.gatewayTopic.ExecuteEvent(key, value)
		default:
			// log + error handling
		}
	}
}
