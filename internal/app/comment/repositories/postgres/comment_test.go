package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jaswdr/faker"
	mockEntities "github.com/mikalai-mitsin/example/internal/app/comment/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/postgres"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/jmoiron/sqlx"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewCommentRepository(t *testing.T) {
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
		want  *CommentRepository
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				database: mockDB,
			},
			want: &CommentRepository{
				database: mockDB,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewCommentRepository(tt.args.database, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentRepository_Create(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "INSERT INTO public.comments (created_at,updated_at,text,author_id,post_id) VALUES ($1,$2,$3,$4,$5)"
	comment := mockEntities.NewComment(t)
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx     context.Context
		comment entities.Comment
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
						comment.UpdatedAt,
						comment.CreatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: nil,
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						comment.UpdatedAt,
						comment.CreatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Create(tt.args.ctx, tt.args.comment)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommentRepository_Get(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	query := "SELECT comments.id, comments.created_at, comments.updated_at, comments.text, comments.author_id, comments.post_id FROM public.comments WHERE id = $1 LIMIT 1"
	comment := mockEntities.NewComment(t)
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
		want    entities.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				rows := newCommentRows(t, []entities.Comment{comment})
				mock.ExpectQuery(query).WithArgs(comment.ID).WillReturnRows(rows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mock.ExpectQuery(query).
					WithArgs(comment.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			want: entities.Comment{},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectQuery(query).WithArgs(comment.ID).WillReturnError(sql.ErrNoRows)
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			want:    entities.Comment{},
			wantErr: errs.NewEntityNotFoundError().WithParam("comment_id", string(comment.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentRepository_List(t *testing.T) {
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
	var listComments []entities.Comment
	for i := 0; i < faker.New().IntBetween(2, 20); i++ {
		listComments = append(listComments, mockEntities.NewComment(t))
	}
	filter := entities.CommentFilter{
		PageSize:   pointer.Pointer(uint64(10)),
		PageNumber: pointer.Pointer(uint64(2)),
		Search:     nil,
		OrderBy:    []string{"id ASC"},
		IDs:        nil,
	}
	query := "SELECT comments.id, comments.created_at, comments.updated_at, comments.text, comments.author_id, comments.post_id FROM public.comments ORDER BY id ASC LIMIT 10 OFFSET 10"
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter entities.CommentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mock.ExpectQuery(query).
					WillReturnRows(newCommentRows(t, listComments))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listComments,
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
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentRepository_Update(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	comment := mockEntities.NewComment(t)
	query := `UPDATE public.comments SET created_at = $1, updated_at = $2, text = $3, author_id = $4, post_id = $5 WHERE id = $6`
	ctx := context.Background()
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx     context.Context
		comment entities.Comment
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
						comment.CreatedAt,
						comment.UpdatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
						comment.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: nil,
		},
		{
			name: "not found",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						comment.CreatedAt,
						comment.UpdatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
						comment.ID,
					).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						comment.CreatedAt,
						comment.UpdatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
						comment.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "unexpected error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						comment.CreatedAt,
						comment.UpdatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
						comment.ID,
					).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec(query).
					WithArgs(
						comment.CreatedAt,
						comment.UpdatedAt,
						comment.Text,
						comment.AuthorId,
						comment.PostId,
						comment.ID,
					).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx:     ctx,
				comment: comment,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Update(tt.args.ctx, tt.args.comment)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommentRepository_Delete(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockLogger := NewMocklogger(ctrl)
	comment := mockEntities.NewComment(t)
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
				mock.ExpectExec("DELETE FROM public.comments WHERE id = $1").
					WithArgs(comment.ID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			wantErr: nil,
		},
		{
			name: "comment not found",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.comments WHERE id = $1").
					WithArgs(comment.ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			wantErr: errs.NewEntityNotFoundError().WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "database error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.comments WHERE id = $1").
					WithArgs(comment.ID).
					WillReturnError(errors.New("test error"))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
		{
			name: "result error",
			setup: func() {
				mock.ExpectExec("DELETE FROM public.comments WHERE id = $1").
					WithArgs(comment.ID).
					WillReturnResult(sqlmock.NewErrorResult(errors.New("test error")))
			},
			fields: fields{
				database: db,
				logger:   mockLogger,
			},
			args: args{
				ctx: context.Background(),
				id:  comment.ID,
			},
			wantErr: errs.FromPostgresError(errors.New("test error")).
				WithParam("comment_id", string(comment.ID)),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			err := r.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommentRepository_Count(t *testing.T) {
	db, mock, err := postgres.NewMockPostgreSQL(t)
	if err != nil {
		t.Fatal(err)
		return
	}
	defer db.Close()
	query := "SELECT count(id) FROM public.comments"
	ctx := context.Background()
	filter := entities.CommentFilter{}
	type fields struct {
		database *sqlx.DB
		logger   logger
	}
	type args struct {
		ctx    context.Context
		filter entities.CommentFilter
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
			r := &CommentRepository{
				database: tt.fields.database,
				logger:   tt.fields.logger,
			}
			got, err := r.Count(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func newCommentRows(t *testing.T, listComments []entities.Comment) *sqlmock.Rows {
	t.Helper()
	rows := sqlmock.NewRows([]string{
		"id",
		"text",
		"author_id",
		"post_id",
		"updated_at",
		"created_at",
	})
	for _, comment := range listComments {
		rows.AddRow(
			comment.ID,
			comment.Text,
			comment.AuthorId,
			comment.PostId,
			comment.UpdatedAt,
			comment.CreatedAt,
		)
	}
	return rows
}
