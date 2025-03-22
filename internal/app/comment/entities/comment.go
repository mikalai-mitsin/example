package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Comment struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Text      string    `json:"text"`
	AuthorId  uuid.UUID `json:"author_id"`
	PostId    uuid.UUID `json:"post_id"`
}

func (m *Comment) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Text, validation.Required),
		validation.Field(&m.AuthorId, validation.Required, is.UUID),
		validation.Field(&m.PostId, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type CommentFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *CommentFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy, validation.Required),
		validation.Field(&m.IDs, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type CommentCreate struct {
	Text     string    `json:"text"`
	AuthorId uuid.UUID `json:"author_id"`
	PostId   uuid.UUID `json:"post_id"`
}

func (m *CommentCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Text, validation.Required),
		validation.Field(&m.AuthorId, validation.Required, is.UUID),
		validation.Field(&m.PostId, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type CommentUpdate struct {
	ID       uuid.UUID  `json:"id"`
	Text     *string    `json:"text"`
	AuthorId *uuid.UUID `json:"author_id"`
	PostId   *uuid.UUID `json:"post_id"`
}

func (m *CommentUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Text),
		validation.Field(&m.AuthorId),
		validation.Field(&m.PostId),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
