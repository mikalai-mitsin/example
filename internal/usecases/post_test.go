package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	mock_repositories "github.com/018bf/example/internal/domain/repositories/mock"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	mock_clock "github.com/018bf/example/pkg/clock/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"
)

func TestNewPostUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		postRepository repositories.PostRepository
		clock          clock.Clock
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.PostUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
			want: &PostUseCase{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPostUseCase(tt.args.postRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	type fields struct {
		postRepository repositories.PostRepository
		logger         log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				postRepository.EXPECT().Get(ctx, post.ID).Return(post, nil)
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
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
				postRepository.EXPECT().Get(ctx, post.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostUseCase{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var posts []*models.Post
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		posts = append(posts, mock_models.NewPost(t))
	}
	filter := mock_models.NewPostFilter(t)
	type fields struct {
		postRepository repositories.PostRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.PostFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Post
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				postRepository.EXPECT().List(ctx, filter).Return(posts, nil)
				postRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    posts,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				postRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
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
				postRepository.EXPECT().List(ctx, filter).Return(posts, nil)
				postRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
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
			u := &PostUseCase{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PostUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestPostUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewPostCreate(t)
	now := time.Now().UTC()
	type fields struct {
		postRepository repositories.PostRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.PostCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				postRepository.EXPECT().
					Create(
						ctx,
						&models.Post{
							Body:      create.Body,
							Title:     create.Title,
							UserId:    create.UserId,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Post{
				ID:        "",
				Body:      create.Body,
				Title:     create.Title,
				UserId:    create.UserId,
				Weight:    create.Weight,
				UpdatedAt: now,
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				postRepository.EXPECT().
					Create(
						ctx,
						&models.Post{
							ID:        "",
							Body:      create.Body,
							Title:     create.Title,
							UserId:    create.UserId,
							Weight:    create.Weight,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		// TODO: Add validation rules or delete this case
		//{
		//	name: "invalid",
		//	setup: func() {
		//	},
		//	fields: fields{
		//		postRepository: postRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.PostCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostUseCase{
				postRepository: tt.fields.postRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewPostUpdate(t)
	now := post.UpdatedAt
	type fields struct {
		postRepository repositories.PostRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.PostUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Post
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				postRepository.EXPECT().
					Get(ctx, update.ID).Return(post, nil)
				postRepository.EXPECT().
					Update(ctx, post).Return(nil)
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
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
				clockMock.EXPECT().Now().Return(now)
				postRepository.EXPECT().
					Get(ctx, update.ID).
					Return(post, nil)
				postRepository.EXPECT().
					Update(ctx, post).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
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
				postRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				postRepository: postRepository,
				clock:          clockMock,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				update: &models.PostUpdate{
					ID: faker.Number().Number(1),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostUseCase{
				postRepository: tt.fields.postRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	postRepository := mock_repositories.NewMockPostRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	type fields struct {
		postRepository repositories.PostRepository
		logger         log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
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
				postRepository.EXPECT().
					Delete(ctx, post.ID).
					Return(nil)
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
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
				postRepository.EXPECT().
					Delete(ctx, post.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				postRepository: postRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  post.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &PostUseCase{
				postRepository: tt.fields.postRepository,
				logger:         tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("PostUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
