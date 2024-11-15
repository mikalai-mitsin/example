package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewUserRepository(t *testing.T) {
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
		want  *UserRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &UserRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewUserRepository(tt.args.database, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.users (created_at,updated_at,first_name,last_name,password,email,group_id) VALUES ($1,$2,$3,$4,$5,$6,$7) RETURNING id"
	user := mock_entities.NewUser(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.User
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
						user.UpdatedAt,
						user.CreatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
					).
					WillReturnRows(sqlmock.NewRows([]string{"id", "created_at"}).
						AddRow(user.ID, user.CreatedAt))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(
						user.UpdatedAt,
						user.CreatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUserRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT users.id, users.created_at, users.updated_at, users.first_name, users.last_name, users.password, users.email, users.group_id FROM public.users WHERE id = $1 LIMIT 1"
	user := mock_entities.NewUser(t)
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
		want    *entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newUserRows(t, []*entities.User{user})
				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			want: nil,
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(user.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError().WithParam("user_id", string(user.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_List(t *testing.T) {
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
	var listUsers []*entities.User
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listUsers = append(listUsers, mock_entities.NewUser(t))
	}
	filter := &entities.UserFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT users.id, users.created_at, users.updated_at, users.first_name, users.last_name, users.password, users.email, users.group_id FROM public.users ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newUserRows(t, listUsers))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listUsers,
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
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	user := mock_entities.NewUser(t)
	query := `UPDATE public.users SET users.created_at = $1, users.updated_at = $2, users.first_name = $3, users.last_name = $4, users.password = $5, users.email = $6, users.group_id = $7 WHERE id = $8`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx  context.Context
		card *entities.User
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
						user.CreatedAt,
						user.UpdatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						user.CreatedAt,
						user.UpdatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
						user.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("user_id", string(user.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						user.CreatedAt,
						user.UpdatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
						user.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						user.CreatedAt,
						user.UpdatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
						user.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						user.CreatedAt,
						user.UpdatedAt,
						user.FirstName,
						user.LastName,
						user.Password,
						user.Email,
						user.GroupID,
						user.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:  ctx,
				card: user,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.card)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	user := mock_entities.NewUser(t)
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
				mock.ExpectExec("DELETE FROM public.users WHERE id = $1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: nil,
		},
		{
			name: "article card not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = $1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("user_id", string(user.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = $1").
					WithArgs(user.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.users WHERE id = $1").
					WithArgs(user.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  user.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("user_id", string(user.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUserRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.users"
	ctx := context.Background()
	filter := &entities.UserFilter{}
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.UserFilter
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
			r := &UserRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newUserRows(t *testing.T, listUsers []*entities.User) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"first_name",
		"last_name",
		"password",
		"email",
		"group_id",
		"updated_at",
		"created_at",
	})
	for _, user := range listUsers {
		rows.AddRow(
			user.ID,
			user.FirstName,
			user.LastName,
			user.Password,
			user.Email,
			user.GroupID,
			user.UpdatedAt,
			user.CreatedAt,
		)
	}
	return rows
}
