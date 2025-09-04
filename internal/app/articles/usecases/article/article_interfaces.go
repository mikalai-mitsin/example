package usecases

//go:generate mockgen -source=article_interfaces.go -package=usecases -destination=article_interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type articleService interface {
	Create(context.Context, dtx.TX, entities.ArticleCreate) (entities.Article, error)
	Get(context.Context, uuid.UUID) (entities.Article, error)
	List(context.Context, entities.ArticleFilter) ([]entities.Article, uint64, error)
	Update(context.Context, dtx.TX, entities.ArticleUpdate) (entities.Article, error)
	Delete(context.Context, dtx.TX, uuid.UUID) error
}
type articleEventProducer interface {
	Created(context.Context, dtx.TX, entities.Article) error
	Updated(context.Context, dtx.TX, entities.Article) error
	Deleted(context.Context, dtx.TX, uuid.UUID) error
}
type logger interface {
	log.Logger
}
type dtxManager interface {
	NewTx() dtx.TX
}
