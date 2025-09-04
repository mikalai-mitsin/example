package services

//go:generate mockgen -source=post_interfaces.go -package=services -destination=post_interfaces_mock.go
import (
	"context"
	"time"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type postRepository interface {
	Create(context.Context, dtx.TX, entities.Post) error
	Get(context.Context, uuid.UUID) (entities.Post, error)
	List(context.Context, entities.PostFilter) ([]entities.Post, error)
	Count(context.Context, entities.PostFilter) (uint64, error)
	Update(context.Context, dtx.TX, entities.Post) error
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
