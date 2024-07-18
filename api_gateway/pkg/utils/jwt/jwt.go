package jwt

import (
	"os"
	"sync"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

type Util interface {
	GenerateRefreshToken(string, string) (string, error)
	GenerateAccessToken(string, string) (string, error)
}

type util struct {
}

type Token = jwt.Token

type AllToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var once sync.Once
var instance Util

func New() Util {
	once.Do(func() {
		instance = &util{}
	})

	return instance
}

func (u *util) GenerateRefreshToken(id, userId string) (string, error) {
	claims := jwt.MapClaims{
		"id":      id,
		"sub":     "refresh_token",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Hour * 24 * 7).Unix(),
		"user_id": userId,
	}

	refreshToken, err := signingToken(claims)
	if err != nil {
		return "", e.Wrap(err)
	}

	return refreshToken, nil
}

func (u *util) GenerateAccessToken(id, userId string) (string, error) {
	claims := jwt.MapClaims{
		"id":      id,
		"sub":     "access_token",
		"iat":     time.Now().Unix(),
		"exp":     time.Now().Add(time.Minute * 30).Unix(),
		"user_id": userId,
	}

	token, err := signingToken(claims)
	if err != nil {
		return "", e.Wrap(err)
	}

	return token, nil
}

func signingToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GetJwtConfig(redisUtil redis.Util) fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(os.Getenv("JWT_SECRET")),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)
			id := claims["id"].(string)
			userId := claims["user_id"].(string)

			c.Locals("id", id)
			c.Locals("user_id", userId)
			return c.Next()
		},
		// accessToken validation 실패했을떄 401 unauthorized 주는 로직 구현해야 함
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Next()
		},
		// Filter: func(c *fiber.Ctx) bool {
		// 	return false
		// },
	})
}
