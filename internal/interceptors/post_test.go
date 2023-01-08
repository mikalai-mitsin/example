package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/018bf/example/internal/domain/errs"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	mock_usecases "github.com/018bf/example/internal/domain/usecases/mock"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/log"
)

func TestNewPostInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase usecases.AuthUseCase
		postUseCase usecases.PostUseCase
		logger      log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.PostInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			want: &PostInterceptor{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewPostInterceptor(tt.args.postUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	type fields struct {
		authUseCase usecases.AuthUseCase
		postUseCase usecases.PostUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		id          string
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDetail).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostDetail, post).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDetail).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostDetail, post).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Post not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDetail).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostInterceptor{
				postUseCase: tt.fields.postUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	create := mock_models.NewPostCreate(t)
	type fields struct {
		postUseCase usecases.PostUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.PostCreate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostCreate, create).
					Return(nil)
				postUseCase.EXPECT().Create(ctx, create).Return(post, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostCreate, create).
					Return(nil)
				postUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostInterceptor{
				postUseCase: tt.fields.postUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	update := mock_models.NewPostUpdate(t)
	type fields struct {
		postUseCase usecases.PostUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.PostUpdate
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostUpdate).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostUpdate, post).
					Return(nil)
				postUseCase.EXPECT().Update(ctx, update).Return(post, nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    post,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostUpdate).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostUpdate, post).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostUpdate).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostUpdate).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostUpdate, post).
					Return(nil)
				postUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostInterceptor{
				postUseCase: tt.fields.postUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPostInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	post := mock_models.NewPost(t)
	type fields struct {
		postUseCase usecases.PostUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		id          string
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDelete).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostDelete, post).
					Return(nil)
				postUseCase.EXPECT().
					Delete(ctx, post.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Post not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDelete).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDelete).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostDelete, post).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDelete).
					Return(nil)
				postUseCase.EXPECT().
					Get(ctx, post.ID).
					Return(post, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostDelete, post).
					Return(nil)
				postUseCase.EXPECT().
					Delete(ctx, post.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase: authUseCase,
				postUseCase: postUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				id:          post.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostInterceptor{
				postUseCase: tt.fields.postUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("PostInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPostInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	postUseCase := mock_usecases.NewMockPostUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewPostFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	posts := make([]*models.Post, 0, count)
	for i := uint64(0); i < count; i++ {
		posts = append(posts, mock_models.NewPost(t))
	}
	type fields struct {
		postUseCase usecases.PostUseCase
		authUseCase usecases.AuthUseCase
		logger      log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.PostFilter
		requestUser *models.User
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
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostList, filter).
					Return(nil)
				postUseCase.EXPECT().
					List(ctx, filter).
					Return(posts, count, nil)
			},
			fields: fields{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    posts,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDPostList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDPostList, filter).
					Return(nil)
				postUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				postUseCase: postUseCase,
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewUnexpectedBehaviorError("l e"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &PostInterceptor{
				postUseCase: tt.fields.postUseCase,
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("PostInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PostInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PostInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
