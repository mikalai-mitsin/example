package interceptors

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/models"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserInterceptor struct {
	userUseCase UserUseCase
	logger      Logger
}

func NewUserInterceptor(userUseCase UserUseCase, logger Logger) *UserInterceptor {
	return &UserInterceptor{userUseCase: userUseCase, logger: logger}
}

func (i *UserInterceptor) Create(
	ctx context.Context,
	create *models.UserCreate,
) (*models.User, error) {
	user, err := i.userUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, uint64, error) {
	items, count, err := i.userUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *UserInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) Update(
	ctx context.Context,
	update *models.UserUpdate,
) (*models.User, error) {
	updated, err := i.userUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *UserInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.userUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
