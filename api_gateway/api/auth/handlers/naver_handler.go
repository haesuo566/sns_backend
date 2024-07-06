package handlers

import (
	"context"
	"os"
	"sync"

	"github.com/gofiber/fiber/v3"
	authHandler "github.com/haesuo566/sns_backend/api_gateway/api/auth"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/jwt"
	"golang.org/x/oauth2"
)

type naverHandler struct {
	naverSerivce auth.Service
	jwtUtil      jwt.JwtUtil
}

var naverConfig oauth2.Config
var naverOnce sync.Once
var naverInstance authHandler.Handler = nil

func NewNaverHandler(naverSerivce auth.Service, jwtUtil jwt.JwtUtil) authHandler.Handler {
	naverOnce.Do(func() {
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
			jwtUtil,
		}
	})

	return naverInstance
}

func (n *naverHandler) Login(ctx fiber.Ctx) error {
	state := authHandler.GenerateToken(ctx)
	url := naverConfig.AuthCodeURL(state)
	return ctx.Redirect().Status(fiber.StatusTemporaryRedirect).To(url)
}

func (n *naverHandler) Callback(ctx fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		return nil
	}

	code := ctx.FormValue("code")
	token, err := naverConfig.Exchange(context.Background(), code)
	if err != nil {
		return e.Wrap(err)
	}

	user, err := n.naverSerivce.Test(token)
	if err != nil {
		return e.Wrap(err)
	}

	allToken, err := n.jwtUtil.GenerateToken(user)
	if err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(allToken)
}
