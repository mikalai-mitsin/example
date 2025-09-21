package entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewMockPost(t *testing.T) Post {
	t.Helper()
	return Post{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		DeletedAt: pointer.Of(faker.New().Time().Time(time.Now())),
		Body:      faker.New().Lorem().Sentence(15),
	}
}
func NewMockPostFilter(t *testing.T) PostFilter {
	t.Helper()
	return PostFilter{
		PageSize:   pointer.Of(faker.New().UInt64()),
		PageNumber: pointer.Of(faker.New().UInt64()),
		Search:     pointer.Of(faker.New().Lorem().Sentence(15)),
		OrderBy: []PostOrdering{
			PostOrderingCreatedAtDESC,
			PostOrderingUpdatedAtASC,
			PostOrderingUpdatedAtDESC,
			PostOrderingDeletedAtASC,
			PostOrderingBodyASC,
			PostOrderingBodyDESC,
			PostOrderingIdASC,
			PostOrderingDeletedAtDESC,
			PostOrderingIdDESC,
			PostOrderingCreatedAtASC,
		},
		IsDeleted: pointer.Of(faker.New().Bool()),
	}
}
func NewMockPostCreate(t *testing.T) PostCreate {
	t.Helper()
	return PostCreate{Body: faker.New().Lorem().Sentence(15)}
}
func NewMockPostUpdate(t *testing.T) PostUpdate {
	t.Helper()
	return PostUpdate{ID: uuid.NewUUID(), Body: pointer.Of(faker.New().Lorem().Sentence(15))}
}
func NewMockPostDelete(t *testing.T) PostDelete {
	t.Helper()
	return PostDelete{ID: uuid.NewUUID()}
}
