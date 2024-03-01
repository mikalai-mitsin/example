package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type UserInterceptor struct {
	userUseCase UserUseCase
	logger      log.Logger
	authUseCase AuthUseCase
}

func NewUserInterceptor(
	userUseCase UserUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *UserInterceptor {
	return &UserInterceptor{userUseCase: userUseCase, logger: logger, authUseCase: authUseCase}
}

func (i *UserInterceptor) Create(
	ctx context.Context,
	create *models.UserCreate,
) (*models.User, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDUserCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDUserCreate, create); err != nil {
		return nil, err
	}
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
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDUserList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDUserList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.userUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *UserInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.User, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDUserDetail); err != nil {
		return nil, err
	}
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDUserDetail, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (i *UserInterceptor) Update(
	ctx context.Context,
	update *models.UserUpdate,
) (*models.User, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDUserUpdate); err != nil {
		return nil, err
	}
	user, err := i.userUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDUserUpdate, user); err != nil {
		return nil, err
	}
	updated, err := i.userUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *UserInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDUserDelete); err != nil {
		return err
	}
	user, err := i.userUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDUserDelete, user); err != nil {
		return err
	}
	if err := i.userUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
