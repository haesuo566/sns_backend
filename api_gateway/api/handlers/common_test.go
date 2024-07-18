package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
	"github.com/joho/godotenv"
)

func TestRefreshToken(t *testing.T) {
	app := fiber.New()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Error(err)
	}

	redisUtil := redis.New()
	jwtUtil := jwt.New()
	handler := NewCommonHandler(jwtUtil, redisUtil)

	app.Get("/test", jwt.GetJwtConfig(redisUtil), handler.RefreshToken)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjE5MTQ3MjEsImlhdCI6MTcyMTMwOTkyMSwiaWQiOiIyYjZiYmRjNzg3ZDE0NmJhYmEzYmVkODUyNGJmZjY1MyIsInN1YiI6InJlZnJlc2hfdG9rZW4iLCJ1c2VyX2lkIjoiNDA1MDc2OWE1ZmE0ZmY0MTliMGVlYjE4OGQ4YzhhM2U4ODlkMDAzNjE4MTViMmNkNmJkZjhhNzQyN2EzZjA2OCJ9.NwWZJAoXXPMhak-kQGe8PfN9Hc2crkVMZIqz795Knuw")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	defer resp.Body.Close()

	// 상태 코드 확인
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Log(string(data))
}
