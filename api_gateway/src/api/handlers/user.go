package handlers

import (
	"fmt"
	"path/filepath"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/haesuo566/sns_backend/api_gateway/src/pkg/utils/errx"
)

type UserHandler struct {
}

var userSyncInit sync.Once
var userInstance *UserHandler

func NewUserHandler() *UserHandler {
	userSyncInit.Do(func() {
		userInstance = &UserHandler{}
	})
	return userInstance
}

func (u *UserHandler) ChangeUserProfileImage(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return errx.Trace(err)
	}

	// 파일 확장자 확인하는 로직 필요
	// allowedTypes := []string{"image/jpeg", "image/png"}
	// if !fiber.Includes(allowedTypes, file.Header.Get("Content-Type")) {
	// 	// Handle invalid file type
	// 	return errors.New("Invalid file type")
	// }

	path, err := filepath.Abs("../uploads")
	if err != nil {
		return errx.Trace(err)
	}

	// 이름은 user unique 값으로 저장해야 될듯
	destination := fmt.Sprintf("%s/%s", path, file.Filename)
	if err := ctx.SaveFile(file, destination); err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (u *UserHandler) GetUserProfile(ctx *fiber.Ctx) error {
	return nil
}
