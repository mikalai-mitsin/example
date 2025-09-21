package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type LikeRepository struct {
	readDB  database
	writeDB database
	logger  logger
}

func NewLikeRepository(readDB database, writeDB database, logger logger) *LikeRepository {
	return &LikeRepository{readDB: readDB, writeDB: writeDB, logger: logger}
}

var orderByMap = map[entities.LikeOrdering]string{
	entities.LikeOrderingUpdatedAtASC:  "likes.updated_at ASC",
	entities.LikeOrderingDeletedAtASC:  "likes.deleted_at ASC",
	entities.LikeOrderingPostIdASC:     "likes.post_id ASC",
	entities.LikeOrderingValueDESC:     "likes.value DESC",
	entities.LikeOrderingUserIdASC:     "likes.user_id ASC",
	entities.LikeOrderingIdASC:         "likes.id ASC",
	entities.LikeOrderingCreatedAtDESC: "likes.created_at DESC",
	entities.LikeOrderingUpdatedAtDESC: "likes.updated_at DESC",
	entities.LikeOrderingDeletedAtDESC: "likes.deleted_at DESC",
	entities.LikeOrderingPostIdDESC:    "likes.post_id DESC",
	entities.LikeOrderingValueASC:      "likes.value ASC",
	entities.LikeOrderingUserIdDESC:    "likes.user_id DESC",
	entities.LikeOrderingIdDESC:        "likes.id DESC",
	entities.LikeOrderingCreatedAtASC:  "likes.created_at ASC",
}

func encodeOrderBy(orderBy []entities.LikeOrdering) []string {
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

type LikeDTO struct {
	ID        uuid.UUID  `db:"id,omitempty"`
	UpdatedAt time.Time  `db:"updated_at,omitempty"`
	CreatedAt time.Time  `db:"created_at,omitempty"`
	DeletedAt *time.Time `db:"deleted_at"`
	PostId    uuid.UUID  `db:"post_id"`
	Value     string     `db:"value"`
	UserId    uuid.UUID  `db:"user_id"`
}
type LikeListDTO []LikeDTO

func (list LikeListDTO) toEntities() []entities.Like {
	items := make([]entities.Like, len(list))
	for i := range list {
		items[i] = list[i].toEntity()
	}
	return items
}
func NewLikeDTOFromEntity(entity entities.Like) LikeDTO {
	dto := LikeDTO{
		ID:        entity.ID,
		CreatedAt: entity.CreatedAt,
		UpdatedAt: entity.UpdatedAt,
		DeletedAt: entity.DeletedAt,
		PostId:    entity.PostId,
		Value:     entity.Value,
		UserId:    entity.UserId,
	}
	return dto
}
func (dto LikeDTO) toEntity() entities.Like {
	entity := entities.Like{
		ID:        dto.ID,
		CreatedAt: dto.CreatedAt,
		UpdatedAt: dto.UpdatedAt,
		DeletedAt: dto.DeletedAt,
		PostId:    dto.PostId,
		Value:     dto.Value,
		UserId:    dto.UserId,
	}
	return entity
}
func (r *LikeRepository) Create(ctx context.Context, tx dtx.TX, entity entities.Like) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewLikeDTOFromEntity(entity)
	q := sq.Insert("public.likes").
		Columns("id", "created_at", "updated_at", "deleted_at", "post_id", "value", "user_id").
		Values(dto.ID, dto.CreatedAt, dto.UpdatedAt, dto.DeletedAt, dto.PostId, dto.Value, dto.UserId)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if _, err := tx.GetSQLTx().ExecContext(ctx, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}
func (r *LikeRepository) Get(ctx context.Context, id uuid.UUID) (entities.Like, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &LikeDTO{}
	q := sq.Select("likes.id", "likes.created_at", "likes.updated_at", "likes.deleted_at", "likes.post_id", "likes.value", "likes.user_id").
		From("public.likes").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.readDB.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("like_id", id.String())
		return entities.Like{}, e
	}
	return dto.toEntity(), nil
}

func (r *LikeRepository) List(
	ctx context.Context,
	filter entities.LikeFilter,
) ([]entities.Like, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto LikeListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Of(pageSize)
	}
	q := sq.Select("likes.id", "likes.created_at", "likes.updated_at", "likes.deleted_at", "likes.post_id", "likes.value", "likes.user_id").
		From("public.likes").
		Limit(pageSize)
	if filter.IsDeleted != nil {
		if *filter.IsDeleted {
			q = q.Where(sq.NotEq{"likes.deleted_at": nil})
		} else {
			q = q.Where(sq.Eq{"likes.deleted_at": nil})
		}
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
func (r *LikeRepository) Count(ctx context.Context, filter entities.LikeFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.likes")
	if filter.IsDeleted != nil {
		if *filter.IsDeleted {
			q = q.Where(sq.NotEq{"likes.deleted_at": nil})
		} else {
			q = q.Where(sq.Eq{"likes.deleted_at": nil})
		}
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	var count uint64
	if err := r.readDB.GetContext(ctx, &count, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
func (r *LikeRepository) Update(ctx context.Context, tx dtx.TX, entity entities.Like) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewLikeDTOFromEntity(entity)
	q := sq.Update("public.likes").Where(sq.Eq{"id": entity.ID})
	{
		q = q.Set("created_at", dto.CreatedAt)
		q = q.Set("updated_at", dto.UpdatedAt)
		q = q.Set("deleted_at", dto.DeletedAt)
		q = q.Set("post_id", dto.PostId)
		q = q.Set("value", dto.Value)
		q = q.Set("user_id", dto.UserId)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("like_id", fmt.Sprint(entity.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("like_id", fmt.Sprint(entity.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("like_id", fmt.Sprint(entity.ID))
		return e
	}
	return nil
}
func (r *LikeRepository) Delete(ctx context.Context, tx dtx.TX, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.likes").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := tx.GetSQLTx().ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("like_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("like_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("like_id", fmt.Sprint(id))
		return e
	}
	return nil
}
