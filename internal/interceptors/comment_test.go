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

func TestNewCommentInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase    usecases.AuthUseCase
		commentUseCase usecases.CommentUseCase
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  interceptors.CommentInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			want: &CommentInterceptor{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewCommentInterceptor(tt.args.commentUseCase, tt.args.authUseCase, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	type fields struct {
		authUseCase    usecases.AuthUseCase
		commentUseCase usecases.CommentUseCase
		logger         log.Logger
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
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDetail).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentDetail, comment).
					Return(nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDetail).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentDetail, comment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Comment not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDetail).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentInterceptor{
				commentUseCase: tt.fields.commentUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	create := mock_models.NewCommentCreate(t)
	type fields struct {
		commentUseCase usecases.CommentUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx         context.Context
		create      *models.CommentCreate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentCreate, create).
					Return(nil)
				commentUseCase.EXPECT().Create(ctx, create).Return(comment, nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				create:      create,
				requestUser: requestUser,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentCreate, create).
					Return(nil)
				commentUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
			i := &CommentInterceptor{
				commentUseCase: tt.fields.commentUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	update := mock_models.NewCommentUpdate(t)
	type fields struct {
		commentUseCase usecases.CommentUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx         context.Context
		update      *models.CommentUpdate
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentUpdate).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentUpdate, comment).
					Return(nil)
				commentUseCase.EXPECT().Update(ctx, update).Return(comment, nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				update:      update,
				requestUser: requestUser,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentUpdate).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentUpdate, comment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentUpdate).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentUpdate).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentUpdate, comment).
					Return(nil)
				commentUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
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
			i := &CommentInterceptor{
				commentUseCase: tt.fields.commentUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	type fields struct {
		commentUseCase usecases.CommentUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentDelete).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentDelete, comment).
					Return(nil)
				commentUseCase.EXPECT().
					Delete(ctx, comment.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDelete).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDelete).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentDelete, comment).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDelete).
					Return(nil)
				commentUseCase.EXPECT().
					Get(ctx, comment.ID).
					Return(comment, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentDelete, comment).
					Return(nil)
				commentUseCase.EXPECT().
					Delete(ctx, comment.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				commentUseCase: commentUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				id:          comment.ID,
				requestUser: requestUser,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &CommentInterceptor{
				commentUseCase: tt.fields.commentUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id, tt.args.requestUser); !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCommentInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_usecases.NewMockAuthUseCase(ctrl)
	requestUser := mock_models.NewUser(t)
	commentUseCase := mock_usecases.NewMockCommentUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewCommentFilter(t)
	count := uint64(faker.Number().NumberInt64(2))
	comments := make([]*models.Comment, 0, count)
	for i := uint64(0); i < count; i++ {
		comments = append(comments, mock_models.NewComment(t))
	}
	type fields struct {
		commentUseCase usecases.CommentUseCase
		authUseCase    usecases.AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx         context.Context
		filter      *models.CommentFilter
		requestUser *models.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Comment
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentList, filter).
					Return(nil)
				commentUseCase.EXPECT().
					List(ctx, filter).
					Return(comments, count, nil)
			},
			fields: fields{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			args: args{
				ctx:         ctx,
				filter:      filter,
				requestUser: requestUser,
			},
			want:    comments,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, models.PermissionIDCommentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
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
					HasPermission(ctx, requestUser, models.PermissionIDCommentList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, models.PermissionIDCommentList, filter).
					Return(nil)
				commentUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				commentUseCase: commentUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
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
			i := &CommentInterceptor{
				commentUseCase: tt.fields.commentUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter, tt.args.requestUser)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CommentInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
