package handlers

import (
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/common"
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
	service := common.NewService(jwtUtil, redisUtil)
	handler := NewCommonHandler(service, redisUtil)

	app.Get("/refresh-token", jwt.GetJwtConfig(redisUtil), handler.RefreshToken)

	req := httptest.NewRequest("GET", "/refresh-token", nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfaWQiOiI0MWRiMjRmMzQ0OTc0NmJjYTUyZDdmOTMwZmNhNjY3OSIsImVtYWlsX2hhc2giOiI0MDUwNzY5YTVmYTRmZjQxOWIwZWViMTg4ZDhjOGEzZTg4OWQwMDM2MTgxNWIyY2Q2YmRmOGE3NDI3YTNmMDY4IiwiZXhwIjoxNzIyNjkyOTkxLCJpYXQiOjE3MjIwODgxOTEsImlkIjoiZjZiMzRkZjY5MDI0NGU0YjhjYmI1ZjI0ZTI0ZjBiYTEiLCJzdWIiOiJyZWZyZXNoX3Rva2VuIn0.33Nvd9FZsTeWak5i0AoSqAwIO1vQIx0rFsyvIGuLmqc")

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

func TestLogout(t *testing.T) {
	app := fiber.New()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Error(err)
	}

	redisUtil := redis.New()
	jwtUtil := jwt.New()
	service := common.NewService(jwtUtil, redisUtil)
	handler := NewCommonHandler(service, redisUtil)

	app.Get("/logout", jwt.GetJwtConfig(redisUtil), handler.Logout)

	req := httptest.NewRequest("GET", "/logout", nil)
	req.Header.Add("Authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJhY2Nlc3NfaWQiOiJiZDgzM2RhNzY2N2I0NDhkYmI5Yzc3NmMxNDQ1MGY0MCIsImVtYWlsX2hhc2giOiI0MDUwNzY5YTVmYTRmZjQxOWIwZWViMTg4ZDhjOGEzZTg4OWQwMDM2MTgxNWIyY2Q2YmRmOGE3NDI3YTNmMDY4IiwiZXhwIjoxNzIyNjkzMTI4LCJpYXQiOjE3MjIwODgzMjgsImlkIjoiNTg4OTgwODI2OTg4NGI0OGJhZDU5NWFjNmE3NGYxODIiLCJzdWIiOiJyZWZyZXNoX3Rva2VuIn0.1dRCMwHcYfZdfVqPhytYFXHJ5zADCmbhZWghl9C8LTE")

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
