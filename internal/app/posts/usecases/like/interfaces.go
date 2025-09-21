package usecases

//go:generate mockgen -source=interfaces.go -package=usecases -destination=interfaces_mock.go
import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type likeService interface {
	Create(context.Context, dtx.TX, entities.LikeCreate) (entities.Like, error)
	Get(context.Context, uuid.UUID) (entities.Like, error)
	List(context.Context, entities.LikeFilter) ([]entities.Like, uint64, error)
	Update(context.Context, dtx.TX, entities.LikeUpdate) (entities.Like, error)
	Delete(context.Context, dtx.TX, entities.LikeDelete) (entities.Like, error)
}
type likeEventService interface {
	Send(context.Context, dtx.TX, entities.Like) error
}
type logger interface {
	log.Logger
}
type dtxManager interface {
	NewTx() dtx.TX
}
