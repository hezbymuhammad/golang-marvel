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

	val, err := c.Client.Get(ctx, "marvel-characters-page-"+fmt.Sprint(page)).Result()
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

	val, err := c.Client.Get(ctx, "marvel-character-id-"+fmt.Sprint(id)).Result()
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
