package domain

import (
	"context"
	"time"
)

type Character struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	FetchedAt   time.Time `json:"fetchedAt"`
}

type CharacterUsecase interface {
	Fetch(ctx context.Context, page int) ([]int, error)
	GetByID(ctx context.Context, id int) (Character, error)
}

type CharacterReadRepository interface {
        Fetch(ctx context.Context, page int) ([]int, error)
        GetByID(ctx context.Context, id int) (Character, error)
}

type CharacterWriteRepository interface {
        StoreByPage(ctx context.Context, page int) error
        StoreByID(ctx context.Context, id int) error
}
