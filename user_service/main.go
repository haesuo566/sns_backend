package main

import (
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
	}

	// c := consumer.New("group")
	// e := events.New(c)

	// e.Test()
}
