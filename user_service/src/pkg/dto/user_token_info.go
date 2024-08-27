package dto

import "github.com/haesuo566/sns_backend/user_service/src/pkg/entities"

type JwtTokenInfo struct {
	User      *entities.User `json:"user"`
	AccessId  string         `json:"accessId"`
	RefreshId string         `json:"refreshId"`
}
