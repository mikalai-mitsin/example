package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
}

func (m *Post) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Body, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostFilter struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	Search     *string  `json:"search"`
	OrderBy    []string `json:"order_by"`
}

func (m *PostFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(
			&m.PageSize,
			validation.Each(
				validation.In(
					"posts.id ASC",
					"posts.id DESC",
					"posts.created_at ASC",
					"posts.created_at DESC",
					"posts.updated_at ASC",
					"posts.updated_at DESC",
					"posts.body ASC",
					"posts.body DESC",
				),
			),
		),
		validation.Field(
			&m.PageNumber,
			validation.Each(
				validation.In(
					"posts.id ASC",
					"posts.id DESC",
					"posts.created_at ASC",
					"posts.created_at DESC",
					"posts.updated_at ASC",
					"posts.updated_at DESC",
					"posts.body ASC",
					"posts.body DESC",
				),
			),
		),
		validation.Field(
			&m.Search,
			validation.Each(
				validation.In(
					"posts.id ASC",
					"posts.id DESC",
					"posts.created_at ASC",
					"posts.created_at DESC",
					"posts.updated_at ASC",
					"posts.updated_at DESC",
					"posts.body ASC",
					"posts.body DESC",
				),
			),
		),
		validation.Field(
			&m.OrderBy,
			validation.Each(
				validation.In(
					"posts.id ASC",
					"posts.id DESC",
					"posts.created_at ASC",
					"posts.created_at DESC",
					"posts.updated_at ASC",
					"posts.updated_at DESC",
					"posts.body ASC",
					"posts.body DESC",
				),
			),
		),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostCreate struct {
	Body string `json:"body"`
}

func (m *PostCreate) Validate() error {
	err := validation.ValidateStruct(m, validation.Field(&m.Body, validation.Required))
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostUpdate struct {
	ID   uuid.UUID `json:"id"`
	Body *string   `json:"body"`
}

func (m *PostUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.Body),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
