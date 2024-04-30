package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/postgres"
	"github.com/018bf/example/internal/pkg/uuid"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

type PlanRepository struct {
	database *sqlx.DB
	logger   log.Logger
}

func NewPlanRepository(database *sqlx.DB, logger log.Logger) *PlanRepository {
	return &PlanRepository{database: database, logger: logger}
}

type PlanDTO struct {
	ID          string    `db:"id,omitempty"`
	UpdatedAt   time.Time `db:"updated_at,omitempty"`
	CreatedAt   time.Time `db:"created_at,omitempty"`
	Name        string    `db:"name"`
	Repeat      int64     `db:"repeat"`
	EquipmentID string    `db:"equipment_id"`
}
type PlanListDTO []*PlanDTO

func (list PlanListDTO) ToModels() []*models.Plan {
	items := make([]*models.Plan, len(list))
	for i := range list {
		items[i] = list[i].ToModel()
	}
	return items
}
func NewPlanDTOFromModel(model *models.Plan) *PlanDTO {
	dto := &PlanDTO{
		ID:          string(model.ID),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		Name:        model.Name,
		Repeat:      int64(model.Repeat),
		EquipmentID: model.EquipmentID,
	}
	return dto
}
func (dto *PlanDTO) ToModel() *models.Plan {
	model := &models.Plan{
		ID:          uuid.UUID(dto.ID),
		CreatedAt:   dto.CreatedAt,
		UpdatedAt:   dto.UpdatedAt,
		Name:        dto.Name,
		Repeat:      uint64(dto.Repeat),
		EquipmentID: dto.EquipmentID,
	}
	return model
}
func (r *PlanRepository) Create(ctx context.Context, model *models.Plan) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPlanDTOFromModel(model)
	q := sq.Insert("public.plans").
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

func (r *PlanRepository) List(
	ctx context.Context,
	filter *models.PlanFilter,
) ([]*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	var dto PlanListDTO
	const pageSize = uint64(10)
	if filter.PageSize == nil {
		filter.PageSize = pointer.Pointer(pageSize)
	}
	q := sq.Select("plans.id", "plans.created_at", "plans.updated_at", "plans.name", "plans.repeat", "plans.equipment_id").
		From("public.plans").
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
func (r *PlanRepository) Count(ctx context.Context, filter *models.PlanFilter) (uint64, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Select("count(id)").From("public.plans")
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
func (r *PlanRepository) Get(ctx context.Context, id uuid.UUID) (*models.Plan, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := &PlanDTO{}
	q := sq.Select("plans.id", "plans.created_at", "plans.updated_at", "plans.name", "plans.repeat", "plans.equipment_id").
		From("public.plans").
		Where(sq.Eq{"id": id}).
		Limit(1)
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	if err := r.database.GetContext(ctx, dto, query, args...); err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", string(id))
		return nil, e
	}
	return dto.ToModel(), nil
}
func (r *PlanRepository) Update(ctx context.Context, model *models.Plan) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	dto := NewPlanDTOFromModel(model)
	q := sq.Update("public.plans").Where(sq.Eq{"id": model.ID})
	{
		q = q.Set("plans.created_at", dto.CreatedAt)
		q = q.Set("plans.updated_at", dto.UpdatedAt)
		q = q.Set("plans.name", dto.Name)
		q = q.Set("plans.repeat", dto.Repeat)
		q = q.Set("plans.equipment_id", dto.EquipmentID)
	}
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(model.ID))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(model.ID))
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("plan_id", fmt.Sprint(model.ID))
		return e
	}
	return nil
}
func (r *PlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	q := sq.Delete("public.plans").Where(sq.Eq{"id": id})
	query, args := q.PlaceholderFormat(sq.Dollar).MustSql()
	result, err := r.database.ExecContext(ctx, query, args...)
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	affected, err := result.RowsAffected()
	if err != nil {
		e := errs.FromPostgresError(err).WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	if affected == 0 {
		e := errs.NewEntityNotFoundError().WithParam("plan_id", fmt.Sprint(id))
		return e
	}
	return nil
}
