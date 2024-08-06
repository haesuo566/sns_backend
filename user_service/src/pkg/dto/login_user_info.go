package dto

import "github.com/haesuo566/sns_backend/user_service/src/pkg/entities"

type JwtTokenInfo struct {
	User      *entities.User
	AccessId  string
	RefreshId string
}
