package db

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/errx"
	_ "github.com/lib/pq"
)

var once sync.Once
var instance *sql.DB = nil

func New() *sql.DB {
	once.Do(func() {
		url := os.Getenv("DATABASE_URL")

		var err error
		instance, err = sql.Open("postgres", url)
		if err != nil {
			panic(err)
		}

		// instance.SetConnMaxIdleTime()
		// instance.SetMaxIdleConns()
		// instance.SetMaxOpenConns()

		if err := instance.Ping(); err != nil {
			panic(err)
		}
	})
	return instance
}

// Redis, Postgresql를 트랜잭션으로 사용함
func StartTx(method func(tx *sql.Tx) (interface{}, error)) (i interface{}, e error) {
	if instance == nil {
		return nil, errx.Trace(fmt.Errorf(""))
	}

	tx, err := instance.Begin()
	if err != nil {
		return nil, errx.Trace(err)
	}

	defer func() {
		if e != nil {
			tx.Rollback()
		}
	}()

	result, err := method(tx)
	if err != nil {
		return nil, errx.Trace(err)
	}

	// 에러가 날 일이 있나 싶긴하네
	if err := tx.Commit(); err != nil {
		return nil, errx.Trace(err)
	}

	return result, nil
}
