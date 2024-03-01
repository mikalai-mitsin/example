package usecases

import (
	"context"

	"github.com/018bf/example/internal/app/session/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// SessionRepository - domain layer repository interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . SessionRepository
type SessionRepository interface {
	Create(ctx context.Context, session *models.Session) error
	List(ctx context.Context, filter *models.SessionFilter) ([]*models.Session, error)
	Count(ctx context.Context, filter *models.SessionFilter) (uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Session, error)
	Update(ctx context.Context, session *models.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
}
