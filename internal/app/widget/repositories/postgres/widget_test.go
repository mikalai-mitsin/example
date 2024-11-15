package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/widget/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewWidgetRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		database *sqlx.DB
		logger   logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *WidgetRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &WidgetRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewWidgetRepository(tt.args.database, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.widgets (created_at,updated_at,form_screen_id,name,ordering,is_optional,ui_settings,deleted_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING id"
	widget := mock_entities.NewWidget(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.Widget
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
						widget.UpdatedAt,
						widget.CreatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(widget.ID, widget.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						widget.UpdatedAt,
						widget.CreatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWidgetRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT widgets.id, widgets.created_at, widgets.updated_at, widgets.form_screen_id, widgets.name, widgets.ordering, widgets.is_optional, widgets.ui_settings, widgets.deleted_at FROM public.widgets WHERE id = $1 LIMIT 1"
	widget := mock_entities.NewWidget(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
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
		want    *entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newWidgetRows(t, []*entities.Widget{widget})
				mock.ExpectQuery(query).WithArgs(widget.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  widget.ID,
			},
			want:    widget,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(widget.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(widget.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError().WithParam("widget_id", string(widget.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetRepository_List(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listWidgets []*entities.Widget
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listWidgets = append(listWidgets, mock_entities.NewWidget(t))
	}
	filter := &entities.WidgetFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT widgets.id, widgets.created_at, widgets.updated_at, widgets.form_screen_id, widgets.name, widgets.ordering, widgets.is_optional, widgets.ui_settings, widgets.deleted_at FROM public.widgets ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.WidgetFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.Widget
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newWidgetRows(t, listWidgets))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listWidgets,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
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
				logger:   mockLogger,
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
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestWidgetRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	widget := mock_entities.NewWidget(t)
	query := `UPDATE public.widgets SET widgets.created_at = $1, widgets.updated_at = $2, widgets.form_screen_id = $3, widgets.name = $4, widgets.ordering = $5, widgets.is_optional = $6, widgets.ui_settings = $7, widgets.deleted_at = $8 WHERE id = $9`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.Widget
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
						widget.CreatedAt,
						widget.UpdatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
						widget.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						widget.CreatedAt,
						widget.UpdatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
						widget.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						widget.CreatedAt,
						widget.UpdatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
						widget.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						widget.CreatedAt,
						widget.UpdatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
						widget.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						widget.CreatedAt,
						widget.UpdatedAt,
						widget.FormScreenId,
						widget.Name,
						widget.Ordering,
						widget.IsOptional,
						widget.UiSettings,
						widget.DeletedAt,
						widget.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: widget,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWidgetRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	widget := mock_entities.NewWidget(t)
	type fields struct {
		database *sqlx.DB
		logger   logger
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
				logger:   mockLogger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.widgets WHERE id = $1").
					WithArgs(widget.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.widgets WHERE id = $1").
					WithArgs(widget.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.widgets WHERE id = $1").
					WithArgs(widget.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.widgets WHERE id = $1").
					WithArgs(widget.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  widget.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("widget_id", string(widget.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestWidgetRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.widgets"
	ctx := context.Background()
	filter := &entities.WidgetFilter{}
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.WidgetFilter
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
			r := &WidgetRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newWidgetRows(t *testing.T, listWidgets []*entities.Widget) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"form_screen_id",
		"name",
		"ordering",
		"is_optional",
		"ui_settings",
		"deleted_at",
		"updated_at",
		"created_at",
	})
	for _, widget := range listWidgets {
		rows.AddRow(
			widget.ID,
			widget.FormScreenId,
			widget.Name,
			widget.Ordering,
			widget.IsOptional,
			widget.UiSettings,
			widget.DeletedAt,
			widget.UpdatedAt,
			widget.CreatedAt,
		)
	}
	return rows
}
