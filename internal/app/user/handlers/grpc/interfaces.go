package handlers

//go:generate mockgen -source=interfaces.go -package=handlers -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type userUseCase interface {
	Create(context.Context, *entities.UserCreate) (*entities.User, error)
	Get(context.Context, uuid.UUID) (*entities.User, error)
	List(context.Context, *entities.UserFilter) ([]*entities.User, uint64, error)
	Update(context.Context, *entities.UserUpdate) (*entities.User, error)
	Delete(context.Context, uuid.UUID) error
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
