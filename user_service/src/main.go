package main

import (
	"github.com/haesuo566/sns_backend/user_service/src/events"
	"github.com/haesuo566/sns_backend/user_service/src/events/topics"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/domain/gateway"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/consumer"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/kafka/producer"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/db"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/redis"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	db := db.NewDatabase()
	redis := redis.New()

	gatewayRepository := gateway.NewRepository(db)
	gatewayService := gateway.NewService(gatewayRepository, redis)

	consumer := consumer.New()
	producer := producer.New()
	gatewayTopic := topics.NewGateWayTopic(gatewayService)

	event := events.New(consumer, producer, gatewayTopic)
	event.Execute()
	// c := consumer.New("group")
	// e := events.New(c)

	// e.Test()
}
