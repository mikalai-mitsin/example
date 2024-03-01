package interceptors

import (
	"context"

	"github.com/018bf/example/internal/app/session/models"
	userModels "github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type SessionInterceptor struct {
	sessionUseCase SessionUseCase
	logger         log.Logger
	authUseCase    AuthUseCase
}

func NewSessionInterceptor(
	sessionUseCase SessionUseCase,
	logger log.Logger,
	authUseCase AuthUseCase,
) *SessionInterceptor {
	return &SessionInterceptor{
		sessionUseCase: sessionUseCase,
		logger:         logger,
		authUseCase:    authUseCase,
	}
}

func (i *SessionInterceptor) Create(
	ctx context.Context,
	create *models.SessionCreate,
) (*models.Session, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDSessionCreate); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionCreate, create); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, uint64, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDSessionList); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionList, filter); err != nil {
		return nil, 0, err
	}
	items, count, err := i.sessionUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *SessionInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDSessionDetail); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDetail, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) Update(
	ctx context.Context,
	update *models.SessionUpdate,
) (*models.Session, error) {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate); err != nil {
		return nil, err
	}
	session, err := i.sessionUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate, session); err != nil {
		return nil, err
	}
	updated, err := i.sessionUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *SessionInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	requestUser, err := i.authUseCase.GetUser(ctx)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete); err != nil {
		return err
	}
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	if err := i.authUseCase.HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDelete, session); err != nil {
		return err
	}
	if err := i.sessionUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
