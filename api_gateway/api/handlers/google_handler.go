package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"os"
	"sync"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/domains/auth/google"
	"github.com/valyala/fasthttp"
	"golang.org/x/oauth2"
)

type GoogleHandler interface {
	Login(fiber.Ctx) error
	Callback(fiber.Ctx) error
}

type googleHandler struct {
	googleService google.GoogleService
}

var googleConfig oauth2.Config = oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_ID"),
	ClientSecret: os.Getenv("GOOGLE_SECRET"),
	RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://accounts.google.com/o/oauth2/auth",
		TokenURL: "https://oauth2.googleapis.com/token",
	},
}

var googleOnce sync.Once
var googleInstance GoogleHandler = nil

func NewGoogleHandler(googleServive google.GoogleService) GoogleHandler {
	googleOnce.Do(func() {
		googleInstance = &googleHandler{
			googleServive,
		}
	})

	return googleInstance
}

func (g *googleHandler) Login(ctx fiber.Ctx) error {
	state := GenerateToken(ctx)
	url := googleConfig.AuthCodeURL(state)
	ctx.Context().Redirect(url, fasthttp.StatusPermanentRedirect)
	return nil
}

func (g *googleHandler) Callback(ctx fiber.Ctx) error {
	state := ctx.Cookies("state")
	if ctx.FormValue("state") != state {
		ctx.Context().Redirect("/google/login", fasthttp.StatusBadRequest)
		return nil
	}

	code := ctx.FormValue("code")
	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.Context().Redirect("/google/login", fasthttp.StatusBadRequest)
		return nil
	}

	g.googleService.Test(token)
	// if err != nil {
	// 	ctx.Context().Redirect("/google/login", fasthttp.StatusBadRequest)
	// 	return exception.WrapError(err.Error())
	// }

	// responseToken, err := jwt.GenerateResponseToken(user)
	// if err != nil {
	// 	return exception.WrapError(err.Error())
	// }

	return nil
}

func GenerateToken(ctx fiber.Ctx) string {
	data := make([]byte, 16)
	rand.Read(data)
	csrfToken := base64.URLEncoding.EncodeToString(data)

	ctx.Cookie(&fiber.Cookie{
		Name:    "state",
		Value:   csrfToken,
		Expires: time.Now().Add(time.Hour * 24),
	})
	return csrfToken
}
