package kakao

import (
	"encoding/json"
	"io"
	"net/http"
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
var instance auth.TemplateService

func NewService() auth.TemplateService {
	once.Do(func() {
		instance = &service{}
	})
	return instance
}

func (s *service) GetOauthUser(token *oauth2.Token) (*entities.User, error) {
	request, err := http.NewRequest("GET", "https://kapi.kakao.com/v2/user/me", nil)
	if err != nil {
		return nil, errx.Trace(err)
	}

	request.Header.Set("Authorization", "Bearer "+token.AccessToken)
	request.Header.Set("Content-Type", "application/json;charset=utf-8")

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

	userInfo := &userInfo{}
	if err := json.Unmarshal(data, userInfo); err != nil {
		return nil, errx.Trace(err)
	}

	return &entities.User{
		Name:     userInfo.Properties.Nickname,
		Email:    "", // 카카오 심사 통과시 얻을 수 있음
		UserTag:  uuid.NewString(),
		Platform: "KAKAO",
	}, nil
}
