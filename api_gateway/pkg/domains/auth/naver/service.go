package naver

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/redis"
	"golang.org/x/oauth2"
)

type service struct {
	authRepository auth.Repository
	jwtUtil        jwt.Util
	redisUtil      redis.Util
}

type userInfo struct {
	ResultCode string `json:"resultCode"`
	Message    string `json:"message"`
	Response   struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"response"`
}

var once sync.Once
var instance auth.Service

func NewService(authRepository auth.Repository, jwtUtil jwt.Util, redisUtil redis.Util) auth.Service {
	once.Do(func() {
		instance = &service{
			authRepository,
			jwtUtil,
			redisUtil,
		}
	})
	return instance
}

func (s *service) GetJwtToken(token *oauth2.Token) (*jwt.AllToken, error) {
	request, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	if err != nil {
		return nil, err
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	userInfo := userInfo{}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, err
	}

	var randomString string
	if rand, err := uuid.NewRandom(); err != nil {
		return nil, e.Wrap(err)
	} else {
		randomString = strings.ReplaceAll(rand.String(), "-", "")
	}

	return s.SaveUser(&entities.User{
		Name:      userInfo.Response.Name,
		Email:     userInfo.Response.Email,
		AccountId: fmt.Sprintf("@%s", randomString),
		// Provider: user.NAVER,
	})
}

func (s *service) SaveUser(user *entities.User) (*jwt.AllToken, error) {
	user, err := s.authRepository.Save(user)
	if err != nil {
		return nil, e.Wrap(err)
	}

	jwtToken, err := s.jwtUtil.GenerateAllToken()
	if err != nil {
		return nil, e.Wrap(err)
	}

	u, err := json.Marshal(user)
	if err != nil {
		return nil, e.Wrap(err)
	}

	if err := s.redisUtil.Set(context.Background(), jwtToken.RefreshToken, u, 0).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	return jwtToken, nil
}
