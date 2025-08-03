package postgres

//go:generate mockgen -source=like_interfaces.go -package=postgres -destination=like_interfaces_mock.go
import "github.com/mikalai-mitsin/example/internal/pkg/log"

type logger interface {
	Debug(msg string, fields ...log.Field)
	Info(msg string, fields ...log.Field)
	Print(msg string, fields ...log.Field)
	Warn(msg string, fields ...log.Field)
	Error(msg string, fields ...log.Field)
	Fatal(msg string, fields ...log.Field)
	Panic(msg string, fields ...log.Field)
}
