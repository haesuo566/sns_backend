package jwt

import (
	"os"
	"sync"
	"time"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
)

type Util interface {
	GenerateAllToken() (*AllToken, error)
	GenerateToken(string) (string, error)
}

type util struct {
}

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

func (u *util) GenerateAllToken() (*AllToken, error) {
	AccessToken, err := u.GenerateToken("access_token")
	if err != nil {
		return nil, e.Wrap(err)
	}

	RefreshToken, err := u.GenerateToken("refresh_token")
	if err != nil {
		return nil, e.Wrap(err)
	}

	return &AllToken{
		AccessToken,
		RefreshToken,
	}, nil
}

func (u *util) GenerateToken(sub string) (string, error) {
	claims := jwt.MapClaims{
		"iss": "haesuo",
		"sub": sub,
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

func GetJwtConfig() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			JWTAlg: jwtware.HS256,
			Key:    []byte(os.Getenv("JWT_SECRET")),
		},
	})
}
