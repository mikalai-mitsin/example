package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/post/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewPostService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	type args struct {
		postRepository postRepository
		clock          clock
		logger         logger
		uuid           UUIDGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *PostService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			want: &PostService{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewPostService(
				tt.args.postRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := mock_entities.NewPost(t)
	type fields struct {
		postRepository postRepository
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
		want    *entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostRepository.EXPECT().Get(ctx, post.ID).Return(post, nil)
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "Post not found",
			setup: func() {
				mockPostRepository.EXPECT().
					Get(ctx, post.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostService{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listPosts []*entities.Post
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listPosts = append(listPosts, mock_entities.NewPost(t))
	}
	filter := mock_entities.NewPostFilter(t)
	type fields struct {
		postRepository postRepository
		logger         logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.Post
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostRepository.EXPECT().List(ctx, filter).Return(listPosts, nil)
				mockPostRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listPosts,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				mockPostRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
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
				mockPostRepository.EXPECT().List(ctx, filter).Return(listPosts, nil)
				mockPostRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
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
			u := &PostService{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestPostService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	ctx := context.Background()
	create := mock_entities.NewPostCreate(t)
	now := time.Now().UTC()
	type fields struct {
		postRepository postRepository
		clock          clock
		logger         logger
		uuid           UUIDGenerator
	}
	type args struct {
		ctx    context.Context
		create *entities.PostCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test"))
				mockPostRepository.EXPECT().
					Create(
						ctx,
						&entities.Post{
							ID:         uuid.UUID("test"),
							Title:      create.Title,
							Order:      create.Order,
							IsOptional: create.IsOptional,
							UpdatedAt:  now,
							CreatedAt:  now,
						},
					).
					Return(nil)
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &entities.Post{
				ID:         uuid.UUID("test"),
				Title:      create.Title,
				Order:      create.Order,
				IsOptional: create.IsOptional,
				UpdatedAt:  now,
				CreatedAt:  now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test 2"))
				mockPostRepository.EXPECT().
					Create(
						ctx,
						&entities.Post{
							ID:         uuid.UUID("test 2"),
							Title:      create.Title,
							Order:      create.Order,
							IsOptional: create.IsOptional,
							UpdatedAt:  now,
							CreatedAt:  now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
				clock:          mockClock,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: &entities.PostCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "title", Value: "cannot be blank"},
				errs.Param{Key: "order", Value: "cannot be blank"},
				errs.Param{Key: "is_optional", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostService{
				postRepository: tt.fields.postRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
				uuid:           tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := mock_entities.NewPost(t)
	mockClock := NewMockclock(ctrl)
	update := mock_entities.NewPostUpdate(t)
	now := post.UpdatedAt
	type fields struct {
		postRepository postRepository
		clock          clock
		logger         logger
	}
	type args struct {
		ctx    context.Context
		update *entities.PostUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockPostRepository.EXPECT().
					Get(ctx, update.ID).Return(post, nil)
				mockPostRepository.EXPECT().
					Update(ctx, post).Return(nil)
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockPostRepository.EXPECT().
					Get(ctx, update.ID).
					Return(post, nil)
				mockPostRepository.EXPECT().
					Update(ctx, post).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Post not found",
			setup: func() {
				mockPostRepository.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				postRepository: mockPostRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				update: &entities.PostUpdate{
					ID: uuid.UUID("baduuid"),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostService{
				postRepository: tt.fields.postRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostRepository := NewMockpostRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := mock_entities.NewPost(t)
	type fields struct {
		postRepository postRepository
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
				mockPostRepository.EXPECT().
					Delete(ctx, post.ID).
					Return(nil)
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: nil,
		},
		{
			name: "Post not found",
			setup: func() {
				mockPostRepository.EXPECT().
					Delete(ctx, post.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				postRepository: mockPostRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostService{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
