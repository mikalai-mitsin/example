package services

//go:generate mockgen -source=like_interfaces.go -package=services -destination=like_interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type likeRepository interface {
	Create(context.Context, dtx.TX, entities.Like) error
	Get(context.Context, uuid.UUID) (entities.Like, error)
	List(context.Context, entities.LikeFilter) ([]entities.Like, error)
	Count(context.Context, entities.LikeFilter) (uint64, error)
	Update(context.Context, dtx.TX, entities.Like) error
	Delete(context.Context, dtx.TX, uuid.UUID) error
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
