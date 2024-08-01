package jwt

import (
	"errors"
	"os"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	e "github.com/haesuo566/sns_backend/user_service/pkg/utils/erorr"
)

type Util struct {
}

type Token = jwt.Token
type MapClaims = jwt.MapClaims

var once sync.Once
var instance *Util

func New() *Util {
	once.Do(func() {
		instance = &Util{}
	})

	return instance
}

// validation token
func (j *Util) Validation(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, e.Wrap(errors.New(""))
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		e.Wrap(err)
	}

	mapClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, e.Wrap(errors.New(""))
	}

	return mapClaims, nil
}
