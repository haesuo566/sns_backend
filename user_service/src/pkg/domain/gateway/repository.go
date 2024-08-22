package gateway

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sync"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/db"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/errx"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/redis"
)

type Repository struct {
	db    *sql.DB
	redis *redis.Client
}

var repositorySyncInit sync.Once
var repositoryInstance *Repository

func NewRepository() *Repository {
	repositorySyncInit.Do(func() {
		repositoryInstance = &Repository{
			db:    db.New(),
			redis: redis.New(),
		}
	})
	return repositoryInstance
}

// go-cache vs redis -> msa니까 redis로 하긴하는데 성능 문제가 생시면 go-cache로 바꿔야할 듯
func (r *Repository) SaveUser(ctx context.Context, user *entities.User) (*entities.User, error) {
	result, err := db.StartTx(func(tx *sql.Tx) (interface{}, error) {
		// user caching
		if val := r.redis.HGet(ctx, "User", user.Email).Val(); len(val) > 0 {
			var cachedUser entities.User
			if err := json.Unmarshal([]byte(val), &cachedUser); err != nil {
				return nil, errx.Trace(err)
			}
			return &cachedUser, nil
		}

		// timestamp, userTag 넣어줘야 할듯 여기서
		var returningId int = 0
		// 여기서 timestamp까지 같이 받으면 되지 않나???
		if err := tx.QueryRowContext(ctx, "INSERT INTO sns_user (name, email, user_tag, platform) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.UserTag, user.Platform).Scan(&returningId); err != nil {
			return nil, errx.Trace(err)
		}

		if returningId < 0 {
			return nil, errx.Trace(errors.New("fail to insert user"))
		}
		user.Id = returningId // User에 timestamp 넣어서 insert 해야함

		data, err := json.Marshal(user)
		if err != nil {
			return nil, errx.Trace(err)
		}

		if err := r.redis.HSet(ctx, "User", user.Email, string(data)).Err(); err != nil {
			return nil, errx.Trace(err)
		}

		return user, nil
	})

	if err != nil {
		return nil, errx.Trace(err)
	}

	return result.(*entities.User), nil
}

func (r *Repository) UpdateName(ctx context.Context, user *entities.User) error {
	_, err := db.StartTx(func(tx *sql.Tx) (interface{}, error) {
		if _, err := tx.ExecContext(ctx, "UPDATE sns_user SET name = $1 WHERE id = $2", user.Name, user.Id); err != nil {
			return nil, errx.Trace(err)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return nil, errx.Trace(err)
		}

		if err := r.redis.HSet(ctx, "User", user.Email, string(data)).Err(); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return errx.Trace(err)
	}

	return nil
}

func (r *Repository) UpdateTag(ctx context.Context, user *entities.User) error {
	_, err := db.StartTx(func(tx *sql.Tx) (interface{}, error) {
		rows, err := tx.QueryContext(ctx, "SELECT id FROM sns_user WHERE user_tag = $1", user.UserTag)
		if err != nil {
			return nil, errx.Trace(err)
		}
		defer rows.Close()

		if rows.NextResultSet() {
			return nil, errx.Trace(errors.New("duplicated user tag"))
		}

		if _, err := tx.ExecContext(ctx, "UPDATE sns_user SET user_tag = $1 WHERE id = $2", user.UserTag, user.Id); err != nil {
			return nil, errx.Trace(err)
		}

		data, err := json.Marshal(user)
		if err != nil {
			return nil, errx.Trace(err)
		}

		if err := r.redis.HSet(ctx, "User", user.Email, string(data)).Err(); err != nil {
			return nil, err
		}

		return nil, nil
	})

	if err != nil {
		return errx.Trace(err)
	}

	return nil
}
