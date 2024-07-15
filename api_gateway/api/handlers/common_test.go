package handlers

import (
	"context"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	jwtUtil "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
	"github.com/joho/godotenv"
)

func TestRefreshToken(t *testing.T) {
	app := fiber.New()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Error(err)
	}

	redisUtil := redis.New()

	// Protected route
	app.Get("/test", jwtUtil.GetJwtConfig(redisUtil), func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwtUtil.Token)
		t.Log(user)
		return nil
	})

	mockUser := &entities.User{
		Id:      1,
		Name:    "test",
		UserTag: "asdsasadsad",
	}
	data, err := json.Marshal(mockUser)
	if err != nil {
		t.Error(err)
	}

	token, _ := jwtUtil.New().GenerateAccessToken()
	if err := redisUtil.Set(context.Background(), token, string(data), 0).Err(); err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to perform request: %v", err)
	}

	defer resp.Body.Close()

	// 상태 코드 확인
	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status code %d, got %d", fiber.StatusOK, resp.StatusCode)
	}
}
