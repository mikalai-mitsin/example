package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Like struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	PostId    uuid.UUID `json:"post_id"`
	Value     string    `json:"value"`
	UserId    uuid.UUID `json:"user_id"`
}

func (m *Like) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.PostId, validation.Required),
		validation.Field(&m.Value, validation.Required),
		validation.Field(&m.UserId, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeFilter struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	Search     *string  `json:"search"`
	OrderBy    []string `json:"order_by"`
}

func (m *LikeFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(
			&m.OrderBy,
			validation.Each(
				validation.In(
					"likes.id ASC",
					"likes.id DESC",
					"likes.created_at ASC",
					"likes.created_at DESC",
					"likes.updated_at ASC",
					"likes.updated_at DESC",
					"likes.post_id ASC",
					"likes.post_id DESC",
					"likes.value ASC",
					"likes.value DESC",
					"likes.user_id ASC",
					"likes.user_id DESC",
				),
			),
		),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeCreate struct {
	PostId uuid.UUID `json:"post_id"`
	Value  string    `json:"value"`
	UserId uuid.UUID `json:"user_id"`
}

func (m *LikeCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PostId, validation.Required, uuid.Required),
		validation.Field(&m.Value, validation.Required),
		validation.Field(&m.UserId, validation.Required, uuid.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type LikeUpdate struct {
	ID     uuid.UUID  `json:"id"`
	PostId *uuid.UUID `json:"post_id"`
	Value  *string    `json:"value"`
	UserId *uuid.UUID `json:"user_id"`
}

func (m *LikeUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.PostId),
		validation.Field(&m.Value),
		validation.Field(&m.UserId),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
