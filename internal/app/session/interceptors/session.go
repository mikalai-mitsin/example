package interceptors

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/session/models"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type SessionInterceptor struct {
	sessionUseCase SessionUseCase
	logger         Logger
}

func NewSessionInterceptor(sessionUseCase SessionUseCase, logger Logger) *SessionInterceptor {
	return &SessionInterceptor{sessionUseCase: sessionUseCase, logger: logger}
}

func (i *SessionInterceptor) Create(
	ctx context.Context,
	create *models.SessionCreate,
) (*models.Session, error) {
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
	items, count, err := i.sessionUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return items, count, nil
}
func (i *SessionInterceptor) Get(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	session, err := i.sessionUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (i *SessionInterceptor) Update(
	ctx context.Context,
	update *models.SessionUpdate,
) (*models.Session, error) {
	updated, err := i.sessionUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updated, nil
}
func (i *SessionInterceptor) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.sessionUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
