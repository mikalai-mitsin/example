package handlers

//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type articleUseCase interface {
	Create(context.Context, entities.ArticleCreate) (entities.Article, error)
	Get(context.Context, uuid.UUID) (entities.Article, error)
	List(context.Context, entities.ArticleFilter) ([]entities.Article, uint64, error)
	Update(context.Context, entities.ArticleUpdate) (entities.Article, error)
	Delete(context.Context, entities.ArticleDelete) (entities.Article, error)
}
type logger interface {
	log.Logger
}
