package google

import (
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
	"golang.org/x/oauth2"
)

type service struct {
	authRepository auth.Repository
}

type userInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Picture       string `json:"picture"`
}

var serviceOnce sync.Once
var serviceInstance auth.Service

func NewService(authRepository auth.Repository) auth.Service {
	serviceOnce.Do(func() {
		serviceInstance = &service{
			authRepository,
		}
	})

	return serviceInstance
}

func (g *service) Test(token *oauth2.Token) (*entities.User, error) {
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

	googleUser := &entities.User{
		Name:      randomString,
		Email:     userInfo.Email,
		AccountId: fmt.Sprintf("@%s", randomString),
	}

	user, err := g.authRepository.Save(googleUser)
	if err != nil {
		return nil, e.Wrap(err)
	}

	return user, nil
}
