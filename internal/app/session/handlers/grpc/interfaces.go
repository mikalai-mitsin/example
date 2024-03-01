package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/session/models"
	"github.com/018bf/example/internal/pkg/uuid"
)

// SessionInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . SessionInterceptor
type SessionInterceptor interface {
	Create(ctx context.Context, create *models.SessionCreate) (*models.Session, error)
	List(ctx context.Context, filter *models.SessionFilter) ([]*models.Session, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Session, error)
	Update(ctx context.Context, update *models.SessionUpdate) (*models.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
}
