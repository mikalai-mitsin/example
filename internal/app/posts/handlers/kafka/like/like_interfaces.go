package handlers

//go:generate mockgen -source=like_interfaces.go -package=handlers -destination=like_interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type likeUseCase interface {
	Create(context.Context, entities.LikeCreate) (entities.Like, error)
	Get(context.Context, uuid.UUID) (entities.Like, error)
	List(context.Context, entities.LikeFilter) ([]entities.Like, uint64, error)
	Update(context.Context, entities.LikeUpdate) (entities.Like, error)
	Delete(context.Context, entities.LikeDelete) (entities.Like, error)
}
type logger interface {
	log.Logger
}
