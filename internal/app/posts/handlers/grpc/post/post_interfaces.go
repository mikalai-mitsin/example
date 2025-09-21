package handlers

//go:generate mockgen -package=$GOPACKAGE -source=$GOFILE -destination=mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type postUseCase interface {
	Create(context.Context, entities.PostCreate) (entities.Post, error)
	Get(context.Context, uuid.UUID) (entities.Post, error)
	List(context.Context, entities.PostFilter) ([]entities.Post, uint64, error)
	Update(context.Context, entities.PostUpdate) (entities.Post, error)
	Delete(context.Context, entities.PostDelete) (entities.Post, error)
}
type logger interface {
	log.Logger
}
