package usecases

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/arch/models"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArchUseCase struct {
	archRepository ArchRepository
	clock          Clock
	logger         Logger
	uuid           UUIDGenerator
}

func NewArchUseCase(
	archRepository ArchRepository,
	clock Clock,
	logger Logger,
	uuid UUIDGenerator,
) *ArchUseCase {
	return &ArchUseCase{archRepository: archRepository, clock: clock, logger: logger, uuid: uuid}
}
func (u *ArchUseCase) Create(ctx context.Context, create *models.ArchCreate) (*models.Arch, error) {
	if err := create.Validate(); err != nil {
		return nil, err
	}
	now := u.clock.Now().UTC()
	arch := &models.Arch{
		ID:          u.uuid.NewUUID(),
		UpdatedAt:   now,
		CreatedAt:   now,
		Name:        create.Name,
		Title:       create.Title,
		Subtitle:    create.Subtitle,
		Tags:        create.Tags,
		Versions:    create.Versions,
		OldVersions: create.OldVersions,
		Release:     create.Release,
		Tested:      create.Tested,
		Mark:        create.Mark,
		Submarine:   create.Submarine,
		Numb:        create.Numb,
	}
	if err := u.archRepository.Create(ctx, arch); err != nil {
		return nil, err
	}
	return arch, nil
}

func (u *ArchUseCase) List(
	ctx context.Context,
	filter *models.ArchFilter,
) ([]*models.Arch, uint64, error) {
	arch, err := u.archRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.archRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return arch, count, nil
}
func (u *ArchUseCase) Get(ctx context.Context, id uuid.UUID) (*models.Arch, error) {
	arch, err := u.archRepository.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return arch, nil
}
func (u *ArchUseCase) Update(ctx context.Context, update *models.ArchUpdate) (*models.Arch, error) {
	if err := update.Validate(); err != nil {
		return nil, err
	}
	arch, err := u.archRepository.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	{
		if update.Name != nil {
			arch.Name = *update.Name
		}
		if update.Title != nil {
			arch.Title = *update.Title
		}
		if update.Subtitle != nil {
			arch.Subtitle = *update.Subtitle
		}
		if update.Tags != nil {
			arch.Tags = *update.Tags
		}
		if update.Versions != nil {
			arch.Versions = *update.Versions
		}
		if update.OldVersions != nil {
			arch.OldVersions = *update.OldVersions
		}
		if update.Release != nil {
			arch.Release = *update.Release
		}
		if update.Tested != nil {
			arch.Tested = *update.Tested
		}
		if update.Mark != nil {
			arch.Mark = *update.Mark
		}
		if update.Submarine != nil {
			arch.Submarine = *update.Submarine
		}
		if update.Numb != nil {
			arch.Numb = *update.Numb
		}
	}
	arch.UpdatedAt = u.clock.Now().UTC()
	if err := u.archRepository.Update(ctx, arch); err != nil {
		return nil, err
	}
	return arch, nil
}
func (u *ArchUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.archRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
