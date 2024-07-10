package handlers

import (
	"context"
	"os"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/api/impls"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
	"golang.org/x/oauth2"
)

type googleHandler struct {
	googleService auth.Service
}

var googleConfig oauth2.Config

var googleOnce sync.Once
var googleInstance impls.AuthHandler

func NewGoogleHandler(googleServive auth.Service) impls.AuthHandler {
	googleOnce.Do(func() {
		googleConfig = oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_ID"),
			ClientSecret: os.Getenv("GOOGLE_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://accounts.google.com/o/oauth2/auth",
				TokenURL: "https://oauth2.googleapis.com/token",
			},
		}

		googleInstance = &googleHandler{
			googleServive,
		}
	})

	return googleInstance
}

func (g *googleHandler) Login(ctx *fiber.Ctx) error {
	state := impls.GenerateToken(ctx)
	url := googleConfig.AuthCodeURL(state)
	return ctx.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (g *googleHandler) Callback(ctx *fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		return nil
	}

	code := ctx.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		return e.Wrap(err)
	}

	jwtToken, err := g.googleService.GetJwtToken(token)
	if err != nil {
		return e.Wrap(err)
	}

	return ctx.JSON(jwtToken)
}
