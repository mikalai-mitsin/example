package uuid

import "github.com/google/uuid"

type UUID string

func NewUUID() UUID {
	return UUID(uuid.New().String())
}

func (uuid UUID) String() string {
	return string(uuid)
}
