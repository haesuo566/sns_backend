package gateway

import (
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load("../../../../.env"); err != nil {
		panic(err)
	}
}

func TestSaveUser(t *testing.T) {
	// db := db.NewDatabase()
	// redis := redis.New()
	// r := NewRepository(db, redis)

	// tempString := "Test2"
	// user := &entities.User{
	// 	Name:     tempString,
	// 	Email:    tempString,
	// 	UserTag:  tempString,
	// 	Platform: "Test",
	// }

	// if _, err := r.SaveUser(user); err != nil {
	// 	t.Fatal(err)
	// }
}
