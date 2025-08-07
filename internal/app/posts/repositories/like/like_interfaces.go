package repositories

//go:generate mockgen -source=like_interfaces.go -package=repositories -destination=like_interfaces_mock.go
import (
	"context"
	"database/sql"

	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}
type database interface {
	ExecContext(ctx context.Context, query string, args ...interface {
	}) (sql.Result, error)
	GetContext(ctx context.Context, dest any, query string, args ...interface {
	}) error
	SelectContext(ctx context.Context, dest any, query string, args ...interface {
	}) error
	QueryRowContext(ctx context.Context, query string, args ...interface {
	}) *sql.Row
}
