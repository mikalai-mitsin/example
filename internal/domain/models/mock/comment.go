package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewComment(t *testing.T) *models.Comment {
	t.Helper()
	return &models.Comment{
		ID:        uuid.New().String(),
		Body:      faker.Lorem().String(),
		PostId:    faker.Lorem().String(),
		UserId:    faker.Lorem().String(),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewCommentCreate(t *testing.T) *models.CommentCreate {
	t.Helper()
	return &models.CommentCreate{
		Body:   faker.Lorem().String(),
		PostId: faker.Lorem().String(),
		UserId: faker.Lorem().String(),
	}
}

func NewCommentUpdate(t *testing.T) *models.CommentUpdate {
	t.Helper()
	return &models.CommentUpdate{
		ID:     uuid.New().String(),
		Body:   utils.Pointer(faker.Lorem().String()),
		PostId: utils.Pointer(faker.Lorem().String()),
		UserId: utils.Pointer(faker.Lorem().String()),
	}
}

func NewCommentFilter(t *testing.T) *models.CommentFilter {
	t.Helper()
	return &models.CommentFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
