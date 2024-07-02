package handlers

import (
	"os"
	"sync"

	"golang.org/x/oauth2"
)

type NaverHandler interface {
}

type naverHandler struct {
}

var naverConfig oauth2.Config = oauth2.Config{
	ClientID:     os.Getenv("NAVER_ID"),
	ClientSecret: os.Getenv("NAVER_SECRET"),
	RedirectURL:  os.Getenv("NAVER_REDIRECT_URL"),
	Scopes:       []string{"https://openapi.naver.com/v1/nid/me"},
	Endpoint: oauth2.Endpoint{
		AuthURL:  "https://nid.naver.com/oauth2.0/authorize",
		TokenURL: "https://nid.naver.com/oauth2.0/token",
	},
}

var naverOnce sync.Once
var naverInstance NaverHandler = nil

func NewNaverHandler() NaverHandler {
	naverOnce.Do(func() {
		naverInstance = &naverHandler{}
	})

	return naverInstance
}

func (n *naverHandler) Login() {

}

func (n *naverHandler) Callback() {

}
