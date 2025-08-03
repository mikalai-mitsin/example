package services

//go:generate mockgen -source=article_interfaces.go -package=services -destination=article_interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type articleRepository interface {
	Create(context.Context, entities.Article) error
	Get(context.Context, uuid.UUID) (entities.Article, error)
	List(context.Context, entities.ArticleFilter) ([]entities.Article, error)
	Count(context.Context, entities.ArticleFilter) (uint64, error)
	Update(context.Context, entities.Article) error
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
