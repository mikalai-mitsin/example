package mock_models

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewArch(t *testing.T) *models.Arch {
	t.Helper()
	return &models.Arch{
		ID:          uuid.NewUUID(),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		Name:        faker.New().Lorem().Sentence(15),
		Title:       faker.New().Lorem().Sentence(15),
		Subtitle:    faker.New().Lorem().Sentence(15),
		Tags:        []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		Versions:    []uint{faker.New().UInt(), faker.New().UInt()},
		OldVersions: []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
		Mark:        faker.New().Lorem().Sentence(15),
		Submarine:   faker.New().Lorem().Sentence(15),
		Numb:        faker.New().UInt64(),
	}
}
func NewArchFilter(t *testing.T) *models.ArchFilter {
	t.Helper()
	return &models.ArchFilter{
		PageSize:   pointer.Pointer(faker.New().UInt64()),
		PageNumber: pointer.Pointer(faker.New().UInt64()),
		Search:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		OrderBy:    []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		IDs:        []uuid.UUID{uuid.NewUUID(), uuid.NewUUID()},
	}
}
func NewArchCreate(t *testing.T) *models.ArchCreate {
	t.Helper()
	return &models.ArchCreate{
		Name:        faker.New().Lorem().Sentence(15),
		Title:       faker.New().Lorem().Sentence(15),
		Subtitle:    faker.New().Lorem().Sentence(15),
		Tags:        []string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		Versions:    []uint{faker.New().UInt(), faker.New().UInt()},
		OldVersions: []uint64{faker.New().UInt64(), faker.New().UInt64()},
		Release:     faker.New().Time().Time(time.Now()),
		Tested:      faker.New().Time().Time(time.Now()),
		Mark:        faker.New().Lorem().Sentence(15),
		Submarine:   faker.New().Lorem().Sentence(15),
		Numb:        faker.New().UInt64(),
	}
}
func NewArchUpdate(t *testing.T) *models.ArchUpdate {
	t.Helper()
	return &models.ArchUpdate{
		ID:       uuid.NewUUID(),
		Name:     pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Title:    pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Subtitle: pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Tags: pointer.Pointer(
			[]string{faker.New().Lorem().Sentence(15), faker.New().Lorem().Sentence(15)},
		),
		Versions:    pointer.Pointer([]uint{faker.New().UInt(), faker.New().UInt()}),
		OldVersions: pointer.Pointer([]uint64{faker.New().UInt64(), faker.New().UInt64()}),
		Release:     pointer.Pointer(faker.New().Time().Time(time.Now())),
		Tested:      pointer.Pointer(faker.New().Time().Time(time.Now())),
		Mark:        pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Submarine:   pointer.Pointer(faker.New().Lorem().Sentence(15)),
		Numb:        pointer.Pointer(faker.New().UInt64()),
	}
}
