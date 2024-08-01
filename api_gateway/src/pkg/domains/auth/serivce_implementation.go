package auth

import (
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/entities"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"golang.org/x/oauth2"
)

type Service interface {
	GetJwtToken(*oauth2.Token) (*jwt.AllToken, error)
	SaveUser(*entities.User) (*jwt.AllToken, error)
}
