package usecases

//go:generate mockgen -source=like_interfaces.go -package=usecases -destination=like_interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type likeService interface {
	Create(context.Context, entities.LikeCreate) (entities.Like, error)
	Get(context.Context, uuid.UUID) (entities.Like, error)
	List(context.Context, entities.LikeFilter) ([]entities.Like, uint64, error)
	Update(context.Context, entities.LikeUpdate) (entities.Like, error)
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
type likeEventProducer interface {
	Created(context.Context, entities.Like) error
	Updated(context.Context, entities.Like) error
	Deleted(context.Context, uuid.UUID) error
}
