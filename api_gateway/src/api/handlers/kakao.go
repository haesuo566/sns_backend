package handlers

import (
	"context"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/impls"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/user"
	e "github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/erorr"
	"golang.org/x/oauth2"
)

type kakaoHandler struct {
	kakaoService *user.Service
}

var kakaoConfig oauth2.Config

var kakaoOnce sync.Once
var kakaoInstance impls.AuthHandler

func NewKakaoHandler(kakaoService *user.Service) impls.AuthHandler {
	kakaoOnce.Do(func() {
		kakaoConfig = oauth2.Config{
			ClientID:     os.Getenv("KAKAO_ID"),
			ClientSecret: os.Getenv("KAKAO_SECRET"),
			RedirectURL:  os.Getenv("KAKAO_REDIRECT_URL"),
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://kauth.kakao.com/oauth/authorize",
				TokenURL: "https://kauth.kakao.com/oauth/token",
			},
		}

		kakaoInstance = &kakaoHandler{
			kakaoService,
		}
	})

	return kakaoInstance
}

func (n *kakaoHandler) Login(ctx *fiber.Ctx) error {
	state := impls.GenerateToken(ctx)
	url := kakaoConfig.AuthCodeURL(state)
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (n *kakaoHandler) Callback(ctx *fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		return nil
	}

	code := ctx.FormValue("code")
	token, err := kakaoConfig.Exchange(context.Background(), code)
	if err != nil {
		return e.Wrap(err)
	}

	jwtToken, err := n.kakaoService.GetJwtToken(token)
	if err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(jwtToken)
}
