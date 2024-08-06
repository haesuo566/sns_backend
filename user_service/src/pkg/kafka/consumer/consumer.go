package consumer

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var syncInit sync.Once
var instance *kafka.Consumer

func New() *kafka.Consumer {
	syncInit.Do(func() {
		var err error
		instance, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost",
			"group.id":          "user_service_group",
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			panic(err)
		}

		topics := []string{"gateway"} // service 분류
		if err := instance.SubscribeTopics(topics, nil); err != nil {
			panic(err)
		}
	})
	return instance
}
