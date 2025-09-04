package repositories

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewLikeRepository(t *testing.T) {
	mockDB, _, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	type args struct {
		writeDB database
		readDB  database
		logger  logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *LikeRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				writeDB: mockDB,
				readDB:  mockDB,
			},
			want: &LikeRepository{
				writeDB: mockDB,
				readDB:  mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewLikeRepository(tt.args.readDB, tt.args.writeDB, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeRepository_Create(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockTxManager := dtx.NewManager(mockDB)
	mock.ExpectBegin()
	mockTX := mockTxManager.NewTx()
	query := "INSERT INTO public.likes (id,created_at,updated_at,post_id,value,user_id) VALUES ($1,$2,$3,$4,$5,$6)"
	like := entities.NewMockLike(t)
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx  context.Context
		tx   dtx.TX
		like entities.Like
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
						like.ID,
						like.UpdatedAt,
						like.CreatedAt,
						like.PostId,
						like.Value,
						like.UserId,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						like.ID,
						like.UpdatedAt,
						like.CreatedAt,
						like.PostId,
						like.Value,
						like.UserId,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.tx, tt.args.like)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeRepository_Get(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT likes.id, likes.created_at, likes.updated_at, likes.post_id, likes.value, likes.user_id FROM public.likes WHERE id = $1 LIMIT 1"
	like := entities.NewMockLike(t)
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
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
		want    entities.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newLikeRows(t, []entities.Like{like})
				mock.ExpectQuery(query).WithArgs(like.ID).WillReturnRows(rows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  like.ID,
			},
			want:    like,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(like.ID).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  like.ID,
			},
			want: entities.Like{},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(like.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  like.ID,
			},
			want:    entities.Like{},
			wantErr: errs.NewEntityNotFoundError().WithParam("like_id", like.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeRepository_List(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listLikes []entities.Like
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listLikes = append(listLikes, entities.NewMockLike(t))
	}
	filter := entities.LikeFilter{
		PageSize:   pointer.Of(uint64(10)),
		PageNumber: pointer.Of(uint64(2)),
		Search:     nil,
		OrderBy:    []entities.LikeOrdering{"id"},
	}
	query := "SELECT likes.id, likes.created_at, likes.updated_at, likes.post_id, likes.value, likes.user_id FROM public.likes ORDER BY likes.id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.LikeFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Like
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newLikeRows(t, listLikes))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listLikes,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
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
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
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
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLikeRepository_Update(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockTxManager := dtx.NewManager(mockDB)
	mock.ExpectBegin()
	mockTX := mockTxManager.NewTx()
	like := entities.NewMockLike(t)
	query := `UPDATE public.likes SET created_at = $1, updated_at = $2, post_id = $3, value = $4, user_id = $5 WHERE id = $6`
	ctx := context.Background()
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx  context.Context
		tx   dtx.TX
		like entities.Like
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
						like.CreatedAt,
						like.UpdatedAt,
						like.PostId,
						like.Value,
						like.UserId,
						like.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						like.CreatedAt,
						like.UpdatedAt,
						like.PostId,
						like.Value,
						like.UserId,
						like.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("like_id", like.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						like.CreatedAt,
						like.UpdatedAt,
						like.PostId,
						like.Value,
						like.UserId,
						like.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						like.CreatedAt,
						like.UpdatedAt,
						like.PostId,
						like.Value,
						like.UserId,
						like.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						like.CreatedAt,
						like.UpdatedAt,
						like.PostId,
						like.Value,
						like.UserId,
						like.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx:  ctx,
				tx:   mockTX,
				like: like,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.tx, tt.args.like)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeRepository_Delete(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	mockTxManager := dtx.NewManager(mockDB)
	mock.ExpectBegin()
	mockTX := mockTxManager.NewTx()
	like := entities.NewMockLike(t)
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx context.Context
		tx  dtx.TX
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
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			setup: func() {
				mock.ExpectExec("DELETE FROM public.likes WHERE id = $1").
					WithArgs(like.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				tx:  mockTX,
				id:  like.ID,
			},
			wantErr: nil,
		},
		{
			name: "like not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.likes WHERE id = $1").
					WithArgs(like.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				tx:  mockTX,
				id:  like.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("like_id", like.ID.String()),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.likes WHERE id = $1").
					WithArgs(like.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				tx:  mockTX,
				id:  like.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.likes WHERE id = $1").
					WithArgs(like.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				writeDB: mockDB,
				readDB:  mockDB,
				logger:  mockLogger,
			},
			args: args{
				ctx: context.Background(),
				tx:  mockTX,
				id:  like.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("like_id", like.ID.String()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.tx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestLikeRepository_Count(t *testing.T) {
	mockDB, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer mockDB.Close()
	query := "SELECT count(id) FROM public.likes"
	ctx := context.Background()
	filter := entities.LikeFilter{}
	type fields struct {
		writeDB database
		readDB  database
		logger  logger
	}
	type args struct {
		ctx    context.Context
		filter entities.LikeFilter
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
				writeDB: mockDB,
				readDB:  mockDB,
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
				writeDB: mockDB,
				readDB:  mockDB,
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
				writeDB: mockDB,
				readDB:  mockDB,
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
			r := &LikeRepository{
				writeDB: tt.fields.writeDB,
				readDB:  tt.fields.readDB,
				logger:  tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newLikeRows(t *testing.T, listLikes []entities.Like) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"post_id",
		"value",
		"user_id",
		"updated_at",
		"created_at",
	})
	for _, like := range listLikes {
		rows.AddRow(
			like.ID,
			like.PostId,
			like.Value,
			like.UserId,
			like.UpdatedAt,
			like.CreatedAt,
		)
	}
	return rows
}
