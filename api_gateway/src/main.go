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
)

// 유저 db에 저장하는거 user service로 옮기는게 나을것 같음
// 그래야 api_gateway의 복잡성이 낮아짐
func main() {
	app := fiber.New()

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
	defer file.Close()

	multiWriter := io.MultiWriter(os.Stdout, file)

	// middlewares
	app.Use(logger.New(logger.Config{
		Output: multiWriter,
	}))

	if err := godotenv.Load("../.env"); err != nil {
		panic(err)
	}

	// oauth2 group
	authRouter := app.Group("/oauth2")

	routers.GoogleRouter(authRouter)
	routers.NaverRouter(authRouter)
	routers.KakaoRouter(authRouter)

	tokenRouter := app.Group("/common")

	routers.CommonRouter(tokenRouter)

	if err := app.Listen(":12121"); err != nil {
		panic(err)
	}
}
