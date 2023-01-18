package models

import (
    "time"

    "github.com/018bf/example/internal/domain/errs"
    validation "github.com/go-ozzo/ozzo-validation/v4"
    "github.com/go-ozzo/ozzo-validation/v4/is"
)

type Mark struct {
    ID        string `json:"id" db:"id,omitempty" form:"id"`
    Name string `json:"name" db:"name" form:"name"`
    Title string `json:"title" db:"title" form:"title"`
    Weight int `json:"weight" db:"weight" form:"weight"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at,omitempty" form:"updated_at"`
    CreatedAt time.Time `json:"created_at" db:"created_at,omitempty" form:"created_at,omitempty"`
}

func (c *Mark) Validate() error {
    err := validation.ValidateStruct(
        c,
        validation.Field(&c.ID, is.UUID),
        validation.Field(&c.Name),
        validation.Field(&c.Title),
        validation.Field(&c.Weight),
    )
    if err != nil {
        return errs.FromValidationError(err)
    }
    return nil
}

type MarkFilter struct {
    IDs        []string `json:"ids" form:"ids"`
    PageSize   *uint64  `json:"page_size" form:"page_size"`
    PageNumber *uint64  `json:"page_number" form:"page_number"`
    OrderBy    []string `json:"order_by" form:"order_by"`
}

func (c *MarkFilter) Validate() error {
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

type MarkCreate struct {
    Name string `json:"name" form:"name"`
    Title string `json:"title" form:"title"`
    Weight int `json:"weight" form:"weight"`
}

func (c *MarkCreate) Validate() error {
    err := validation.ValidateStruct(
        c,
        validation.Field(&c.Name),
        validation.Field(&c.Title),
        validation.Field(&c.Weight),
    )
    if err != nil {
        return errs.FromValidationError(err)
    }
    return nil
}

type MarkUpdate struct {
    ID string `json:"id"`
    Name *string `json:"name" form:"name"`
    Title *string `json:"title" form:"title"`
    Weight *int `json:"weight" form:"weight"`
}

func (c *MarkUpdate) Validate() error {
    err := validation.ValidateStruct(
        c,
        validation.Field(&c.ID, validation.Required, is.UUID),
        validation.Field(&c.Name),
        validation.Field(&c.Title),
        validation.Field(&c.Weight),
    )
    if err != nil {
        return errs.FromValidationError(err)
    }
    return nil
}

const (
    PermissionIDMarkList PermissionID = "mark_list"
    PermissionIDMarkDetail PermissionID = "mark_detail"
    PermissionIDMarkCreate PermissionID = "mark_create"
    PermissionIDMarkUpdate PermissionID = "mark_update"
    PermissionIDMarkDelete PermissionID = "mark_delete"
)
