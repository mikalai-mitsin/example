package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Post struct {
	ID        uuid.UUID  `json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Body      string     `json:"body"`
}

func (m *Post) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.DeletedAt),
		validation.Field(&m.Body, validation.Required),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type PostOrdering string

func (o PostOrdering) Validate() error {
	if err := validation.Validate(o.String(), validation.In(PostOrderingIdASC.String(), PostOrderingCreatedAtASC.String(), PostOrderingCreatedAtDESC.String(), PostOrderingUpdatedAtDESC.String(), PostOrderingBodyASC.String(), PostOrderingBodyDESC.String(), PostOrderingIdDESC.String(), PostOrderingUpdatedAtASC.String(), PostOrderingDeletedAtASC.String(), PostOrderingDeletedAtDESC.String())); err != nil {
		return err
	}
	return nil
}
func (o PostOrdering) String() string {
	return string(o)
}

const PostOrderingDeletedAtDESC PostOrdering = "-deleted_at"
const PostOrderingBodyASC PostOrdering = "body"
const PostOrderingIdDESC PostOrdering = "-id"
const PostOrderingBodyDESC PostOrdering = "-body"
const PostOrderingIdASC PostOrdering = "id"
const PostOrderingCreatedAtASC PostOrdering = "created_at"
const PostOrderingCreatedAtDESC PostOrdering = "-created_at"
const PostOrderingUpdatedAtASC PostOrdering = "updated_at"
const PostOrderingUpdatedAtDESC PostOrdering = "-updated_at"
const PostOrderingDeletedAtASC PostOrdering = "deleted_at"

type PostFilter struct {
	PageSize   *uint64        `json:"page_size"`
	PageNumber *uint64        `json:"page_number"`
	Search     *string        `json:"search"`
	OrderBy    []PostOrdering `json:"order_by"`
	IsDeleted  *bool          `json:"is_deleted"`
}

func (m *PostFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy),
		validation.Field(&m.IsDeleted),
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
