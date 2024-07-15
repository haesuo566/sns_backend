package jwt

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
)

type Util interface {
	GenerateJwtToken() (*AllToken, error)
	GenerateAccessToken() (string, error)
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

func (u *util) GenerateJwtToken() (*AllToken, error) {
	AccessToken, err := u.GenerateAccessToken()
	if err != nil {
		return nil, e.Wrap(err)
	}

	claims := jwt.MapClaims{
		"iss": "haesuo",
		"sub": "refresh_token",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	RefreshToken, err := signingToken(claims)
	if err != nil {
		return nil, e.Wrap(err)
	}

	return &AllToken{
		AccessToken,
		RefreshToken,
	}, nil
}

func (u *util) GenerateAccessToken() (string, error) {
	claims := jwt.MapClaims{
		"iss": "haesuo",
		"sub": "access_token",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Minute * 30).Unix(),
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
			token := c.Locals("user").(*jwt.Token).Raw
			data, err := redisUtil.Get(context.Background(), token).Result()
			if err != nil {
				return e.Wrap(err)
			}

			user := &entities.User{}
			if err := json.Unmarshal([]byte(data), user); err != nil {
				return e.Wrap(err)
			}

			c.Locals("sns_user", user)
			return c.Next()
		},
	})
}
