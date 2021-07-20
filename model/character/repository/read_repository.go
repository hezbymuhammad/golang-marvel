package repository

import (
	"context"
	"encoding/json"
	"fmt"

	redis "github.com/go-redis/redis/v8"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
)

type CharacterReadRepository struct {
	Client redis.Cmdable
}

func NewCharacterReadRepository(Conn redis.Cmdable) domain.CharacterReadRepository {
	return &CharacterReadRepository{
		Client: Conn,
	}
}

func (c *CharacterReadRepository) Fetch(ctx context.Context, page int) ([]int, error) {
	var data []int
        key := "marvel-characters-page-"+fmt.Sprint(page)

        isEmpty, err := c.checkRedisKeyEmpty(ctx, key)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
        if isEmpty {
                return nil, domain.ErrCacheKeyEmpty
        }

	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, domain.ErrInternalServerError
	}
	if len(val) == 0 {
		return nil, domain.ErrNotFound
	}

	err = json.Unmarshal([]byte(val), &data)
	if err != nil {
		return nil, domain.ErrInternalServerError
	}

	return data, nil
}

func (c *CharacterReadRepository) GetByID(ctx context.Context, id int) (domain.Character, error) {
	var character domain.Character
        key := "marvel-character-id-"+fmt.Sprint(id)

        isEmpty, err := c.checkRedisKeyEmpty(ctx, key)
	if err != nil {
		return domain.Character{}, domain.ErrInternalServerError
	}
        if isEmpty {
                return domain.Character{}, domain.ErrCacheKeyEmpty
        }

	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return domain.Character{}, domain.ErrInternalServerError
	}
	if len(val) == 0 {
		return domain.Character{}, domain.ErrNotFound
	}

	err = json.Unmarshal([]byte(val), &character)
	if err != nil {
		return domain.Character{}, domain.ErrInternalServerError
	}

	return character, nil
}

func (c *CharacterReadRepository) checkRedisKeyEmpty(ctx context.Context, str string) (bool, error) {
        val, err := c.Client.Exists(ctx, str).Result()
	if err != nil {
		return false, domain.ErrInternalServerError
	}

        return val == 0, nil
}
