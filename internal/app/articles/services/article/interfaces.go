package services

//go:generate mockgen -source=interfaces.go -package=services -destination=interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type articleRepository interface {
	Create(context.Context, dtx.TX, entities.Article) error
	Get(context.Context, uuid.UUID) (entities.Article, error)
	List(context.Context, entities.ArticleFilter) ([]entities.Article, error)
	Count(context.Context, entities.ArticleFilter) (uint64, error)
	Update(context.Context, dtx.TX, entities.Article) error
	Delete(context.Context, dtx.TX, uuid.UUID) error
}
type articleEventProducer interface {
	Send(context.Context, entities.Article) error
}

// clock - clock interface
type clock interface {
	Now() time.Time
}
type logger interface {
	log.Logger
}
type uuidGenerator interface {
	NewUUID() uuid.UUID
}
