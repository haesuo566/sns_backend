package naver

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/entities"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/errx"
	"golang.org/x/oauth2"
)

type service struct {
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
var instance auth.TemplateService

func NewService() auth.TemplateService {
	once.Do(func() {
		instance = &service{}
	})
	return instance
}

func (s *service) GetOauthUser(token *oauth2.Token) (*entities.User, error) {
	request, err := http.NewRequest("GET", "https://openapi.naver.com/v1/nid/me", nil)
	if err != nil {
		return nil, errx.Trace(err)
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	client := http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, errx.Trace(err)
	}
	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errx.Trace(err)
	}

	userInfo := userInfo{}
	if err := json.Unmarshal(data, &userInfo); err != nil {
		return nil, errx.Trace(err)
	}

	h := sha256.New()
	if _, err := h.Write([]byte(userInfo.Response.Email)); err != nil {
		return nil, errx.Trace(err)
	}

	return &entities.User{
		Name:     userInfo.Response.Name,
		Email:    hex.EncodeToString(h.Sum(nil)),
		UserTag:  strings.ReplaceAll(uuid.NewString(), "-", ""),
		Platform: "NAVER",
	}, nil
}
