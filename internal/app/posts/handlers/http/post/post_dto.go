package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostDTO struct {
	ID        uuid.UUID  `json:"id"`
	UpdatedAt time.Time  `json:"updated_at"`
	CreatedAt time.Time  `json:"created_at"`
	DeletedAt *time.Time `json:"deleted_at"`
	Body      string     `json:"body"`
}

func NewPostDTO(entity entities.Post) (PostDTO, error) {
	dto := PostDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		Body:      entity.Body,
	}
	return dto, nil
}

type PostListDTO struct {
	Items []PostDTO `json:"items"`
	Count uint64    `json:"count"`
}

func NewPostListDto(posts []entities.Post, count uint64) (PostListDTO, error) {
	response := PostListDTO{Items: make([]PostDTO, len(posts)), Count: count}
	for i, post := range posts {
		dto, err := NewPostDTO(post)
		if err != nil {
			return PostListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type PostFilterDTO struct {
	PageSize   *uint64  `json:"page_size"`
	PageNumber *uint64  `json:"page_number"`
	OrderBy    []string `json:"order_by"`
	IsDeleted  *bool    `json:"is_deleted"`
}

func NewPostFilterDTO(r *http.Request) (PostFilterDTO, error) {
	filter := PostFilterDTO{PageSize: nil, PageNumber: nil, OrderBy: nil, IsDeleted: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return PostFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Of(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return PostFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Of(uint64(pageNumber))
	}
	if r.URL.Query().Has("is_deleted") {
		isDeleted, err := strconv.ParseBool(r.URL.Query().Get("is_deleted"))
		if err != nil {
			return PostFilterDTO{}, errs.NewInvalidFormError().
				WithParam("is_deleted", "Invalid page_number.").
				WithCause(err)
		}
		filter.IsDeleted = pointer.Of(isDeleted)
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	return filter, nil
}
func (dto PostFilterDTO) toEntity() (entities.PostFilter, error) {
	filter := entities.PostFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		IsDeleted:  dto.IsDeleted,
		OrderBy:    []entities.PostOrdering{},
	}
	for _, orderBy := range dto.OrderBy {
		filter.OrderBy = append(filter.OrderBy, entities.PostOrdering(orderBy))
	}
	return filter, nil
}

type PostUpdateDTO struct {
	ID   uuid.UUID `json:"id"`
	Body *string   `json:"body"`
}

func NewPostUpdateDTO(r *http.Request) (PostUpdateDTO, error) {
	update := PostUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return PostUpdateDTO{}, err
	}
	update.ID = uuid.MustParse(chi.URLParam(r, "id"))
	return update, nil
}
func (dto PostUpdateDTO) toEntity() (entities.PostUpdate, error) {
	update := entities.PostUpdate{ID: dto.ID, Body: dto.Body}
	return update, nil
}

type PostCreateDTO struct {
	Body string `json:"body"`
}

func NewPostCreateDTO(r *http.Request) (PostCreateDTO, error) {
	create := PostCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return PostCreateDTO{}, err
	}
	return create, nil
}
func (dto PostCreateDTO) toEntity() (entities.PostCreate, error) {
	create := entities.PostCreate{Body: dto.Body}
	return create, nil
}
