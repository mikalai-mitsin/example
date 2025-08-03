package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Tag struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostId    uuid.UUID `json:"post_id"`
	Value     string    `json:"value"`
}

func (m *Tag) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.Value, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *TagFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(
			&m.PageSize,
			validation.Each(
				validation.In(
					"tags.id ASC",
					"tags.id DESC",
					"tags.created_at ASC",
					"tags.created_at DESC",
					"tags.updated_at ASC",
					"tags.updated_at DESC",
					"tags.post_id ASC",
					"tags.post_id DESC",
					"tags.value ASC",
					"tags.value DESC",
				),
			),
		),
		validation.Field(
			&m.PageNumber,
			validation.Each(
				validation.In(
					"tags.id ASC",
					"tags.id DESC",
					"tags.created_at ASC",
					"tags.created_at DESC",
					"tags.updated_at ASC",
					"tags.updated_at DESC",
					"tags.post_id ASC",
					"tags.post_id DESC",
					"tags.value ASC",
					"tags.value DESC",
				),
			),
		),
		validation.Field(
			&m.Search,
			validation.Each(
				validation.In(
					"tags.id ASC",
					"tags.id DESC",
					"tags.created_at ASC",
					"tags.created_at DESC",
					"tags.updated_at ASC",
					"tags.updated_at DESC",
					"tags.post_id ASC",
					"tags.post_id DESC",
					"tags.value ASC",
					"tags.value DESC",
				),
			),
		),
		validation.Field(
			&m.OrderBy,
			validation.Each(
				validation.In(
					"tags.id ASC",
					"tags.id DESC",
					"tags.created_at ASC",
					"tags.created_at DESC",
					"tags.updated_at ASC",
					"tags.updated_at DESC",
					"tags.post_id ASC",
					"tags.post_id DESC",
					"tags.value ASC",
					"tags.value DESC",
				),
			),
		),
		validation.Field(
			&m.IDs,
			validation.Each(
				validation.In(
					"tags.id ASC",
					"tags.id DESC",
					"tags.created_at ASC",
					"tags.created_at DESC",
					"tags.updated_at ASC",
					"tags.updated_at DESC",
					"tags.post_id ASC",
					"tags.post_id DESC",
					"tags.value ASC",
					"tags.value DESC",
				),
			),
		),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagCreate struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
}

func (m *TagCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PostId, validation.Required, uuid.Required),
		validation.Field(&m.Value, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type TagUpdate struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
}

func (m *TagUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.PostId),
		validation.Field(&m.Value),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
