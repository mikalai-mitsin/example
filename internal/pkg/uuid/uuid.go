package uuid

import "github.com/google/uuid"

type UUID string

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func (uuid UUID) String() string {
	return string(uuid)
}

//go:generate mockgen -build_flags=-mod=mod -destination mock/uuid_mock.go . Generator
type Generator interface {
	NewUUID() UUID
}

type UUIDv4Generator struct{}

func NewUUIDv4Generator() *UUIDv4Generator {
	return &UUIDv4Generator{}
}

func (m *UUIDv4Generator) NewUUID() UUID {
	return UUID(uuid.New().String())
}
