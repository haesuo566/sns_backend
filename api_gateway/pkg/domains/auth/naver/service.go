package naver

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

	h := sha256.New()
	if _, err := h.Write([]byte(userInfo.Response.Email)); err != nil {
		return nil, e.Wrap(err)
	}

	return s.SaveUser(&entities.User{
		Name:     userInfo.Response.Name,
		Email:    hex.EncodeToString(h.Sum(nil)),
		UserTag:  uuid.NewString(),
		Platform: "NAVER",
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
