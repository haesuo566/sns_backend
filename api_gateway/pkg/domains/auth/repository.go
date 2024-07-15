package auth

import (
	"context"
	"sync"

	"github.com/haesuo566/sns_backend/api_gateway/pkg/entities"
	"github.com/haesuo566/sns_backend/api_gateway/pkg/utils/db"
	e "github.com/haesuo566/sns_backend/api_gateway/pkg/utils/erorr"
)

type Repository interface {
	Save(*entities.User) (*entities.User, error)
	// UpdateName(string) (*entities.User, error)
}

type repository struct {
	db db.Database
}

var repositoryOnce sync.Once
var repositoryInstance Repository = nil

func NewRepository(db db.Database) Repository {
	repositoryOnce.Do(func() {
		repositoryInstance = &repository{
			db,
		}
	})

	return repositoryInstance
}

func (a *repository) Save(user *entities.User) (*entities.User, error) {
	tx, err := a.db.Begin()
	if err != nil {
		return nil, e.Wrap(err)
	}

	rows, err := tx.QueryContext(context.Background(), "select id, name, email, created_at, user_tag, platform from sns_user where email = ?", user.Email)
	if err != nil {
		return nil, e.Wrap(err)
	}
	defer rows.Close()

	if rows.Next() {
		user := &entities.User{}
		if err := rows.Scan(&user.Id, &user.Name, &user.Email, &user.CreatedAt, &user.UserTag, &user.Platform); err != nil {
			return nil, e.Wrap(err)
		}

		return user, nil
	}

	result, err := tx.ExecContext(context.Background(), "insert into sns_user(name, email, user_tag, platform) values(?, ?, ?, ?)", user.Name, user.Email, user.UserTag, user.Platform)
	if err != nil {
		return nil, e.Wrap(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, e.Wrap(err)
	}

	rows, err = tx.QueryContext(context.Background(), "select id, name, email, created_at, user_tag, platform from sns_user where id = ?", id)
	if err != nil {
		return nil, e.Wrap(err)
	}
	defer rows.Close()

	newUser := &entities.User{}
	if rows.Next() {
		if err := rows.Scan(&newUser.Id, &newUser.Name, &newUser.Email, &newUser.CreatedAt, &newUser.UserTag, &newUser.Platform); err != nil {
			return nil, e.Wrap(err)
		}
	}

	if err := tx.Commit(); err != nil {
		if err := tx.Rollback(); err != nil {
			return nil, e.Wrap(err)
		}
		return nil, e.Wrap(err)
	}

	return newUser, nil
}

// func (a *authRepository) UpdateName(name string) (*entities.User, error) {
// 	return nil, nil
// }
