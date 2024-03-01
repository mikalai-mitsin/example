package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/session/models"
	"github.com/018bf/example/internal/pkg/clock"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/uuid"
)

type SessionUseCase struct {
	sessionRepository SessionRepository
	clock             clock.Clock
	logger            log.Logger
}

func NewSessionUseCase(
	sessionRepository SessionRepository,
	clock clock.Clock,
	logger log.Logger,
) *SessionUseCase {
	return &SessionUseCase{sessionRepository: sessionRepository, clock: clock, logger: logger}
}

func (u *SessionUseCase) Create(
	ctx context.Context,
	create *models.SessionCreate,
) (*models.Session, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	session := &models.Session{
		ID:          "",
		UpdatedAt:   now,
		CreatedAt:   now,
		Title:       create.Title,
		Description: create.Description,
	}
	if err := u.sessionRepository.Create(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}

func (u *SessionUseCase) List(
	ctx context.Context,
	filter *models.SessionFilter,
) ([]*models.Session, uint64, error) {
	session, err := u.sessionRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.sessionRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return session, count, nil
}
func (u *SessionUseCase) Get(ctx context.Context, id uuid.UUID) (*models.Session, error) {
	session, err := u.sessionRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (u *SessionUseCase) Update(
	ctx context.Context,
	update *models.SessionUpdate,
) (*models.Session, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	session, err := u.sessionRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.Title != nil {
			session.Title = *update.Title
		}
		if update.Description != nil {
			session.Description = *update.Description
		}
	}
	session.UpdatedAt = u.clock.Now().UTC()
	if err := u.sessionRepository.Update(ctx, session); err != nil {
		return nil, err
	}
	return session, nil
}
func (u *SessionUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.sessionRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
