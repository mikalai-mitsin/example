package entities

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type Article struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Body        string    `json:"body"`
	IsPublished bool      `json:"is_published"`
}

func (m *Article) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.CreatedAt, validation.Required),
		validation.Field(&m.UpdatedAt, validation.Required),
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Subtitle, validation.Required),
		validation.Field(&m.Body, validation.Required),
		validation.Field(&m.IsPublished),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type ArticleOrdering string

func (o ArticleOrdering) Validate() error {
	if err := validation.Validate(o.String(), validation.In(ArticleOrderingCreatedAtDESC.String(), ArticleOrderingTitleDESC.String(), ArticleOrderingSubtitleASC.String(), ArticleOrderingSubtitleDESC.String(), ArticleOrderingBodyASC.String(), ArticleOrderingBodyDESC.String(), ArticleOrderingIdDESC.String(), ArticleOrderingUpdatedAtASC.String(), ArticleOrderingUpdatedAtDESC.String(), ArticleOrderingTitleASC.String(), ArticleOrderingIsPublishedASC.String(), ArticleOrderingIsPublishedDESC.String(), ArticleOrderingIdASC.String(), ArticleOrderingCreatedAtASC.String())); err != nil {
		return err
	}
	return nil
}
func (o ArticleOrdering) String() string {
	return string(o)
}

const ArticleOrderingBodyDESC ArticleOrdering = "-body"
const ArticleOrderingIsPublishedASC ArticleOrdering = "is_published"
const ArticleOrderingCreatedAtASC ArticleOrdering = "created_at"
const ArticleOrderingCreatedAtDESC ArticleOrdering = "-created_at"
const ArticleOrderingUpdatedAtASC ArticleOrdering = "updated_at"
const ArticleOrderingTitleASC ArticleOrdering = "title"
const ArticleOrderingTitleDESC ArticleOrdering = "-title"
const ArticleOrderingSubtitleASC ArticleOrdering = "subtitle"
const ArticleOrderingIsPublishedDESC ArticleOrdering = "-is_published"
const ArticleOrderingIdASC ArticleOrdering = "id"
const ArticleOrderingIdDESC ArticleOrdering = "-id"
const ArticleOrderingUpdatedAtDESC ArticleOrdering = "-updated_at"
const ArticleOrderingSubtitleDESC ArticleOrdering = "-subtitle"
const ArticleOrderingBodyASC ArticleOrdering = "body"

type ArticleFilter struct {
	PageSize   *uint64           `json:"page_size"`
	PageNumber *uint64           `json:"page_number"`
	Search     *string           `json:"search"`
	OrderBy    []ArticleOrdering `json:"order_by"`
}

func (m *ArticleFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(&m.OrderBy),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type ArticleCreate struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Body        string `json:"body"`
	IsPublished bool   `json:"is_published"`
}

func (m *ArticleCreate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.Title, validation.Required),
		validation.Field(&m.Subtitle, validation.Required),
		validation.Field(&m.Body, validation.Required),
		validation.Field(&m.IsPublished),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}

type ArticleUpdate struct {
	ID          uuid.UUID `json:"id"`
	Title       *string   `json:"title"`
	Subtitle    *string   `json:"subtitle"`
	Body        *string   `json:"body"`
	IsPublished *bool     `json:"is_published"`
}

func (m *ArticleUpdate) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.ID, validation.Required),
		validation.Field(&m.Title),
		validation.Field(&m.Subtitle),
		validation.Field(&m.Body),
		validation.Field(&m.IsPublished),
	)
	if err != nil {
		return errs.NewFromValidationError(err)
	}
	return nil
}
