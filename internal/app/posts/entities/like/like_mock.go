package entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewMockLike(t *testing.T) Like {
	t.Helper()
	return Like{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		DeletedAt: pointer.Of(faker.New().Time().Time(time.Now())),
		PostId:    uuid.NewUUID(),
		Value:     faker.New().Lorem().Sentence(15),
		UserId:    uuid.NewUUID(),
	}
}
func NewMockLikeFilter(t *testing.T) LikeFilter {
	t.Helper()
	return LikeFilter{
		PageSize:   pointer.Of(faker.New().UInt64()),
		PageNumber: pointer.Of(faker.New().UInt64()),
		Search:     pointer.Of(faker.New().Lorem().Sentence(15)),
		OrderBy: []LikeOrdering{
			LikeOrderingIdASC,
			LikeOrderingIdDESC,
			LikeOrderingCreatedAtASC,
			LikeOrderingUpdatedAtASC,
			LikeOrderingDeletedAtDESC,
			LikeOrderingPostIdASC,
			LikeOrderingUserIdASC,
			LikeOrderingUserIdDESC,
			LikeOrderingCreatedAtDESC,
			LikeOrderingUpdatedAtDESC,
			LikeOrderingDeletedAtASC,
			LikeOrderingPostIdDESC,
			LikeOrderingValueASC,
			LikeOrderingValueDESC,
		},
		IsDeleted: pointer.Of(faker.New().Bool()),
	}
}
func NewMockLikeCreate(t *testing.T) LikeCreate {
	t.Helper()
	return LikeCreate{
		PostId: uuid.NewUUID(),
		Value:  faker.New().Lorem().Sentence(15),
		UserId: uuid.NewUUID(),
	}
}
func NewMockLikeUpdate(t *testing.T) LikeUpdate {
	t.Helper()
	return LikeUpdate{
		ID:     uuid.NewUUID(),
		PostId: pointer.Of(uuid.NewUUID()),
		Value:  pointer.Of(faker.New().Lorem().Sentence(15)),
		UserId: pointer.Of(uuid.NewUUID()),
	}
}
func NewMockLikeDelete(t *testing.T) LikeDelete {
	t.Helper()
	return LikeDelete{ID: uuid.NewUUID()}
}
