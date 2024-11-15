package services

//go:generate mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
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

// AuthRepository - domain layer repository interface
type AuthRepository interface {
	Create(ctx context.Context, user *userEntities.User) (*entities.TokenPair, error)
	Validate(ctx context.Context, token entities.Token) error
	RefreshToken(ctx context.Context, token entities.Token) (*entities.TokenPair, error)
	GetSubject(ctx context.Context, token entities.Token) (string, error)
}

// UserRepository - domain layer repository interface
type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (*userEntities.User, error)
	Get(ctx context.Context, id uuid.UUID) (*userEntities.User, error)
}
