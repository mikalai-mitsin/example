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

type ArticleFilter struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	Search     *string  `json:"search"`
	OrderBy    []string `json:"order_by"`
}

func (m *ArticleFilter) Validate() error {
	err := validation.ValidateStruct(
		m,
		validation.Field(&m.PageSize),
		validation.Field(&m.PageNumber),
		validation.Field(&m.Search),
		validation.Field(
			&m.OrderBy,
			validation.Each(
				validation.In(
					"articles.id ASC",
					"articles.id DESC",
					"articles.created_at ASC",
					"articles.created_at DESC",
					"articles.updated_at ASC",
					"articles.updated_at DESC",
					"articles.title ASC",
					"articles.title DESC",
					"articles.subtitle ASC",
					"articles.subtitle DESC",
					"articles.body ASC",
					"articles.body DESC",
					"articles.is_published ASC",
					"articles.is_published DESC",
				),
			),
		),
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
