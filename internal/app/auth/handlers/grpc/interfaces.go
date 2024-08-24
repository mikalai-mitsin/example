package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

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

// AuthInterceptor - domain layer interceptor interface
//
//go:generate mockgen -build_flags=-mod=mod -destination mock/interfaces.go . AuthInterceptor
type AuthInterceptor interface {
	CreateToken(ctx context.Context, login *models.Login) (*models.TokenPair, error)
	RefreshToken(ctx context.Context, refresh models.Token) (*models.TokenPair, error)
}
