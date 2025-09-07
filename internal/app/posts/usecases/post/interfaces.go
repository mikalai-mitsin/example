package usecases

//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type postService interface {
	Create(context.Context, dtx.TX, entities.PostCreate) (entities.Post, error)
	Get(context.Context, uuid.UUID) (entities.Post, error)
	List(context.Context, entities.PostFilter) ([]entities.Post, uint64, error)
	Update(context.Context, dtx.TX, entities.PostUpdate) (entities.Post, error)
	Delete(context.Context, dtx.TX, uuid.UUID) error
}
type postEventService interface {
	Created(context.Context, dtx.TX, entities.Post) error
	Updated(context.Context, dtx.TX, entities.Post) error
	Deleted(context.Context, dtx.TX, uuid.UUID) error
}
type logger interface {
	log.Logger
}
type dtxManager interface {
	NewTx() dtx.TX
}
