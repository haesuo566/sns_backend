package db

import (
	"context"
	"database/sql"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

type Database interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
	Begin() (*sql.Tx, error)
}

var once sync.Once
var instance Database = nil

func NewDatabase() Database {
	once.Do(func() {
		url := os.Getenv("DATABASE_URL")

		db, err := sql.Open("postgres", url)
		if err != nil {
			panic(err)
		}

		// db.SetConnMaxIdleTime()
		// db.SetMaxIdleConns()
		// db.SetMaxOpenConns()

		if err := db.Ping(); err != nil {
			panic(err)
		}

		instance = db
	})
	return instance
}
