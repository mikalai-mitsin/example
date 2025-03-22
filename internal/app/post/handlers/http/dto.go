package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type PostDTO struct {
	ID          uuid.UUID `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	AuthorId    uuid.UUID `json:"author_id"`
}

func NewPostDTO(entity entities.Post) (PostDTO, error) {
	dto := PostDTO{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Title:       entity.Title,
		Body:        entity.Body,
		IsPrivate:   entity.IsPrivate,
		Tags:        []string{},
		PublishedAt: entity.PublishedAt,
		AuthorId:    entity.AuthorId,
	}
	for _, param := range entity.Tags {
		dto.Tags = append(dto.Tags, param)
	}
	return dto, nil
}

type PostListDTO struct {
	Items []PostDTO `json:"items"`
	Count uint64    `json:"count"`
}

func NewPostListDTO(posts []entities.Post, count uint64) (PostListDTO, error) {
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
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func NewPostFilterDTO(r *http.Request) (PostFilterDTO, error) {
	filter := PostFilterDTO{IDs: nil, PageSize: nil, PageNumber: nil, OrderBy: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return PostFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Pointer(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return PostFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Pointer(uint64(pageNumber))
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	if r.URL.Query().Has("ids") {
		ids := strings.Split(r.URL.Query().Get("ids"), ",")
		filter.IDs = make([]uuid.UUID, len(ids), len(ids))
		for i, id := range ids {
			filter.IDs[i] = uuid.UUID(id)
		}
	}
	return filter, nil
}
func (dto PostFilterDTO) toEntity() (entities.PostFilter, error) {
	filter := entities.PostFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		OrderBy:    dto.OrderBy,
		IDs:        dto.IDs,
	}
	return filter, nil
}

type PostUpdateDTO struct {
	ID          uuid.UUID  `json:"id"`
	Title       *string    `json:"title"`
	Body        *string    `json:"body"`
	IsPrivate   *bool      `json:"is_private"`
	Tags        *[]string  `json:"tags"`
	PublishedAt *time.Time `json:"published_at"`
	AuthorId    *uuid.UUID `json:"author_id"`
}

func NewPostUpdateDTO(r *http.Request) (PostUpdateDTO, error) {
	update := PostUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return PostUpdateDTO{}, err
	}
	update.ID = uuid.UUID(chi.URLParam(r, "id"))
	return update, nil
}
func (dto PostUpdateDTO) toEntity() (entities.PostUpdate, error) {
	update := entities.PostUpdate{
		ID:          dto.ID,
		Title:       dto.Title,
		Body:        dto.Body,
		IsPrivate:   dto.IsPrivate,
		Tags:        dto.Tags,
		PublishedAt: dto.PublishedAt,
		AuthorId:    dto.AuthorId,
	}
	return update, nil
}

type PostCreateDTO struct {
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	IsPrivate   bool      `json:"is_private"`
	Tags        []string  `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
	AuthorId    uuid.UUID `json:"author_id"`
}

func NewPostCreateDTO(r *http.Request) (PostCreateDTO, error) {
	create := PostCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return PostCreateDTO{}, err
	}
	return create, nil
}
func (dto PostCreateDTO) toEntity() (entities.PostCreate, error) {
	create := entities.PostCreate{
		Title:       dto.Title,
		Body:        dto.Body,
		IsPrivate:   dto.IsPrivate,
		Tags:        dto.Tags,
		PublishedAt: dto.PublishedAt,
		AuthorId:    dto.AuthorId,
	}
	return create, nil
}
