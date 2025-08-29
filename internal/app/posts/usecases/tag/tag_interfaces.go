package usecases

//go:generate mockgen -source=tag_interfaces.go -package=usecases -destination=tag_interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type tagService interface {
	Create(context.Context, entities.TagCreate) (entities.Tag, error)
	Get(context.Context, uuid.UUID) (entities.Tag, error)
	List(context.Context, entities.TagFilter) ([]entities.Tag, uint64, error)
	Update(context.Context, entities.TagUpdate) (entities.Tag, error)
	Delete(context.Context, uuid.UUID) error
}
type logger interface {
	log.Logger
}
type tagEventProducer interface {
	Created(context.Context, entities.Tag) error
	Updated(context.Context, entities.Tag) error
	Deleted(context.Context, uuid.UUID) error
}
