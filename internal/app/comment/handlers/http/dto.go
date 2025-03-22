package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type CommentDTO struct {
	ID        uuid.UUID `json:"id"`
	UpdatedAt time.Time `json:"updated_at"`
	CreatedAt time.Time `json:"created_at"`
	Text      string    `json:"text"`
	AuthorId  uuid.UUID `json:"author_id"`
	PostId    uuid.UUID `json:"post_id"`
}

func NewCommentDTO(entity entities.Comment) (CommentDTO, error) {
	dto := CommentDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		Text:      entity.Text,
		AuthorId:  entity.AuthorId,
		PostId:    entity.PostId,
	}
	return dto, nil
}

type CommentListDTO struct {
	Items []CommentDTO `json:"items"`
	Count uint64       `json:"count"`
}

func NewCommentListDTO(comments []entities.Comment, count uint64) (CommentListDTO, error) {
	response := CommentListDTO{Items: make([]CommentDTO, len(comments)), Count: count}
	for i, comment := range comments {
		dto, err := NewCommentDTO(comment)
		if err != nil {
			return CommentListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type CommentFilterDTO struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func NewCommentFilterDTO(r *http.Request) (CommentFilterDTO, error) {
	filter := CommentFilterDTO{IDs: nil, PageSize: nil, PageNumber: nil, OrderBy: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return CommentFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Pointer(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return CommentFilterDTO{}, errs.NewInvalidFormError().
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
func (dto CommentFilterDTO) toEntity() (entities.CommentFilter, error) {
	filter := entities.CommentFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		OrderBy:    dto.OrderBy,
		IDs:        dto.IDs,
	}
	return filter, nil
}

type CommentUpdateDTO struct {
	ID       uuid.UUID  `json:"id"`
	Text     *string    `json:"text"`
	AuthorId *uuid.UUID `json:"author_id"`
	PostId   *uuid.UUID `json:"post_id"`
}

func NewCommentUpdateDTO(r *http.Request) (CommentUpdateDTO, error) {
	update := CommentUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return CommentUpdateDTO{}, err
	}
	update.ID = uuid.UUID(chi.URLParam(r, "id"))
	return update, nil
}
func (dto CommentUpdateDTO) toEntity() (entities.CommentUpdate, error) {
	update := entities.CommentUpdate{
		ID:       dto.ID,
		Text:     dto.Text,
		AuthorId: dto.AuthorId,
		PostId:   dto.PostId,
	}
	return update, nil
}

type CommentCreateDTO struct {
	Text     string    `json:"text"`
	AuthorId uuid.UUID `json:"author_id"`
	PostId   uuid.UUID `json:"post_id"`
}

func NewCommentCreateDTO(r *http.Request) (CommentCreateDTO, error) {
	create := CommentCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return CommentCreateDTO{}, err
	}
	return create, nil
}
func (dto CommentCreateDTO) toEntity() (entities.CommentCreate, error) {
	create := entities.CommentCreate{Text: dto.Text, AuthorId: dto.AuthorId, PostId: dto.PostId}
	return create, nil
}
