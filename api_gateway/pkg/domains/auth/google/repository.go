package google

import (
	"context"
	"sync"

	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
)

type GoogleRepository interface {
	Save(*entities.User) (*entities.User, error)
}

type googleRepository struct {
	db db.Database
}

var repositoryOnce sync.Once
var repositoryInstance GoogleRepository = nil

func NewGoogleRepository(db db.Database) GoogleRepository {
	repositoryOnce.Do(func() {
		repositoryInstance = &googleRepository{
			db,
		}
	})

	return repositoryInstance
}

func (g *googleRepository) Save(user *entities.User) (*entities.User, error) {
	g.db.QueryContext(context.Background(), "select ")

	return nil, nil
}
