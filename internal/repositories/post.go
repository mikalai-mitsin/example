package repositories

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/018bf/example/pkg/log"

	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/repositories"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/jmoiron/sqlx"
)

type PostRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPostRepository(
	database *sqlx.DB,
	logger log.Logger,
) repositories.PostRepository {
	return &PostRepository{
		database: database,
		logger:   logger,
	}
}

func (r *PostRepository) Create(
	ctx context.Context,
	post *models.Post,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Insert("public.posts").
		Columns(
			"body",
			"title",
			"user_id",
			"weight",
			"updated_at",
			"created_at",
		).
		Values(
			post.Body,
			post.Title,
			post.UserId,
			post.Weight,
			post.UpdatedAt,
			post.CreatedAt,
		).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(post); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	return nil
}

func (r *PostRepository) Get(
	ctx context.Context,
	id string,
) (*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	post := &models.Post{}
	q := sq.Select(
		"posts.id",
		"posts.body",
		"posts.title",
		"posts.user_id",
		"posts.weight",
		"posts.updated_at",
		"posts.created_at",
	).
		From("public.posts").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, post, query, args...); err != nil {
		e := errs.FromPostgresError(err).
			WithParam("post_id", id)
		return nil, e
	}
	return post, nil
}

func (r *PostRepository) List(
	ctx context.Context,
	filter *models.PostFilter,
) ([]*models.Post, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var posts []*models.Post
	const pageSize = 10
	q := sq.Select(
		"posts.id",
		"posts.body",
		"posts.title",
		"posts.user_id",
		"posts.weight",
		"posts.updated_at",
		"posts.created_at",
	).
		From("public.posts").
		Limit(pageSize)
	// TODO: add filtering
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	if filter.PageSize != nil {
		q = q.Limit(*filter.PageSize)
	}
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &posts, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return posts, nil
}

func (r *PostRepository) Update(
	ctx context.Context,
	post *models.Post,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Update("public.posts").
		Where(sq.Eq{"id": post.ID}).
		Set("posts.body", post.Body).
		Set("posts.title", post.Title).
		Set("posts.user_id", post.UserId).
		Set("posts.weight", post.Weight).
		Set("updated_at", post.UpdatedAt)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("post_id", fmt.Sprint(post.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).
			WithParam("post_id", fmt.Sprint(post.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("post_id", fmt.Sprint(post.ID))
		return e
	}
	return nil
}

func (r *PostRepository) Delete(
	ctx context.Context,
	id string,
) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.posts").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("post_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).
			WithParam("post_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFound().
			WithParam("post_id", fmt.Sprint(id))
		return e
	}
	return nil
}

func (r *PostRepository) Count(
	ctx context.Context,
	filter *models.PostFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.posts")
	// TODO: add filtering
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result := r.database.QueryRowxContext(ctx, query, args...)
	if err := result.Err(); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	var count uint64
	if err := result.Scan(&count); err != nil {
		e := errs.FromPostgresError(err)
		return 0, e
	}
	return count, nil
}
