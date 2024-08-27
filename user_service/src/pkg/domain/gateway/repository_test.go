package gateway

import (
	"context"
	"testing"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/joho/godotenv"
)

var repository *Repository

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err)
	}

	repository = NewRepository()
}

func TestSaveUser(t *testing.T) {
	tempString := "Test2"
	user := &entities.User{
		Name:     tempString,
		Email:    tempString,
		UserTag:  tempString,
		Platform: "Test",
	}

	if _, err := repository.SaveUser(context.Background(), user); err != nil {
		t.Fatal(err)
	}
}

func TestChangeUserName(t *testing.T) {
	// repository :=
}

func TestChangeUserTag(t *testing.T) {

}
