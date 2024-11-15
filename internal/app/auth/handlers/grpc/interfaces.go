package handlers

//go:generate mockgen -source=interfaces.go -package=handlers -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

// logger - base logger interface
type logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}

// AuthUseCase - domain layer usecase interface
type AuthUseCase interface {
	CreateToken(ctx context.Context, login *entities.Login) (*entities.TokenPair, error)
	RefreshToken(ctx context.Context, refresh entities.Token) (*entities.TokenPair, error)
}
