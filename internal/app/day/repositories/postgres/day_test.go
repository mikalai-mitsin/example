package postgres

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"testing"

	mock_models "github.com/018bf/example/internal/app/day/models/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"
	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/018bf/example/internal/pkg/postgres"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	"github.com/jmoiron/sqlx"
)

func TestNewDayRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		database *sqlx.DB
		logger   log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *DayRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &DayRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewDayRepository(tt.args.database, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewDayRepository() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "INSERT INTO public.days (created_at,updated_at,name,repeat,equipment_id) VALUES ($1,$2,$3,$4,$5) RETURNING id"
	day := mock_models.NewDay(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Day
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						day.UpdatedAt,
						day.CreatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(day.ID, day.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						day.UpdatedAt,
						day.CreatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Create(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("DayRepository.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	query := "SELECT days.id, days.created_at, days.updated_at, days.name, days.repeat, days.equipment_id FROM public.days WHERE id = $1 LIMIT 1"
	day := mock_models.NewDay(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newDayRows(t, []*models.Day{day})
				mock.ExpectQuery(query).WithArgs(day.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: ctx,
				id:  day.ID,
			},
			want:    day,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(day.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(day.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError().WithParam("day_id", string(day.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayRepository.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayRepository.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayRepository_List(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listDays []*models.Day
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listDays = append(listDays, mock_models.NewDay(t))
	}
	filter := &models.DayFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT days.id, days.created_at, days.updated_at, days.name, days.repeat, days.equipment_id FROM public.days ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.DayFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Day
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newDayRows(t, listDays))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listDays,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "error", Value: "test error"}},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("DayRepository.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DayRepository.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDayRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	day := mock_models.NewDay(t)
	query := `UPDATE public.days SET days.created_at = $1, days.updated_at = $2, days.name = $3, days.repeat = $4, days.equipment_id = $5 WHERE id = $6`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx  context.Context
		card *models.Day
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						day.CreatedAt,
						day.UpdatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
						day.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						day.CreatedAt,
						day.UpdatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
						day.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("day_id", string(day.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						day.CreatedAt,
						day.UpdatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
						day.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						day.CreatedAt,
						day.UpdatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
						day.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						day.CreatedAt,
						day.UpdatedAt,
						day.Name,
						day.Repeat,
						day.EquipmentID,
						day.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx:  ctx,
				card: day,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Update(tt.args.ctx, tt.args.card); !errors.Is(err, tt.wantErr) {
				t.Errorf("DayRepository.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	logger := mock_log.NewMockLogger(ctrl)
	day := mock_models.NewDay(t)
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx context.Context
		id  uuid.UUID
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "ok",
			fields: fields{
				database: db,
				logger:   logger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.days WHERE id = $1").
					WithArgs(day.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.days WHERE id = $1").
					WithArgs(day.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("day_id", string(day.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.days WHERE id = $1").
					WithArgs(day.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.days WHERE id = $1").
					WithArgs(day.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   logger,
			},
			args: args{
				ctx: context.Background(),
				id:  day.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("day_id", string(day.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("DayRepository.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDayRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.days"
	ctx := context.Background()
	filter := &models.DayFilter{}
	type fields struct {
		database *sqlx.DB
		logger   log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.DayFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow(1))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "bad return type",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).
						AddRow("one"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want: 0,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: errs.Params{
					{
						Key:   "error",
						Value: "sql: Scan error on column index 0, name \"count\": converting driver.Value type string (\"one\") to a uint64: invalid syntax",
					},
				},
			},
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    0,
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &DayRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Count() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Count() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func newDayRows(t *testing.T, listDays []*models.Day) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"name",
		"repeat",
		"equipment_id",
		"updated_at",
		"created_at",
	})
	for _, day := range listDays {
		rows.AddRow(
			day.ID,
			day.Name,
			day.Repeat,
			day.EquipmentID,
			day.UpdatedAt,
			day.CreatedAt,
		)
	}
	return rows
}
