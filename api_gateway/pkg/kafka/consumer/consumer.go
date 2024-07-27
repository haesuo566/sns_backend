package consumer

import (
	"sync"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

var instanceMap map[string]*kafka.Consumer
var mutex sync.Mutex

func New(group string) *kafka.Consumer {
	mutex.Lock()
	defer mutex.Unlock()

	c, exists := instanceMap[group]
	if !exists {
		var err error
		c, err = kafka.NewConsumer(&kafka.ConfigMap{
			"bootstrap.servers": "localhost",
			"group.id":          group,
			"auto.offset.reset": "earliest",
		})

		if err != nil {
			panic(err)
		}
	}
	return c
}
