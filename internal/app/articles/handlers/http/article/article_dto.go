package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArticleDTO struct {
	ID          uuid.UUID `json:"id"`
	UpdatedAt   time.Time `json:"updated_at"`
	CreatedAt   time.Time `json:"created_at"`
	Title       string    `json:"title"`
	Subtitle    string    `json:"subtitle"`
	Body        string    `json:"body"`
	IsPublished string    `json:"is_published"`
}

func NewArticleDTO(entity entities.Article) (ArticleDTO, error) {
	dto := ArticleDTO{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Title:       entity.Title,
		Subtitle:    entity.Subtitle,
		Body:        entity.Body,
		IsPublished: entity.IsPublished,
	}
	return dto, nil
}

type ArticleListDTO struct {
	Items []ArticleDTO `json:"items"`
	Count uint64       `json:"count"`
}

func NewArticleListDto(articles []entities.Article, count uint64) (ArticleListDTO, error) {
	response := ArticleListDTO{Items: make([]ArticleDTO, len(articles)), Count: count}
	for i, article := range articles {
		dto, err := NewArticleDTO(article)
		if err != nil {
			return ArticleListDTO{}, err
		}
		response.Items[i] = dto
	}
	return response, nil
}

type ArticleFilterDTO struct {
	PageSize   *uint64     `json:"page_size"`
	PageNumber *uint64     `json:"page_number"`
	OrderBy    []string    `json:"order_by"`
	IDs        []uuid.UUID `json:"ids"`
}

func NewArticleFilterDTO(r *http.Request) (ArticleFilterDTO, error) {
	filter := ArticleFilterDTO{IDs: nil, PageSize: nil, PageNumber: nil, OrderBy: nil}
	if r.URL.Query().Has("page_size") {
		pageSize, err := strconv.Atoi(r.URL.Query().Get("page_size"))
		if err != nil {
			return ArticleFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_size", "Invalid page_size.").
				WithCause(err)
		}
		filter.PageSize = pointer.Of(uint64(pageSize))
	}
	if r.URL.Query().Has("page_number") {
		pageNumber, err := strconv.Atoi(r.URL.Query().Get("page_number"))
		if err != nil {
			return ArticleFilterDTO{}, errs.NewInvalidFormError().
				WithParam("page_number", "Invalid page_number.").
				WithCause(err)
		}
		filter.PageNumber = pointer.Of(uint64(pageNumber))
	}
	if r.URL.Query().Has("order_by") {
		filter.OrderBy = strings.Split(r.URL.Query().Get("order_by"), ",")
	}
	if r.URL.Query().Has("ids") {
		ids := strings.Split(r.URL.Query().Get("ids"), ",")
		filter.IDs = make([]uuid.UUID, len(ids))
		for i, id := range ids {
			filter.IDs[i] = uuid.MustParse(id)
		}
	}
	return filter, nil
}
func (dto ArticleFilterDTO) toEntity() (entities.ArticleFilter, error) {
	filter := entities.ArticleFilter{
		PageSize:   dto.PageSize,
		PageNumber: dto.PageNumber,
		OrderBy:    dto.OrderBy,
		IDs:        dto.IDs,
	}
	return filter, nil
}

type ArticleUpdateDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       *string   `json:"title"`
	Subtitle    *string   `json:"subtitle"`
	Body        *string   `json:"body"`
	IsPublished *string   `json:"is_published"`
}

func NewArticleUpdateDTO(r *http.Request) (ArticleUpdateDTO, error) {
	update := ArticleUpdateDTO{}
	if err := render.DecodeJSON(r.Body, &update); err != nil {
		return ArticleUpdateDTO{}, err
	}
	update.ID = uuid.MustParse(chi.URLParam(r, "id"))
	return update, nil
}
func (dto ArticleUpdateDTO) toEntity() (entities.ArticleUpdate, error) {
	update := entities.ArticleUpdate{
		ID:          dto.ID,
		Title:       dto.Title,
		Subtitle:    dto.Subtitle,
		Body:        dto.Body,
		IsPublished: dto.IsPublished,
	}
	return update, nil
}

type ArticleCreateDTO struct {
	Title       string `json:"title"`
	Subtitle    string `json:"subtitle"`
	Body        string `json:"body"`
	IsPublished string `json:"is_published"`
}

func NewArticleCreateDTO(r *http.Request) (ArticleCreateDTO, error) {
	create := ArticleCreateDTO{}
	if err := render.DecodeJSON(r.Body, &create); err != nil {
		return ArticleCreateDTO{}, err
	}
	return create, nil
}
func (dto ArticleCreateDTO) toEntity() (entities.ArticleCreate, error) {
	create := entities.ArticleCreate{
		Title:       dto.Title,
		Subtitle:    dto.Subtitle,
		Body:        dto.Body,
		IsPublished: dto.IsPublished,
	}
	return create, nil
}
