package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	mockEntities "github.com/mikalai-mitsin/example/internal/app/comment/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewCommentUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		commentService commentService
		logger         logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *CommentUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			want: &CommentUseCase{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewCommentUseCase(tt.args.commentService, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	type fields struct {
		commentService commentService
		logger         logger
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
				mockCommentService.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(comment.ID),
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				mockCommentService.EXPECT().
					Get(ctx, comment.ID).
					Return(entities.Comment{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			want:    entities.Comment{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentUseCase{
				commentService: tt.fields.commentService,
				logger:         tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	create := mockEntities.NewCommentCreate(t)
	type fields struct {
		commentService commentService
		logger         logger
	}
	type args struct {
		ctx    context.Context
		create entities.CommentCreate
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
				mockCommentService.EXPECT().Create(ctx, create).Return(comment, nil)
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockCommentService.EXPECT().
					Create(ctx, create).
					Return(entities.Comment{}, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Comment{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentUseCase{
				commentService: tt.fields.commentService,
				logger:         tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	update := mockEntities.NewCommentUpdate(t)
	type fields struct {
		commentService commentService
		logger         logger
	}
	type args struct {
		ctx    context.Context
		update entities.CommentUpdate
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
				mockCommentService.EXPECT().Update(ctx, update).Return(comment, nil)
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockCommentService.EXPECT().
					Update(ctx, update).
					Return(entities.Comment{}, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Comment{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentUseCase{
				commentService: tt.fields.commentService,
				logger:         tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	type fields struct {
		commentService commentService
		logger         logger
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
			setup: func() {
				mockCommentService.EXPECT().
					Delete(ctx, comment.ID).
					Return(nil)
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockCommentService.EXPECT().
					Delete(ctx, comment.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentUseCase{
				commentService: tt.fields.commentService,
				logger:         tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestCommentUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentService := NewMockcommentService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mockEntities.NewCommentFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listComments := make([]entities.Comment, 0, count)
	for i := uint64(0); i < count; i++ {
		listComments = append(listComments, mockEntities.NewComment(t))
	}
	type fields struct {
		commentService commentService
		logger         logger
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
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockCommentService.EXPECT().
					List(ctx, filter).
					Return(listComments, count, nil)
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listComments,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockCommentService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				commentService: mockCommentService,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentUseCase{
				commentService: tt.fields.commentService,
				logger:         tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
