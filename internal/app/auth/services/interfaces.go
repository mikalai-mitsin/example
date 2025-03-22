package services

//go:generate mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}
type authRepository interface {
	Create(ctx context.Context, user userEntities.User) (entities.TokenPair, error)
	Validate(ctx context.Context, token entities.Token) error
	RefreshToken(ctx context.Context, token entities.Token) (entities.TokenPair, error)
	GetSubject(ctx context.Context, token entities.Token) (string, error)
}
type userRepository interface {
	GetByEmail(ctx context.Context, email string) (userEntities.User, error)
	Get(ctx context.Context, id uuid.UUID) (userEntities.User, error)
}
