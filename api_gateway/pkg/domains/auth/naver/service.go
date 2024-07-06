package naver

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
	ResultCode string `json:"resultCode"`
	Message    string `json:"message"`
	Response   struct {
		Id    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	} `json:"response"`
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

func (n *service) Test(token *oauth2.Token) (*entities.User, error) {
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

	var randomString string
	if rand, err := uuid.NewRandom(); err != nil {
		return nil, e.Wrap(err)
	} else {
		randomString = strings.ReplaceAll(rand.String(), "-", "")
	}

	naverUser := &entities.User{
		Name:      userInfo.Response.Name,
		Email:     userInfo.Response.Email,
		AccountId: fmt.Sprintf("@%s", randomString),
		// Provider: user.NAVER,
	}

	user, err := n.authRepository.Save(naverUser)
	if err != nil {
		return nil, e.Wrap(err)
	}

	return user, nil
}
