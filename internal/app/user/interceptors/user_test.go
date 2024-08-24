package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/mikalai-mitsin/example/internal/app/user/interceptors/mock"
	"github.com/mikalai-mitsin/example/internal/app/user/models"
	mock_models "github.com/mikalai-mitsin/example/internal/app/user/models/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewUserInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	type args struct {
		userUseCase UserUseCase
		logger      Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *UserInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				userUseCase: userUseCase,
				logger:      logger,
			},
			want: &UserInterceptor{
				userUseCase: userUseCase,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewUserInterceptor(tt.args.userUseCase, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUserInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userUseCase UserUseCase
		logger      Logger
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
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userUseCase.EXPECT().
					Get(ctx, user.ID).
					Return(user, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(user.ID),
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				userUseCase.EXPECT().
					Get(ctx, user.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	create := mock_models.NewUserCreate(t)
	type fields struct {
		userUseCase UserUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		create *models.UserCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userUseCase.EXPECT().Create(ctx, create).Return(user, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "create error",
			setup: func() {
				userUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("c u"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	update := mock_models.NewUserUpdate(t)
	type fields struct {
		userUseCase UserUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		update *models.UserUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userUseCase.EXPECT().Update(ctx, update).Return(user, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				userUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		userUseCase UserUseCase
		logger      Logger
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
				userUseCase.EXPECT().
					Delete(ctx, user.ID).
					Return(nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: nil,
		},
		{
			name: "delete error",
			setup: func() {
				userUseCase.EXPECT().
					Delete(ctx, user.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUserInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userUseCase := mock_interceptors.NewMockUserUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewUserFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listUsers := make([]*models.User, 0, count)
	for i := uint64(0); i < count; i++ {
		listUsers = append(listUsers, mock_models.NewUser(t))
	}
	type fields struct {
		userUseCase UserUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.User
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userUseCase.EXPECT().
					List(ctx, filter).
					Return(listUsers, count, nil)
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listUsers,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				userUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				userUseCase: userUseCase,
				logger:      logger,
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
			i := &UserInterceptor{
				userUseCase: tt.fields.userUseCase,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("UserInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("UserInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("UserInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
