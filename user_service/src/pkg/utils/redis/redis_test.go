package redis

import (
	"context"
	"testing"

	"github.com/joho/godotenv"
)

func TestSet(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Error(err)
	}

	r := New()

	if err := r.Set(context.Background(), "key", "value", 0).Err(); err != nil {
		t.Error(err)
	}
}

func TestGet(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Error(err)
	}

	r := New()

	result := r.Get(context.Background(), "key")
	if err := result.Err(); err != nil {
		t.Fatal(err)
	}

	t.Log(result.Val())
}

func TestDel(t *testing.T) {
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Error(err)
	}

	r := New()

	result, err := r.Del(context.Background(), "key").Result()
	if err != nil {
		t.Error(err)
	}

	t.Log(result)
}

func TestTransaction(t *testing.T) {
	if err := godotenv.Load("../../../../.env"); err != nil {
		t.Error(err)
	}

	rdb := New()
	ctx1 := context.Background()
	ctx2 := context.Background()

	pipe1 := rdb.TxPipeline()
	pipe1.Set(ctx1, "key1", "value1", 0)
	pipe1.Set(ctx1, "key2", "value2", 0)

	pipe2 := rdb.TxPipeline()
	pipe2.Set(ctx2, "key3", "value3", 0)
	pipe2.Set(ctx2, "key4", "value4", 0)

	// pipe2 트랜잭션 실행
	if _, err := pipe2.Exec(ctx2); err != nil {
		t.Fatal(err)
	}

	// pipe1 트랜잭션 실행
	if _, err := pipe1.Exec(ctx1); err != nil {
		t.Fatal(err)
	}
}
