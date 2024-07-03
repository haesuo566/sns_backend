package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"golang.org/x/oauth2"
)

type GoogleService interface {
	Test(*oauth2.Token) (*entities.User, error)
}

type googleService struct {
	googleRepository GoogleRepository
}

type googleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

var serviceOnce sync.Once
var serviceInstance GoogleService = nil

func NewGoogleService(googleRepository GoogleRepository) GoogleService {
	serviceOnce.Do(func() {
		serviceInstance = &googleService{
			googleRepository,
		}
	})

	return serviceInstance
}

func (g *googleService) Test(token *oauth2.Token) (*entities.User, error) {
	url := fmt.Sprintf("https://www.googleapis.com/oauth2/v2/userinfo?access_token=%s", token.AccessToken)

	// User Infomation Request
	response, err := http.Get(url)
	if err != nil {
		return nil, e.Wrap(err.Error())
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, e.Wrap(err.Error())
	}

	googleUserInfo := googleUserInfo{}
	if err := json.Unmarshal(body, &googleUserInfo); err != nil {
		return nil, e.Wrap(err.Error())
	}

	// user := &entities.User{}

	return nil, nil
}
