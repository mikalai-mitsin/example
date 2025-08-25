package usecases

import (
	"context"
	"testing"

	"github.com/jaswdr/faker"
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewPostUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		postService       postService
		postEventProducer postEventProducer
		logger            logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *PostUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			want: &PostUseCase{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewPostUseCase(tt.args.postService, tt.args.postEventProducer, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		postService       postService
		postEventProducer postEventProducer
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
		want    entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostService.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(post.ID),
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "Post not found",
			setup: func() {
				mockPostService.EXPECT().
					Get(ctx, post.ID).
					Return(entities.Post{}, errs.NewEntityNotFoundError())
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    entities.Post{},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:       tt.fields.postService,
				postEventProducer: tt.fields.postEventProducer,
				logger:            tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	create := entities.NewMockPostCreate(t)
	type fields struct {
		postService       postService
		postEventProducer postEventProducer
		logger            logger
	}
	type args struct {
		ctx    context.Context
		create entities.PostCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostService.EXPECT().Create(ctx, create).Return(post, nil)
				mockpostEventProducer.EXPECT().Created(ctx, post).Return(nil)
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				mockPostService.EXPECT().
					Create(ctx, create).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    entities.Post{},
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:       tt.fields.postService,
				postEventProducer: tt.fields.postEventProducer,
				logger:            tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	update := entities.NewMockPostUpdate(t)
	type fields struct {
		postService       postService
		postEventProducer postEventProducer
		logger            logger
	}
	type args struct {
		ctx    context.Context
		update entities.PostUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    entities.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostService.EXPECT().Update(ctx, update).Return(post, nil)
				mockpostEventProducer.EXPECT().Updated(ctx, post).Return(nil)
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
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
				mockPostService.EXPECT().
					Update(ctx, update).
					Return(entities.Post{}, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    entities.Post{},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:       tt.fields.postService,
				postEventProducer: tt.fields.postEventProducer,
				logger:            tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestPostUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	post := entities.NewMockPost(t)
	type fields struct {
		postService       postService
		postEventProducer postEventProducer
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
				mockPostService.EXPECT().
					Delete(ctx, post.ID).
					Return(nil)
				mockpostEventProducer.EXPECT().Deleted(ctx, post.ID).Return(nil)
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				mockPostService.EXPECT().
					Delete(ctx, post.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostUseCase{
				postService:       tt.fields.postService,
				postEventProducer: tt.fields.postEventProducer,
				logger:            tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestPostUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockPostService := NewMockpostService(ctrl)
	mockpostEventProducer := NewMockpostEventProducer(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := entities.NewMockPostFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listPosts := make([]entities.Post, 0, count)
	for i := uint64(0); i < count; i++ {
		listPosts = append(listPosts, entities.NewMockPost(t))
	}
	type fields struct {
		postService       postService
		postEventProducer postEventProducer
		logger            logger
	}
	type args struct {
		ctx    context.Context
		filter entities.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []entities.Post
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockPostService.EXPECT().
					List(ctx, filter).
					Return(listPosts, count, nil)
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
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
				mockPostService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				postService:       mockPostService,
				postEventProducer: mockpostEventProducer,
				logger:            mockLogger,
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
			i := &PostUseCase{
				postService:       tt.fields.postService,
				postEventProducer: tt.fields.postEventProducer,
				logger:            tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
