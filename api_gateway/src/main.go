package main

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/routers"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func init() {
	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	// logs directory 생성
	if _, err := os.Stat("../logs"); os.IsNotExist(err) {
		if err := os.Mkdir("../logs", os.ModePerm); err != nil {
			panic(err)
		}
	}

	// uploads directory 생성
	if _, err := os.Stat("../uploads"); os.IsNotExist(err) {
		if err := os.Mkdir("../uploads", os.ModePerm); err != nil {
			panic(err)
		}
	}

	now := time.Now()
	logFileName := fmt.Sprintf("../logs/%s-%s-%s-service.log", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()))
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		// ForceColors:   true,
	})

	logrus.SetOutput(multiWriter)
	logrus.SetReportCaller(true)
}

// logging하는건 좀 중요하게 수정해야할 듯 -> 나중에 elk stack에 붙이려면 어떻게든 해야함
// log를 service로직이랑 handler들어오는거랑 분리해서 받아야할까?
func main() {
	app := fiber.New()

	app.Static("/uploads", "../uploads")

	now := time.Now()
	logFileName := fmt.Sprintf("../logs/%s-%s-%s-handler.log", strconv.Itoa(now.Year()), strconv.Itoa(int(now.Month())), strconv.Itoa(now.Day()))
	file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	// middlewares
	// log format 추가해야 함
	app.Use(logger.New(logger.Config{
		Output: multiWriter,
	}))

	// oauth2 group
	authRouter := app.Group("/oauth2")

	routers.GoogleRouter(authRouter)
	routers.NaverRouter(authRouter)
	routers.KakaoRouter(authRouter)

	tokenRouter := app.Group("/common")

	routers.CommonRouter(tokenRouter)

	userRouter := app.Group("user")
	routers.UserRouter(userRouter)

	if err := app.Listen(":12121"); err != nil {
		panic(err)
	}
}
