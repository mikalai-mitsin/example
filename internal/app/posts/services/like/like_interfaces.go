package services

//go:generate mockgen -source=like_interfaces.go -package=services -destination=like_interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type likeRepository interface {
	Create(context.Context, entities.Like) error
	Get(context.Context, uuid.UUID) (entities.Like, error)
	List(context.Context, entities.LikeFilter) ([]entities.Like, error)
	Count(context.Context, entities.LikeFilter) (uint64, error)
	Update(context.Context, entities.Like) error
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
