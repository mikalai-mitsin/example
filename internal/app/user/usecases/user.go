package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserUseCase struct {
	userService userService
	logger      logger
}

func NewUserUseCase(userService userService, logger logger) *UserUseCase {
	return &UserUseCase{userService: userService, logger: logger}
}

func (i *UserUseCase) Create(
	ctx context.Context,
	create *entities.UserCreate,
) (*entities.User, error) {
	user, err := i.userService.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (i *UserUseCase) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := i.userService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserUseCase) List(
	ctx context.Context,
	filter *entities.UserFilter,
) ([]*entities.User, uint64, error) {
	items, count, err := i.userService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}

func (i *UserUseCase) Update(
	ctx context.Context,
	update *entities.UserUpdate,
) (*entities.User, error) {
	updated, err := i.userService.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.userService.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
