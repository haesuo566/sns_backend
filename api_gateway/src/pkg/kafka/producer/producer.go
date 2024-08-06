package producer

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/gofiber/fiber/v2/log"
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

		// Delivery report handler for produced messages
		go func() {
			for e := range instance.Events() {
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						log.Errorf("Delivery failed: %v\n", ev.TopicPartition)
					} else {
						log.Infof("Delivered message to %v\n", ev.TopicPartition)
					}
				}
			}
		}()
	})
	return instance
}
