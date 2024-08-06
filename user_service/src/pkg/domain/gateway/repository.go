package gateway

import (
	"context"
	"sync"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/db"
	e "github.com/haesuo566/sns_backend/user_service/src/pkg/utils/erorr"
)

type Repository struct {
	db db.Database
}

var repositorySyncInit sync.Once
var repositoryInstance *Repository

func NewRepository(db db.Database) *Repository {
	repositorySyncInit.Do(func() {
		repositoryInstance = &Repository{
			db,
		}
	})
	return repositoryInstance
}

func (r *Repository) Save(user *entities.User) (*entities.User, error) {
	tx, err := r.db.Begin()
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
