package services

//go:generate mockgen -source=tag_interfaces.go -package=services -destination=tag_interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type tagRepository interface {
	Create(context.Context, entities.Tag) error
	Get(context.Context, uuid.UUID) (entities.Tag, error)
	List(context.Context, entities.TagFilter) ([]entities.Tag, error)
	Count(context.Context, entities.TagFilter) (uint64, error)
	Update(context.Context, entities.Tag) error
	Delete(context.Context, uuid.UUID) error
}

// clock - clock interface
type clock interface {
	Now() time.Time
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
type uuidGenerator interface {
	NewUUID() uuid.UUID
}
