package google

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
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var serviceOnce sync.Once
var serviceInstance auth.Service

func NewService(authRepository auth.Repository, jwtUtil jwt.Util, redisUtil redis.Util) auth.Service {
	serviceOnce.Do(func() {
		serviceInstance = &service{
			authRepository,
			jwtUtil,
			redisUtil,
		}
	})

	return serviceInstance
}

func (s *service) GetJwtToken(token *oauth2.Token) (*jwt.AllToken, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken)

	// User Infomation Request
	response, err := http.Get(url)
	if err != nil {
		return nil, e.Wrap(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.Wrap(err)
	}

	userInfo := userInfo{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, e.Wrap(err)
	}

	var randomString string
	if rand, err := uuid.NewRandom(); err != nil {
		return nil, e.Wrap(err)
	} else {
		randomString = strings.ReplaceAll(rand.String(), "-", "")
	}

	return s.SaveUser(&entities.User{
		Name:      randomString,
		Email:     userInfo.Email,
		AccountId: fmt.Sprintf("@%s", randomString),
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
