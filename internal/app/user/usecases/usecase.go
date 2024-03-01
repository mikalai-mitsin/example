package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/clock"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type UserUseCase struct {
	userRepository UserRepository
	clock          clock.Clock
	logger         log.Logger
}

func NewUserUseCase(
	userRepository UserRepository,
	clock clock.Clock,
	logger log.Logger,
) *UserUseCase {
	return &UserUseCase{userRepository: userRepository, clock: clock, logger: logger}
}
func (u *UserUseCase) Create(ctx context.Context, create *models.UserCreate) (*models.User, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	user := &models.User{
		ID:        "",
		UpdatedAt: now,
		CreatedAt: now,
		FirstName: create.FirstName,
		LastName:  create.LastName,
		Password:  create.Password,
		Email:     create.Email,
		GroupID:   create.GroupID,
	}
	if err := u.userRepository.Create(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserUseCase) List(
	ctx context.Context,
	filter *models.UserFilter,
) ([]*models.User, uint64, error) {
	user, err := u.userRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.userRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return user, count, nil
}
func (u *UserUseCase) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserUseCase) Update(ctx context.Context, update *models.UserUpdate) (*models.User, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	user, err := u.userRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.FirstName != nil {
			user.FirstName = *update.FirstName
		}
		if update.LastName != nil {
			user.LastName = *update.LastName
		}
		if update.Password != nil {
			user.Password = *update.Password
		}
		if update.Email != nil {
			user.Email = *update.Email
		}
		if update.GroupID != nil {
			user.GroupID = *update.GroupID
		}
	}
	user.UpdatedAt = u.clock.Now().UTC()
	if err := u.userRepository.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}
func (u *UserUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.userRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
func (u *UserUseCase) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
