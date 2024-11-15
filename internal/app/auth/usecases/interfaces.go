package usecases

//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

// clock - clock interface
type clock interface {
	Now() time.Time
}

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

// AuthService - domain layer usecase interface
type AuthService interface {
	CreateToken(ctx context.Context, login *entities.Login) (*entities.TokenPair, error)
	CreateTokenByUser(ctx context.Context, user *userEntities.User) (*entities.TokenPair, error)
	RefreshToken(ctx context.Context, refresh entities.Token) (*entities.TokenPair, error)
	ValidateToken(ctx context.Context, access entities.Token) error
	Auth(ctx context.Context, access entities.Token) (*userEntities.User, error)
}
