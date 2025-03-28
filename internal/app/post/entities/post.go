package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Post struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	AuthorId    uuid.UUID `json:"author_id"`
}

func (m *Post) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Body, validation.Required),
		validation.Field(&m.IsPrivate),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.PublishedAt, validation.Required),
		validation.Field(&m.AuthorId, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostFilter struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	Search     *string     `json:"search"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func (m *PostFilter) Validate() error {
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

type PostCreate struct {
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	AuthorId    uuid.UUID `json:"author_id"`
}

func (m *PostCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Body, validation.Required),
		validation.Field(&m.IsPrivate),
		validation.Field(&m.Tags, validation.Required),
		validation.Field(&m.PublishedAt, validation.Required),
		validation.Field(&m.AuthorId, validation.Required, is.UUID),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostUpdate struct {
	ID          uuid.UUID  `json:"id"`
	Title       *string    `json:"title"`
	Body        *string    `json:"body"`
	IsPrivate   *bool      `json:"is_private"`
	Tags        *[]string  `json:"tags"`
	PublishedAt *time.Time `json:"published_at"`
	AuthorId    *uuid.UUID `json:"author_id"`
}

func (m *PostUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required, is.UUID),
		validation.Field(&m.Title),
		validation.Field(&m.Body),
		validation.Field(&m.IsPrivate),
		validation.Field(&m.Tags),
		validation.Field(&m.PublishedAt),
		validation.Field(&m.AuthorId),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
