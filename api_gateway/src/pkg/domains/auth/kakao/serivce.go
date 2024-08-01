package kakao

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/jwt"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/redis"
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

	accessId := strings.ReplaceAll(uuid.NewString(), "-", "")
	accessToken, err := s.jwtUtil.GenerateAccessToken(accessId, user.Email)
	if err != nil {
		return nil, e.Wrap(err)
	}

	refreshId := strings.ReplaceAll(uuid.NewString(), "-", "")
	refreshToken, err := s.jwtUtil.GenerateRefreshToken(refreshId, user.Email, accessId)
	if err != nil {
		return nil, e.Wrap(err)
	}

	// 로그아웃 확인을 위해 accessToken을 redis에 저장
	if err := s.redisUtil.Set(context.Background(), accessId, user.Email, jwt.AccessTime).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	// Refresh Token을 도난 당했을때를 대비해 refresh토큰을 rotation해서 저장한 값과 비교함
	if err := s.redisUtil.Set(context.Background(), refreshId, user.Email, jwt.RefreshTime).Err(); err != nil {
		return nil, e.Wrap(err)
	}

	return &jwt.AllToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
