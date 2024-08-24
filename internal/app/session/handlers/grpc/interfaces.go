package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/session/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

// SessionInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interceptor.go . SessionInterceptor
type SessionInterceptor interface {
	Create(ctx context.Context, create *models.SessionCreate) (*models.Session, error)
	List(ctx context.Context, filter *models.SessionFilter) ([]*models.Session, uint64, error)
	Get(ctx context.Context, id uuid.UUID) (*models.Session, error)
	Update(ctx context.Context, update *models.SessionUpdate) (*models.Session, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

// Logger - base logger interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/logger.go . Logger
type Logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}
