package producer

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var once sync.Once
var instance *kafka.Producer

func New() *kafka.Producer {
	once.Do(func() {
		var err error
		instance, err = kafka.NewProducer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost",
		})
		if err != nil {
			panic(err)
		}
	})
	return instance
}
