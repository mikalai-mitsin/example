package usecases

//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type tagService interface {
	Create(context.Context, dtx.TX, entities.TagCreate) (entities.Tag, error)
	Get(context.Context, uuid.UUID) (entities.Tag, error)
	List(context.Context, entities.TagFilter) ([]entities.Tag, uint64, error)
	Update(context.Context, dtx.TX, entities.TagUpdate) (entities.Tag, error)
	Delete(context.Context, dtx.TX, uuid.UUID) error
}
type tagEventService interface {
	Created(context.Context, dtx.TX, entities.Tag) error
	Updated(context.Context, dtx.TX, entities.Tag) error
	Deleted(context.Context, dtx.TX, uuid.UUID) error
}
type logger interface {
	log.Logger
}
type dtxManager interface {
	NewTx() dtx.TX
}
