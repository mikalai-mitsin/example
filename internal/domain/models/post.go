package models

import (
	"time"

	"github.com/018bf/example/internal/domain/errs"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Post struct {
	ID        string    `json:"id" db:"id,omitempty"`
	Body      string    `json:"body" db:"body"`
	Title     string    `json:"title" db:"title"`
	UserId    string    `json:"user_id" db:"user_id"`
	Weight    int       `json:"weight" db:"weight"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty"`
	CreatedAt time.Time `json:"created_at" db:"created_at,omitempty"`
}

func (c *Post) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, is.UUID),
		validation.Field(&c.Body),
		validation.Field(&c.Title),
		validation.Field(&c.UserId),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PostFilter struct {
	IDs        []string `json:"ids"`
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
}

func (c *PostFilter) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.IDs),
		validation.Field(&c.PageSize),
		validation.Field(&c.PageNumber),
		validation.Field(&c.OrderBy),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PostCreate struct {
	Body   string `json:"body"`
	Title  string `json:"title"`
	UserId string `json:"user_id"`
	Weight int    `json:"weight"`
}

func (c *PostCreate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.Body),
		validation.Field(&c.Title),
		validation.Field(&c.UserId),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

type PostUpdate struct {
	ID     string  `json:"id"`
	Body   *string `json:"body"`
	Title  *string `json:"title"`
	UserId *string `json:"user_id"`
	Weight *int    `json:"weight"`
}

func (c *PostUpdate) Validate() error {
	err := validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required, is.UUID),
		validation.Field(&c.Body),
		validation.Field(&c.Title),
		validation.Field(&c.UserId),
		validation.Field(&c.Weight),
	)
	if err != nil {
		return errs.FromValidationError(err)
	}
	return nil
}

const (
	PermissionIDPostList   PermissionID = "post_list"
	PermissionIDPostDetail PermissionID = "post_detail"
	PermissionIDPostCreate PermissionID = "post_create"
	PermissionIDPostUpdate PermissionID = "post_update"
	PermissionIDPostDelete PermissionID = "post_delete"
)
