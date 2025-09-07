package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArticleRepository struct {
	readDB  database
	writeDB database
	logger  logger
}

func NewArticleRepository(readDB database, writeDB database, logger logger) *ArticleRepository {
	return &ArticleRepository{readDB: readDB, writeDB: writeDB, logger: logger}
}

var orderByMap = map[entities.ArticleOrdering]string{
	entities.ArticleOrderingBodyDESC:        "articles.body DESC",
	entities.ArticleOrderingIsPublishedDESC: "articles.is_published DESC",
	entities.ArticleOrderingIdASC:           "articles.id ASC",
	entities.ArticleOrderingIdDESC:          "articles.id DESC",
	entities.ArticleOrderingCreatedAtASC:    "articles.created_at ASC",
	entities.ArticleOrderingUpdatedAtASC:    "articles.updated_at ASC",
	entities.ArticleOrderingUpdatedAtDESC:   "articles.updated_at DESC",
	entities.ArticleOrderingTitleASC:        "articles.title ASC",
	entities.ArticleOrderingSubtitleDESC:    "articles.subtitle DESC",
	entities.ArticleOrderingIsPublishedASC:  "articles.is_published ASC",
	entities.ArticleOrderingCreatedAtDESC:   "articles.created_at DESC",
	entities.ArticleOrderingTitleDESC:       "articles.title DESC",
	entities.ArticleOrderingSubtitleASC:     "articles.subtitle ASC",
	entities.ArticleOrderingBodyASC:         "articles.body ASC",
}

func encodeOrderBy(orderBy []entities.ArticleOrdering) []string {
	columns := make([]string, len(orderBy))
	for i, item := range orderBy {
		column, exists := orderByMap[item]
		if !exists {
			continue
		}
		columns[i] = column
	}
	return columns
}

type ArticleDTO struct {
	ID          uuid.UUID `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Title       string    `db:"title"`
	Subtitle    string    `db:"subtitle"`
	Body        string    `db:"body"`
	IsPublished bool      `db:"is_published"`
}
type ArticleListDTO []ArticleDTO

func (list ArticleListDTO) toEntities() []entities.Article {
	items := make([]entities.Article, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewArticleDTOFromEntity(entity entities.Article) ArticleDTO {
	dto := ArticleDTO{
		ID:          entity.ID,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		Title:       entity.Title,
		Subtitle:    entity.Subtitle,
		Body:        entity.Body,
		IsPublished: entity.IsPublished,
	}
	return dto
}
func (dto ArticleDTO) toEntity() entities.Article {
	entity := entities.Article{
		ID:          dto.ID,
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		Title:       dto.Title,
		Subtitle:    dto.Subtitle,
		Body:        dto.Body,
		IsPublished: dto.IsPublished,
	}
	return entity
}
func (r *ArticleRepository) Create(ctx context.Context, tx dtx.TX, entity entities.Article) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewArticleDTOFromEntity(entity)
	q := sq.Insert("public.articles").
		Columns("id", "created_at", "updated_at", "title", "subtitle", "body", "is_published").
		Values(dto.ID, dto.CreatedAt, dto.UpdatedAt, dto.Title, dto.Subtitle, dto.Body, dto.IsPublished)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := tx.GetSQLTx().ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *ArticleRepository) Get(ctx context.Context, id uuid.UUID) (entities.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &ArticleDTO{}
	q := sq.Select("articles.id", "articles.created_at", "articles.updated_at", "articles.title", "articles.subtitle", "articles.body", "articles.is_published").
		From("public.articles").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("article_id", id.String())
		return entities.Article{}, e
	}
	return dto.toEntity(), nil
}

func (r *ArticleRepository) List(
	ctx context.Context,
	filter entities.ArticleFilter,
) ([]entities.Article, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto ArticleListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Of(pageSize)
	}
	q := sq.Select("articles.id", "articles.created_at", "articles.updated_at", "articles.title", "articles.subtitle", "articles.body", "articles.is_published").
		From("public.articles").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"title", "subtitle", "body"},
			},
		)
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(encodeOrderBy(filter.OrderBy)...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.toEntities(), nil
}

func (r *ArticleRepository) Count(
	ctx context.Context,
	filter entities.ArticleFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.articles")
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{
				Lang:   "english",
				Query:  *filter.Search,
				Fields: []string{"title", "subtitle", "body"},
			},
		)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	var count uint64
	if err := r.readDB.GetContext(ctx, &count, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *ArticleRepository) Update(ctx context.Context, tx dtx.TX, entity entities.Article) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewArticleDTOFromEntity(entity)
	q := sq.Update("public.articles").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("title", dto.Title)
		q = q.Set("subtitle", dto.Subtitle)
		q = q.Set("body", dto.Body)
		q = q.Set("is_published", dto.IsPublished)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("article_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("article_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("article_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *ArticleRepository) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.articles").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("article_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("article_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("article_id", fmt.Sprint(id))
		return e
	}
	return nil
}
