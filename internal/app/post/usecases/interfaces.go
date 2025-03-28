package usecases

//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type postService interface {
	Create(context.Context, entities.PostCreate) (entities.Post, error)
	Get(context.Context, uuid.UUID) (entities.Post, error)
	List(context.Context, entities.PostFilter) ([]entities.Post, uint64, error)
	Update(context.Context, entities.PostUpdate) (entities.Post, error)
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
