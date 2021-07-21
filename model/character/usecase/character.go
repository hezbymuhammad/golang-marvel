package usecase

import (
	"context"
	"time"
        "log"

	"github.com/hezbymuhammad/golang-marvel-demo/domain"
)

type characterUsecase struct {
	characterReadRepo   domain.CharacterReadRepository
	characterWriteRepo  domain.CharacterWriteRepository
	contextTimeout      time.Duration
}

func NewCharacterUsecase(crr domain.CharacterReadRepository, cwr domain.CharacterWriteRepository, timeout time.Duration) domain.CharacterUsecase {
	return &characterUsecase{
		characterReadRepo: crr,
		characterWriteRepo: cwr,
		contextTimeout: timeout,
	}
}

func (cu *characterUsecase) Fetch(c context.Context, page int) ([]int, error) {
        ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

        go func() {
                err := cu.characterWriteRepo.StoreByPage(context.Background(), page)
                if(err != nil) {
                        log.Printf("Error saving page %d: %v", page, err)
                }
        }()

        res, err := cu.characterReadRepo.Fetch(ctx, page)

        if(err != nil) {
                return nil, err
        }

        return res, nil
}

func (cu *characterUsecase) GetByID(c context.Context, id int) (domain.Character, error) {
        ctx, cancel := context.WithTimeout(c, cu.contextTimeout)
	defer cancel()

        go func() {
                err := cu.characterWriteRepo.StoreByID(context.Background(), id)
                if(err != nil) {
                        log.Printf("Error saving id %d: %v", id, err)
                }
        }()

        res, err := cu.characterReadRepo.GetByID(ctx, id)

        if(err != nil) {
                return domain.Character{}, err
        }

        return res, nil
}
