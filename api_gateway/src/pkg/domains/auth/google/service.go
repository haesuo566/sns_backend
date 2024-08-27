package google

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
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
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var syncInit sync.Once
var instance auth.TemplateService

func NewService() auth.TemplateService {
	syncInit.Do(func() {
		instance = &service{}
	})

	return instance
}

func (s *service) GetOauthUser(token *oauth2.Token) (*entities.User, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken)

	// User Infomation Request
	response, err := http.Get(url)
	if err != nil {
		return nil, errx.Trace(err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, errx.Trace(err)
	}

	userInfo := userInfo{}
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errx.Trace(err)
	}

	h := sha256.New()
	if _, err := h.Write([]byte(userInfo.Email)); err != nil {
		return nil, errx.Trace(err)
	}

	return &entities.User{
		Name:     strings.ReplaceAll(uuid.NewString(), "-", ""),
		Email:    hex.EncodeToString(h.Sum(nil)),
		UserTag:  strings.ReplaceAll(uuid.NewString(), "-", ""),
		Platform: "GOOGLE",
	}, nil
}
