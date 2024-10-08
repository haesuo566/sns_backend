package handlers

import (
	"context"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/api/impls"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/domains/auth"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/errx"
	"golang.org/x/oauth2"
)

type naverHandler struct {
	naverSerivce *auth.Service
}

var naverConfig oauth2.Config

var naverSyncInit sync.Once
var naverInstance impls.AuthHandler

func NewNaverHandler(naverSerivce *auth.Service) impls.AuthHandler {
	naverSyncInit.Do(func() {
		naverConfig = oauth2.Config{
			ClientID:     os.Getenv("NAVER_ID"),
			ClientSecret: os.Getenv("NAVER_SECRET"),
			RedirectURL:  os.Getenv("NAVER_REDIRECT_URL"),
			Scopes:       []string{"https://openapi.naver.com/v1/nid/me"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
				TokenURL: "https://nid.naver.com/oauth2.0/token",
			},
		}

		naverInstance = &naverHandler{
			naverSerivce,
		}
	})

	return naverInstance
}

func (n *naverHandler) Login(ctx *fiber.Ctx) error {
	state := impls.GenerateToken(ctx)
	url := naverConfig.AuthCodeURL(state)
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (n *naverHandler) Callback(ctx *fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		return nil
	}

	code := ctx.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		return errx.Trace(err)
	}

	jwtToken, err := n.naverSerivce.GetJwtToken(token)
	if err != nil {
		return errx.Trace(err)
	}

	return ctx.JSON(jwtToken)
}
