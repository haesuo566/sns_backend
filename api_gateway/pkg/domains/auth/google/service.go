package google

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
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

	h := sha256.New()
	if _, err := h.Write([]byte(userInfo.Email)); err != nil {
		return nil, e.Wrap(err)
	}

	// 이거 uuid trigger 든 뭐든 처리하셈
	return s.SaveUser(&entities.User{
		Name:     strings.ReplaceAll(uuid.NewString(), "-", ""),
		Email:    hex.EncodeToString(h.Sum(nil)),
		UserTag:  strings.ReplaceAll(uuid.NewString(), "-", ""),
		Platform: "GOOGLE",
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
