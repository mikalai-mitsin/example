package postgres

import (
	"context"
	"fmt"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type WidgetRepository struct {
	database *sqlx.DB
	logger   logger
}

func NewWidgetRepository(database *sqlx.DB, logger logger) *WidgetRepository {
	return &WidgetRepository{database: database, logger: logger}
}

type WidgetDTO struct {
	ID           string    `db:"id,omitempty"`
	UpdatedAt    time.Time `db:"updated_at,omitempty"`
	CreatedAt    time.Time `db:"created_at,omitempty"`
	FormScreenId string    `db:"form_screen_id"`
	Name         string    `db:"name"`
	Ordering     int64     `db:"ordering"`
	IsOptional   bool      `db:"is_optional"`
	UiSettings   string    `db:"ui_settings"`
	DeletedAt    time.Time `db:"deleted_at"`
}
type WidgetListDTO []*WidgetDTO

func (list WidgetListDTO) ToEntities() []*entities.Widget {
	items := make([]*entities.Widget, len(list))
	for i := range list {
		items[i] = list[i].ToModel()
	}
	return items
}
func NewWidgetDTOFromModel(model *entities.Widget) *WidgetDTO {
	dto := &WidgetDTO{
		ID:           string(model.ID),
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
		FormScreenId: model.FormScreenId,
		Name:         model.Name,
		Ordering:     model.Ordering,
		IsOptional:   model.IsOptional,
		UiSettings:   model.UiSettings,
		DeletedAt:    model.DeletedAt,
	}
	return dto
}
func (dto *WidgetDTO) ToModel() *entities.Widget {
	model := &entities.Widget{
		ID:           uuid.UUID(dto.ID),
		CreatedAt:    dto.CreatedAt,
		UpdatedAt:    dto.UpdatedAt,
		FormScreenId: dto.FormScreenId,
		Name:         dto.Name,
		Ordering:     dto.Ordering,
		IsOptional:   dto.IsOptional,
		UiSettings:   dto.UiSettings,
		DeletedAt:    dto.DeletedAt,
	}
	return model
}
func (r *WidgetRepository) Create(ctx context.Context, model *entities.Widget) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewWidgetDTOFromModel(model)
	q := sq.Insert("public.widgets").
		Columns("created_at", "updated_at", "form_screen_id", "name", "ordering", "is_optional", "ui_settings", "deleted_at").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.FormScreenId, dto.Name, dto.Ordering, dto.IsOptional, dto.UiSettings, dto.DeletedAt).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	model.ID = uuid.UUID(dto.ID)
	return nil
}
func (r *WidgetRepository) Get(ctx context.Context, id uuid.UUID) (*entities.Widget, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &WidgetDTO{}
	q := sq.Select("widgets.id", "widgets.created_at", "widgets.updated_at", "widgets.form_screen_id", "widgets.name", "widgets.ordering", "widgets.is_optional", "widgets.ui_settings", "widgets.deleted_at").
		From("public.widgets").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("widget_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}

func (r *WidgetRepository) List(
	ctx context.Context,
	filter *entities.WidgetFilter,
) ([]*entities.Widget, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto WidgetListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("widgets.id", "widgets.created_at", "widgets.updated_at", "widgets.form_screen_id", "widgets.name", "widgets.ordering", "widgets.is_optional", "widgets.ui_settings", "widgets.deleted_at").
		From("public.widgets").
		Limit(pageSize)
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
	if filter.PageNumber != nil && *filter.PageNumber > 1 {
		q = q.Offset((*filter.PageNumber - 1) * *filter.PageSize)
	}
	q = q.Limit(*filter.PageSize)
	if len(filter.OrderBy) > 0 {
		q = q.OrderBy(filter.OrderBy...)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.SelectContext(ctx, &dto, query, args...); err != nil {
		e := errs.FromPostgresError(err)
		return nil, e
	}
	return dto.ToEntities(), nil
}

func (r *WidgetRepository) Count(
	ctx context.Context,
	filter *entities.WidgetFilter,
) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.widgets")
	if len(filter.IDs) > 0 {
		q = q.Where(sq.Eq{"id": filter.IDs})
	}
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
func (r *WidgetRepository) Update(ctx context.Context, model *entities.Widget) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewWidgetDTOFromModel(model)
	q := sq.Update("public.widgets").Where(sq.Eq{"id": model.ID})
	{
		q = q.Set("widgets.created_at", dto.CreatedAt)
		q = q.Set("widgets.updated_at", dto.UpdatedAt)
		q = q.Set("widgets.form_screen_id", dto.FormScreenId)
		q = q.Set("widgets.name", dto.Name)
		q = q.Set("widgets.ordering", dto.Ordering)
		q = q.Set("widgets.is_optional", dto.IsOptional)
		q = q.Set("widgets.ui_settings", dto.UiSettings)
		q = q.Set("widgets.deleted_at", dto.DeletedAt)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("widget_id", fmt.Sprint(model.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("widget_id", fmt.Sprint(model.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("widget_id", fmt.Sprint(model.ID))
		return e
	}
	return nil
}
func (r *WidgetRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.widgets").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("widget_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("widget_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("widget_id", fmt.Sprint(id))
		return e
	}
	return nil
}
