package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	mockEntities "github.com/mikalai-mitsin/example/internal/app/comment/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewCommentService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	type args struct {
		commentRepository commentRepository
		clock             clock
		logger            logger
		uuid              uuidGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *CommentService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			want: &CommentService{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewCommentService(
				tt.args.commentRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	type fields struct {
		commentRepository commentRepository
		logger            logger
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
				mockCommentRepository.EXPECT().Get(ctx, comment.ID).Return(comment, nil)
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				mockCommentRepository.EXPECT().
					Get(ctx, comment.ID).
					Return(entities.Comment{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
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
			u := &CommentService{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listComments []entities.Comment
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listComments = append(listComments, mockEntities.NewComment(t))
	}
	filter := mockEntities.NewCommentFilter(t)
	type fields struct {
		commentRepository commentRepository
		logger            logger
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
				mockCommentRepository.EXPECT().List(ctx, filter).Return(listComments, nil)
				mockCommentRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
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
				mockCommentRepository.EXPECT().
					List(ctx, filter).
					Return([]entities.Comment{}, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "count error",
			setup: func() {
				mockCommentRepository.EXPECT().List(ctx, filter).Return(listComments, nil)
				mockCommentRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentService{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestCommentService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockuuidGenerator(ctrl)
	ctx := context.Background()
	create := mockEntities.NewCommentCreate(t)
	now := time.Now().UTC()
	type fields struct {
		commentRepository commentRepository
		clock             clock
		logger            logger
		uuid              uuidGenerator
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
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test"))
				mockCommentRepository.EXPECT().
					Create(
						ctx,
						entities.Comment{
							ID:        uuid.UUID("test"),
							Text:      create.Text,
							AuthorId:  create.AuthorId,
							PostId:    create.PostId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: entities.Comment{
				ID:        uuid.UUID("test"),
				Text:      create.Text,
				AuthorId:  create.AuthorId,
				PostId:    create.PostId,
				UpdatedAt: now,
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test 2"))
				mockCommentRepository.EXPECT().
					Create(
						ctx,
						entities.Comment{
							ID:        uuid.UUID("test 2"),
							Text:      create.Text,
							AuthorId:  create.AuthorId,
							PostId:    create.PostId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Comment{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
				clock:             mockClock,
				uuid:              mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: entities.CommentCreate{},
			},
			want: entities.Comment{},
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "text", Value: "cannot be blank"},
				errs.Param{Key: "author_id", Value: "cannot be blank"},
				errs.Param{Key: "post_id", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentService{
				commentRepository: tt.fields.commentRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
				uuid:              tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	mockClock := NewMockclock(ctrl)
	update := mockEntities.NewCommentUpdate(t)
	now := time.Now().UTC()
	updatedComment := entities.Comment{
		ID:        comment.ID,
		CreatedAt: comment.CreatedAt,
		UpdatedAt: now,

		Text:     *update.Text,
		AuthorId: *update.AuthorId,
		PostId:   *update.PostId,
	}
	type fields struct {
		commentRepository commentRepository
		clock             clock
		logger            logger
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
				mockClock.EXPECT().Now().Return(now)
				mockCommentRepository.EXPECT().
					Get(ctx, update.ID).Return(comment, nil)
				mockCommentRepository.EXPECT().
					Update(ctx, updatedComment).Return(nil)
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    updatedComment,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockCommentRepository.EXPECT().
					Get(ctx, update.ID).
					Return(comment, nil)
				mockCommentRepository.EXPECT().
					Update(ctx, updatedComment).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Comment{},
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Comment not found",
			setup: func() {
				mockCommentRepository.EXPECT().
					Get(ctx, update.ID).
					Return(entities.Comment{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Comment{},
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				clock:             mockClock,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				update: entities.CommentUpdate{
					ID: uuid.UUID("baduuid"),
				},
			},
			want:    entities.Comment{},
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentService{
				commentRepository: tt.fields.commentRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCommentService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockCommentRepository := NewMockcommentRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	comment := mockEntities.NewComment(t)
	type fields struct {
		commentRepository commentRepository
		logger            logger
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
				mockCommentRepository.EXPECT().
					Delete(ctx, comment.ID).
					Return(nil)
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				mockCommentRepository.EXPECT().
					Delete(ctx, comment.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				commentRepository: mockCommentRepository,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentService{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
