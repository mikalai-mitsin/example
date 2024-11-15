package services

//go:generate mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
import (
	"context"
	"time"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type widgetRepository interface {
	Create(context.Context, *entities.Widget) error
	Get(context.Context, uuid.UUID) (*entities.Widget, error)
	List(context.Context, *entities.WidgetFilter) ([]*entities.Widget, error)
	Count(context.Context, *entities.WidgetFilter) (uint64, error)
	Update(context.Context, *entities.Widget) error
	Delete(context.Context, uuid.UUID) error
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
