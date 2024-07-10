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
	if err := godotenv.Load("../../../.env"); err != nil {
		t.Error(err)
	}

	r := New()

	result := r.Get(context.Background(), "key")
	if err := result.Err(); err != nil {
		t.Error(err)
	}

	t.Log(result.Val())
}
