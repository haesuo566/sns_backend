package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/haesuo566/sns_backend/user_service/src/events"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/worker"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// logs directory 생성
	if _, err := os.Stat("../logs"); os.IsNotExist(err) {
		if err := os.Mkdir("../logs", os.ModePerm); err != nil {
			panic(err)
		}
	}

	now := time.Now()
	logFileName := fmt.Sprintf("../logs/%s-%s-%s.log", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()))
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	logrus.SetOutput(multiWriter)
	logrus.SetReportCaller(true)
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			logrus.Error(r)
		}
	}()

	// db := db.NewDatabase()
	// redis := redis.New()

	// consumer := consumer.New()
	// producer := producer.New()

	// gatewayRepository := gateway.NewRepository(db, redis)
	// gatewayService := gateway.NewService(gatewayRepository, redis, producer)

	// gatewayTopic := topics.NewGateWayTopic(gatewayService)

	// Kafka WorkerPool Run
	worker.Run()

	event := events.New()
	event.Execute()
	// c := consumer.New("group")
	// e := events.New(c)

	// e.Test()
}
