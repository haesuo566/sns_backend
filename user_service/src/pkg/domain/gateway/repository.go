package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/haesuo566/sns_backend/user_service/src/pkg/entities"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/db"
	e "github.com/haesuo566/sns_backend/user_service/src/pkg/utils/erorr"
	"github.com/haesuo566/sns_backend/user_service/src/pkg/utils/redis"
)

type Repository struct {
	db    db.Database
	redis *redis.Client
}

var repositorySyncInit sync.Once
var repositoryInstance *Repository

func NewRepository() *Repository {
	repositorySyncInit.Do(func() {
		repositoryInstance = &Repository{
			db:    db.NewDatabase(),
			redis: redis.New(),
		}
	})
	return repositoryInstance
}

// go-cache vs redis -> msa니까 redis로 하긴하는데 성능 문제가 생시면 go-cache로 바꿔야할 듯
// early return pattern을 panic하고 defer에서 recover하는 방법으로 좀 간결하게 짜기 가능해짐
func (r *Repository) SaveUser(user *entities.User) (u *entities.User, er error) {
	ctx := context.Background()

	tx, err := r.db.Begin()
	if err != nil {
		return nil, e.Wrap(err)
	}
	// error가 있을 경우 query rollback
	defer func() {
		if er != nil {
			tx.Rollback()
		}
	}()

	// user caching
	if val := r.redis.HGet(ctx, "User", user.Email).Val(); len(val) > 0 {
		var cachedUser entities.User
		if err := json.Unmarshal([]byte(val), &cachedUser); err != nil {
			return nil, e.Wrap(err)
		}
		return &cachedUser, nil
	}

	// timestamp, userTag 넣어줘야 할듯 여기서
	var returningId int = 0
	if err := tx.QueryRowContext(ctx, "INSERT INTO sns_user (name, email, user_tag, platform) VALUES ($1, $2, $3, $4) RETURNING id", user.Name, user.Email, user.UserTag, user.Platform).Scan(&returningId); err != nil {
		return nil, e.Wrap(err)
	}

	if returningId == 0 {
		return nil, e.Wrap(fmt.Errorf("asdasdsasda"))
	}
	user.Id = returningId // User에 timestamp 넣어서 insert 해야함

	if data, err := json.Marshal(user); err != nil {
		return nil, e.Wrap(err)
	} else {
		if err := r.redis.HSet(ctx, "User", user.Email, string(data)).Err(); err != nil {
			return nil, e.Wrap(err)
		}
	}

	if err := tx.Commit(); err != nil {
		if err := r.redis.HDel(ctx, "User", user.Email).Err(); err != nil {
			return nil, e.Wrap(err)
		}

		return nil, e.Wrap(err)
	}

	return user, nil
}

// func (a *authRepository) UpdateName(name string) (*entities.User, error) {
// 	return nil, nil
// }
