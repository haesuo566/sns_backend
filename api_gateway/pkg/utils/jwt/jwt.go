package jwt

import (
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
)

type JwtUtil interface {
	GenerateToken(*entities.User) (*AllToken, error)
}

type jwtUtil struct {
}

type AllToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

var once sync.Once
var instance JwtUtil

func New() JwtUtil {
	once.Do(func() {
		instance = &jwtUtil{}
	})

	return instance
}

func (j *jwtUtil) GenerateToken(user *entities.User) (*AllToken, error) {
	accessClaims := jwt.MapClaims{
		"sub":   "access_token",
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(time.Minute * 30).Unix(),
		"email": user.Email,
	}

	AccessToken, err := signingToken(accessClaims)
	if err != nil {
		return nil, e.Wrap(err)
	}

	refreshClaims := jwt.MapClaims{
		"sub": "refresh_token",
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 25 * 7).Unix(),
	}

	RefreshToken, err := signingToken(refreshClaims)
	if err != nil {
		return nil, e.Wrap(err)
	}

	return &AllToken{
		AccessToken,
		RefreshToken,
	}, nil
}

func signingToken(claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
