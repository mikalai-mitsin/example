package repositories

//go:generate mockgen -source=tag_interfaces.go -package=repositories -destination=tag_interfaces_mock.go
import (
	"context"
	"database/sql"

	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type logger interface {
	log.Logger
}
type database interface {
	ExecContext(ctx context.Context, query string, args ...interface {
	}) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...interface {
	}) error
	SelectContext(ctx context.Context, dest any, query string, args ...interface {
	}) error
}
