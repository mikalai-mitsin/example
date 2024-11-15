package services

//go:generate mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type userRepository interface {
	Create(context.Context, *entities.User) error
	Get(context.Context, uuid.UUID) (*entities.User, error)
	List(context.Context, *entities.UserFilter) ([]*entities.User, error)
	Count(context.Context, *entities.UserFilter) (uint64, error)
	Update(context.Context, *entities.User) error
	Delete(context.Context, uuid.UUID) error
	GetByEmail(context.Context, string) (*entities.User, error)
}

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

// UUIDGenerator - UUID generator interface
type UUIDGenerator interface {
	NewUUID() uuid.UUID
}
