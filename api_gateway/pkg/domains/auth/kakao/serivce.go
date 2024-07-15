package kakao

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
	"time"

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
	Id          int    `json:"id"`
	ConnectedAt string `json:"connected_at"`
	Properties  struct {
		Nickname string `json:"nickname"`
	} `json:"properties"`
	KakaoAccount struct {
		ProfileNicknameNeedsAgreement bool `json:"profile_nickname_needs_agreement"`
		Profile                       struct {
			Nickname          string `json:"nickname"`
			IsDefaultNickname bool   `json:"is_default_nickname"`
		}
	} `json:"kakao_account"`
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
	request, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, e.Wrap(err)
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json;charset=utf-8")

	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, e.Wrap(err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.Wrap(err)
	}

	userInfo := &userInfo{}
	if err := json.Unmarshal(data, userInfo); err != nil {
		return nil, e.Wrap(err)
	}

	return s.SaveUser(&entities.User{
		Name:     userInfo.Properties.Nickname,
		Email:    "", // 카카오 심사 통과시 얻을 수 있음
		UserTag:  uuid.NewString(),
		Platform: "KAKAO",
	})
}

func (s *service) SaveUser(user *entities.User) (*jwt.AllToken, error) {
	user, err := s.authRepository.Save(user)
	if err != nil {
		return nil, e.Wrap(err)
	}

	jwtToken, err := s.jwtUtil.GenerateJwtToken()
	if err != nil {
		return nil, e.Wrap(err)
	}

	u, err := json.Marshal(user)
	if err != nil {
		return nil, e.Wrap(err)
	}

	if err := s.redisUtil.Set(context.Background(), jwtToken.AccessToken, u, time.Minute*30).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	return jwtToken, nil
}
