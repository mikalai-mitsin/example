package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/postgres"
	"github.com/018bf/example/internal/pkg/uuid"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type DayRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewDayRepository(database *sqlx.DB, logger log.Logger) *DayRepository {
	return &DayRepository{database: database, logger: logger}
}

type DayDTO struct {
	ID          string    `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Name        string    `db:"name"`
	Repeat      int       `db:"repeat"`
	EquipmentID string    `db:"equipment_id"`
}
type DayListDTO []*DayDTO

func (list DayListDTO) ToModels() []*models.Day {
	items := make([]*models.Day, len(list))
	for i := range list {
		items[i] = list[i].ToModel()
	}
	return items
}
func NewDayDTOFromModel(model *models.Day) *DayDTO {
	dto := &DayDTO{
		ID:          string(model.ID),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Name:        model.Name,
		Repeat:      model.Repeat,
		EquipmentID: model.EquipmentID,
	}
	return dto
}
func (dto *DayDTO) ToModel() *models.Day {
	model := &models.Day{
		ID:          uuid.UUID(dto.ID),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		Name:        dto.Name,
		Repeat:      dto.Repeat,
		EquipmentID: dto.EquipmentID,
	}
	return model
}
func (r *DayRepository) Create(ctx context.Context, model *models.Day) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewDayDTOFromModel(model)
	q := sq.Insert("public.days").
		Columns("created_at", "updated_at", "name", "repeat", "equipment_id").
		Values(dto.CreatedAt, dto.UpdatedAt, dto.Name, dto.Repeat, dto.EquipmentID).
		Suffix("RETURNING id")
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.QueryRowxContext(ctx, query, args...).StructScan(dto); err != nil {
		e := errs.FromPostgresError(err)
		return e
	}
	model.ID = uuid.UUID(dto.ID)
	return nil
}
func (r *DayRepository) List(ctx context.Context, filter *models.DayFilter) ([]*models.Day, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto DayListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("days.id", "days.created_at", "days.updated_at", "days.name", "days.repeat", "days.equipment_id").
		From("public.days").
		Limit(pageSize)
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{Lang: "english", Query: *filter.Search, Fields: []string{"name"}},
		)
	}
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
	return dto.ToModels(), nil
}
func (r *DayRepository) Count(ctx context.Context, filter *models.DayFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.days")
	if filter.Search != nil {
		q = q.Where(
			postgres.Search{Lang: "english", Query: *filter.Search, Fields: []string{"name"}},
		)
	}
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
func (r *DayRepository) Get(ctx context.Context, id uuid.UUID) (*models.Day, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &DayDTO{}
	q := sq.Select("days.id", "days.created_at", "days.updated_at", "days.name", "days.repeat", "days.equipment_id").
		From("public.days").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("day_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}
func (r *DayRepository) Update(ctx context.Context, model *models.Day) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewDayDTOFromModel(model)
	q := sq.Update("public.days").Where(sq.Eq{"id": model.ID})
	{
		q = q.Set("days.created_at", dto.CreatedAt)
		q = q.Set("days.updated_at", dto.UpdatedAt)
		q = q.Set("days.name", dto.Name)
		q = q.Set("days.repeat", dto.Repeat)
		q = q.Set("days.equipment_id", dto.EquipmentID)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("day_id", fmt.Sprint(model.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("day_id", fmt.Sprint(model.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("day_id", fmt.Sprint(model.ID))
		return e
	}
	return nil
}
func (r *DayRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.days").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("day_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("day_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("day_id", fmt.Sprint(id))
		return e
	}
	return nil
}
