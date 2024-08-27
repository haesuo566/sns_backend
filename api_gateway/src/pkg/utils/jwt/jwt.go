package jwt

import (
	"context"
	"os"
	"strings"
	"sync"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/errx"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
)

type Util = *util

type util struct {
}

type Token = jwt.Token

type AllToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	accessString  = "access_token"
	refreshString = "refresh_token"
)

const (
	AccessTime  = time.Minute * 30
	RefreshTime = time.Hour * 24 * 7
)

var once sync.Once
var instance Util

func New() Util {
	once.Do(func() {
		instance = &util{}
	})

	return instance
}

func (u *util) GenerateRefreshToken(id, emailHash, accessId string) (string, error) {
	claims := jwt.MapClaims{
		"id":         id,
		"sub":        refreshString,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(RefreshTime).Unix(),
		"email_hash": emailHash,
		"access_id":  accessId,
	}

	refreshToken, err := signingToken(claims)
	if err != nil {
		return "", errx.Trace(err)
	}

	return refreshToken, nil
}

func (u *util) GenerateAccessToken(id, emailHash string) (string, error) {
	claims := jwt.MapClaims{
		"id":         id,
		"sub":        accessString,
		"iat":        time.Now().Unix(),
		"exp":        time.Now().Add(AccessTime).Unix(),
		"email_hash": emailHash,
	}

	token, err := signingToken(claims)
	if err != nil {
		return "", errx.Trace(err)
	}

	return token, nil
}

func signingToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func GetJwtConfig() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(os.Getenv("JWT_SECRET")),
		},
		SuccessHandler: func(c *fiber.Ctx) error {
			token := c.Locals("user").(*jwt.Token)
			claims := token.Claims.(jwt.MapClaims)

			id := claims["id"].(string)                // user uuid
			emailHash := claims["email_hash"].(string) // user email hash
			sub := claims["sub"].(string)              // token subject

			// email_hash가 redis에 있는지 조회하는 로직은 일단 추가하지 않음 -> 필요할 경우 추가 -> condition = sub == access_token
			if strings.EqualFold(sub, accessString) {
				if _, err := redis.New().Get(context.Background(), id).Result(); err != nil {
					return errx.Trace(err)
				}
			} else if strings.EqualFold(sub, refreshString) {
				accessId := claims["access_id"].(string)
				c.Locals("access_id", accessId)
			}

			c.Locals("id", id)
			c.Locals("email_hash", emailHash)
			return c.Next()
		},
		// accessToken validation 실패했을떄 401 unauthorized 주는 로직 구현해야 함
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			return c.Status(fiber.StatusUnauthorized).SendString("invalid token")
		},
		// Filter: func(c *fiber.Ctx) bool {
		// 	return false
		// },
	})
}
