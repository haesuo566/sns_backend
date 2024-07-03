package db

import (
	"context"
	"database/sql"
	"os"
	"sync"

	_ "github.com/go-sql-driver/mysql"
)

type Database interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
}

var mutex sync.Mutex
var instance Database = nil

func NewDatabase() Database {
	if instance == nil {
		mutex.Lock()
		defer mutex.Unlock()

		url := os.Getenv("DATABASE_URL")

		d, err := sql.Open("mysql", url)
		if err != nil {
			panic(err)
		}

		instance = d
	}

	return instance
}
