package entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewMockTag(t *testing.T) Tag {
	t.Helper()
	return Tag{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		PostId:    uuid.NewUUID(),
		Value:     faker.New().Lorem().Sentence(15),
	}
}
func NewMockTagFilter(t *testing.T) TagFilter {
	t.Helper()
	return TagFilter{
		PageSize:   pointer.Of(faker.New().UInt64()),
		PageNumber: pointer.Of(faker.New().UInt64()),
		Search:     pointer.Of(faker.New().Lorem().Sentence(15)),
		OrderBy: []TagOrdering{
			TagOrderingIdASC,
			TagOrderingIdDESC,
			TagOrderingCreatedAtASC,
			TagOrderingUpdatedAtASC,
			TagOrderingUpdatedAtDESC,
			TagOrderingPostIdDESC,
			TagOrderingValueASC,
			TagOrderingCreatedAtDESC,
			TagOrderingPostIdASC,
			TagOrderingValueDESC,
		},
	}
}
func NewMockTagCreate(t *testing.T) TagCreate {
	t.Helper()
	return TagCreate{PostId: uuid.NewUUID(), Value: faker.New().Lorem().Sentence(15)}
}
func NewMockTagUpdate(t *testing.T) TagUpdate {
	t.Helper()
	return TagUpdate{
		ID:     uuid.NewUUID(),
		PostId: pointer.Of(uuid.NewUUID()),
		Value:  pointer.Of(faker.New().Lorem().Sentence(15)),
	}
}
