package mock_models // nolint:stylecheck

import (
    "testing"
    "time"

    "github.com/018bf/example/internal/domain/models"
    "github.com/018bf/example/pkg/utils"
    "github.com/google/uuid"
    "syreclabs.com/go/faker"
)

func NewMark(t *testing.T) *models.Mark {
    t.Helper()
    return &models.Mark{
        ID:         uuid.New().String(),
        Name: faker.Lorem().String(),
        Title: faker.Lorem().String(),
        Weight: faker.RandomInt(2, 100),
        UpdatedAt:  faker.Time().Backward(40 * time.Hour).UTC(),
        CreatedAt:  faker.Time().Backward(40 * time.Hour).UTC(),
    }
}

func NewMarkCreate(t *testing.T) *models.MarkCreate {
    t.Helper()
    return &models.MarkCreate{
        Name: faker.Lorem().String(),
        Title: faker.Lorem().String(),
        Weight: faker.RandomInt(2, 100),
    }
}

func NewMarkUpdate(t *testing.T) *models.MarkUpdate {
    t.Helper()
    return &models.MarkUpdate{
        ID: uuid.New().String(),
        Name:  utils.Pointer(faker.Lorem().String()),
        Title:  utils.Pointer(faker.Lorem().String()),
        Weight:  utils.Pointer(faker.RandomInt(2, 100)),
    }
}

func NewMarkFilter(t *testing.T) *models.MarkFilter {
    t.Helper()
    return &models.MarkFilter{
        PageSize:   utils.Pointer(uint64(faker.RandomInt64(2, 100))),
        PageNumber: utils.Pointer(uint64(faker.RandomInt64(2, 100))),
        OrderBy:    faker.Lorem().Words(5),
        IDs:        []string{uuid.New().String(), uuid.New().String(), uuid.New().String()},
    }
}
