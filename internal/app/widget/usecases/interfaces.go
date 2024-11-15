package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type widgetService interface {
	Create(context.Context, *entities.WidgetCreate) (*entities.Widget, error)
	Get(context.Context, uuid.UUID) (*entities.Widget, error)
	List(context.Context, *entities.WidgetFilter) ([]*entities.Widget, uint64, error)
	Update(context.Context, *entities.WidgetUpdate) (*entities.Widget, error)
	Delete(context.Context, uuid.UUID) error
}
type logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}
