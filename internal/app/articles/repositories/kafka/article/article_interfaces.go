package events

//go:generate mockgen -source=article_interfaces.go -package=events -destination=article_interfaces_mock.go
import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type logger interface {
	log.Logger
}
type producer interface {
	Send(ctx context.Context, msg *kafka.Message) error
}
