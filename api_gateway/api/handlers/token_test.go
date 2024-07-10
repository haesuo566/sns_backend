package handlers

import (
	"net/http/httptest"
	"os"
	"testing"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	jwtUtil "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/joho/godotenv"
)

func TestRefreshToken(t *testing.T) {
	app := fiber.New()

	if err := godotenv.Load("../../.env"); err != nil {
		t.Error(err)
	}

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key:    []byte(os.Getenv("JWT_SECRET")),
			JWTAlg: jwtware.HS256,
		},
	}))

	// Protected route
	app.Get("/test", func(c *fiber.Ctx) error {
		user := c.Locals("user").(*jwt.Token)
		claim := user.Claims.(jwt.MapClaims)
		t.Log(claim)
		return nil
	})

	token, _ := jwtUtil.New().GenerateToken("access_token")
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
