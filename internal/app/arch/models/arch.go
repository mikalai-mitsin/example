package models

import (
	"time"

	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/uuid"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Arch struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Tags        []string  `json:"tags"`
	Versions    []uint    `json:"versions"`
	OldVersions []uint64  `json:"old_versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
	Mark        string    `json:"mark"`
	Submarine   string    `json:"submarine"`
	Numb        uint64    `json:"numb"`
}

func (m *Arch) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Subtitle, validation.Required),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.OldVersions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
		validation.Field(&m.Mark, validation.Required),
		validation.Field(&m.Submarine, validation.Required),
		validation.Field(&m.Numb, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *ArchFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy, validation.Required),
		validation.Field(&m.IDs, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchCreate struct {
	Name        string    `json:"name"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Tags        []string  `json:"tags"`
	Versions    []uint    `json:"versions"`
	OldVersions []uint64  `json:"old_versions"`
	Release     time.Time `json:"release"`
	Tested      time.Time `json:"tested"`
	Mark        string    `json:"mark"`
	Submarine   string    `json:"submarine"`
	Numb        uint64    `json:"numb"`
}

func (m *ArchCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Name, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Subtitle, validation.Required),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.Versions, validation.Required),
		validation.Field(&m.OldVersions, validation.Required),
		validation.Field(&m.Release, validation.Required),
		validation.Field(&m.Tested, validation.Required),
		validation.Field(&m.Mark, validation.Required),
		validation.Field(&m.Submarine, validation.Required),
		validation.Field(&m.Numb, validation.Required),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type ArchUpdate struct {
	ID          uuid.UUID  `json:"id"`
	Name        *string    `json:"name"`
	Title       *string    `json:"title"`
	Subtitle    *string    `json:"subtitle"`
	Tags        *[]string  `json:"tags"`
	Versions    *[]uint    `json:"versions"`
	OldVersions *[]uint64  `json:"old_versions"`
	Release     *time.Time `json:"release"`
	Tested      *time.Time `json:"tested"`
	Mark        *string    `json:"mark"`
	Submarine   *string    `json:"submarine"`
	Numb        *uint64    `json:"numb"`
}

func (m *ArchUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Name),
		validation.Field(&m.Title),
		validation.Field(&m.Subtitle),
		validation.Field(&m.Tags),
		validation.Field(&m.Versions),
		validation.Field(&m.OldVersions),
		validation.Field(&m.Release),
		validation.Field(&m.Tested),
		validation.Field(&m.Mark),
		validation.Field(&m.Submarine),
		validation.Field(&m.Numb),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}
