package services

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type UserService struct {
	userRepository userRepository
	clock          clock
	logger         logger
	uuid           UUIDGenerator
}

func NewUserService(
	userRepository userRepository,
	clock clock,
	logger logger,
	uuid UUIDGenerator,
) *UserService {
	return &UserService{userRepository: userRepository, clock: clock, logger: logger, uuid: uuid}
}

func (u *UserService) Create(
	ctx context.Context,
	create *entities.UserCreate,
) (*entities.User, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	user := &entities.User{
		ID:        u.uuid.NewUUID(),
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
func (u *UserService) Get(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	user, err := u.userRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserService) List(
	ctx context.Context,
	filter *entities.UserFilter,
) ([]*entities.User, uint64, error) {
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

func (u *UserService) Update(
	ctx context.Context,
	update *entities.UserUpdate,
) (*entities.User, error) {
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
func (u *UserService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.userRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
func (u *UserService) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	user, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
