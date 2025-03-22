package handlers

//go:generate mockgen -source=interfaces.go -package=handlers -destination=interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type commentUseCase interface {
	Create(context.Context, entities.CommentCreate) (entities.Comment, error)
	Get(context.Context, uuid.UUID) (entities.Comment, error)
	List(context.Context, entities.CommentFilter) ([]entities.Comment, uint64, error)
	Update(context.Context, entities.CommentUpdate) (entities.Comment, error)
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
