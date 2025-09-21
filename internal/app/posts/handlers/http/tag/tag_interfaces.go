package handlers

//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type tagUseCase interface {
	Create(context.Context, entities.TagCreate) (entities.Tag, error)
	Get(context.Context, uuid.UUID) (entities.Tag, error)
	List(context.Context, entities.TagFilter) ([]entities.Tag, uint64, error)
	Update(context.Context, entities.TagUpdate) (entities.Tag, error)
	Delete(context.Context, entities.TagDelete) (entities.Tag, error)
}
type logger interface {
	log.Logger
}
