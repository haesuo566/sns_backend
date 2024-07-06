package auth

import (
	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	"golang.org/x/oauth2"
)

type Service interface {
	Test(*oauth2.Token) (*entities.User, error)
}
