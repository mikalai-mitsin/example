package mock_models // nolint:stylecheck

import (
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/pkg/utils"
	"github.com/google/uuid"
	"syreclabs.com/go/faker"
)

func NewPost(t *testing.T) *models.Post {
	t.Helper()
	return &models.Post{
		ID:        uuid.New().String(),
		Body:      faker.Lorem().String(),
		Title:     faker.Lorem().String(),
		UserId:    faker.Lorem().String(),
		Weight:    int(faker.RandomInt(2, 100)),
		UpdatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
		CreatedAt: faker.Time().Backward(40 * time.Hour).UTC(),
	}
}

func NewPostCreate(t *testing.T) *models.PostCreate {
	t.Helper()
	return &models.PostCreate{
		Body:   faker.Lorem().String(),
		Title:  faker.Lorem().String(),
		UserId: faker.Lorem().String(),
		Weight: int(faker.RandomInt(2, 100)),
	}
}

func NewPostUpdate(t *testing.T) *models.PostUpdate {
	t.Helper()
	return &models.PostUpdate{
		ID:     uuid.New().String(),
		Body:   utils.Pointer(faker.Lorem().String()),
		Title:  utils.Pointer(faker.Lorem().String()),
		UserId: utils.Pointer(faker.Lorem().String()),
		Weight: utils.Pointer(int(faker.RandomInt(2, 100))),
	}
}

func NewPostFilter(t *testing.T) *models.PostFilter {
	t.Helper()
	return &models.PostFilter{
		PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
		OrderBy:    faker.Lorem().Words(5),
		IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
	}
}
