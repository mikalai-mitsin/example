package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type TagRepository struct {
	readDB  database
	writeDB database
	logger  logger
}

func NewTagRepository(readDB database, writeDB database, logger logger) *TagRepository {
	return &TagRepository{readDB: readDB, writeDB: writeDB, logger: logger}
}

var orderByMap = map[entities.TagOrdering]string{
	entities.TagOrderingCreatedAtDESC: "tags.created_at DESC",
	entities.TagOrderingUpdatedAtASC:  "tags.updated_at ASC",
	entities.TagOrderingPostIdASC:     "tags.post_id ASC",
	entities.TagOrderingIdASC:         "tags.id ASC",
	entities.TagOrderingIdDESC:        "tags.id DESC",
	entities.TagOrderingUpdatedAtDESC: "tags.updated_at DESC",
	entities.TagOrderingPostIdDESC:    "tags.post_id DESC",
	entities.TagOrderingValueASC:      "tags.value ASC",
	entities.TagOrderingValueDESC:     "tags.value DESC",
	entities.TagOrderingCreatedAtASC:  "tags.created_at ASC",
}

func encodeOrderBy(orderBy []entities.TagOrdering) []string {
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

type TagDTO struct {
	ID        uuid.UUID `db:"id,omitempty"`
	UpdatedAt time.Time `db:"updated_at,omitempty"`
	CreatedAt time.Time `db:"created_at,omitempty"`
	PostId    uuid.UUID `db:"post_id"`
	Value     string    `db:"value"`
}
type TagListDTO []TagDTO

func (list TagListDTO) toEntities() []entities.Tag {
	items := make([]entities.Tag, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewTagDTOFromEntity(entity entities.Tag) TagDTO {
	dto := TagDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		PostId:    entity.PostId,
		Value:     entity.Value,
	}
	return dto
}
func (dto TagDTO) toEntity() entities.Tag {
	entity := entities.Tag{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		PostId:    dto.PostId,
		Value:     dto.Value,
	}
	return entity
}
func (r *TagRepository) Create(ctx context.Context, tx dtx.TX, entity entities.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewTagDTOFromEntity(entity)
	q := sq.Insert("public.tags").
		Columns("id", "created_at", "updated_at", "post_id", "value").
		Values(dto.ID, dto.CreatedAt, dto.UpdatedAt, dto.PostId, dto.Value)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := tx.GetSQLTx().ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *TagRepository) Get(ctx context.Context, id uuid.UUID) (entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &TagDTO{}
	q := sq.Select("tags.id", "tags.created_at", "tags.updated_at", "tags.post_id", "tags.value").
		From("public.tags").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", id.String())
		return entities.Tag{}, e
	}
	return dto.toEntity(), nil
}

func (r *TagRepository) List(
	ctx context.Context,
	filter entities.TagFilter,
) ([]entities.Tag, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto TagListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Of(pageSize)
	}
	q := sq.Select("tags.id", "tags.created_at", "tags.updated_at", "tags.post_id", "tags.value").
		From("public.tags").
		Limit(pageSize)
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
func (r *TagRepository) Count(ctx context.Context, filter entities.TagFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.tags")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	var count uint64
	if err := r.readDB.GetContext(ctx, &count, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *TagRepository) Update(ctx context.Context, tx dtx.TX, entity entities.Tag) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewTagDTOFromEntity(entity)
	q := sq.Update("public.tags").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("post_id", dto.PostId)
		q = q.Set("value", dto.Value)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("tag_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *TagRepository) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.tags").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("tag_id", fmt.Sprint(id))
		return e
	}
	return nil
}
