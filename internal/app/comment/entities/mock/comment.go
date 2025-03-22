package mock_entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewComment(t *testing.T) entities.Comment {
	t.Helper()
	return entities.Comment{
		ID:        uuid.NewUUID(),
		CreatedAt: faker.New().Time().Time(time.Now()),
		UpdatedAt: faker.New().Time().Time(time.Now()),
		Text:      faker.New().Lorem().Sentence(15),
		AuthorId:  uuid.NewUUID(),
		PostId:    uuid.NewUUID(),
	}
}
func NewCommentFilter(t *testing.T) entities.CommentFilter {
	t.Helper()
	return entities.CommentFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewCommentCreate(t *testing.T) entities.CommentCreate {
	t.Helper()
	return entities.CommentCreate{
		Text:     faker.New().Lorem().Sentence(15),
		AuthorId: uuid.NewUUID(),
		PostId:   uuid.NewUUID(),
	}
}
func NewCommentUpdate(t *testing.T) entities.CommentUpdate {
	t.Helper()
	return entities.CommentUpdate{
		ID:       uuid.NewUUID(),
		Text:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		AuthorId: pointer.Pointer(uuid.NewUUID()),
		PostId:   pointer.Pointer(uuid.NewUUID()),
	}
}
