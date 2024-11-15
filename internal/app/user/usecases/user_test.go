package usecases

import (
	"context"
	"testing"

	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"github.com/jaswdr/faker"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func TestNewUserUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		userService userService
		logger      logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *UserUseCase
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				userService: mockUserService,
				logger:      mockLogger,
			},
			want: &UserUseCase{
				userService: mockUserService,
				logger:      mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewUserUseCase(tt.args.userService, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	type fields struct {
		userService userService
		logger      logger
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
		want    *entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserService.EXPECT().
					Get(ctx, user.ID).
					Return(user, nil)
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
				mockUserService.EXPECT().
					Get(ctx, user.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
			i := &UserUseCase{
				userService: tt.fields.userService,
				logger:      tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	create := mock_entities.NewUserCreate(t)
	type fields struct {
		userService userService
		logger      logger
	}
	type args struct {
		ctx    context.Context
		create *entities.UserCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserService.EXPECT().Create(ctx, create).Return(user, nil)
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
				mockUserService.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
			i := &UserUseCase{
				userService: tt.fields.userService,
				logger:      tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	update := mock_entities.NewUserUpdate(t)
	type fields struct {
		userService userService
		logger      logger
	}
	type args struct {
		ctx    context.Context
		update *entities.UserUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserService.EXPECT().Update(ctx, update).Return(user, nil)
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
				mockUserService.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
			i := &UserUseCase{
				userService: tt.fields.userService,
				logger:      tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	type fields struct {
		userService userService
		logger      logger
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
				mockUserService.EXPECT().
					Delete(ctx, user.ID).
					Return(nil)
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
				mockUserService.EXPECT().
					Delete(ctx, user.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
			i := &UserUseCase{
				userService: tt.fields.userService,
				logger:      tt.fields.logger,
			}
			err := i.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}

func TestUserUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserService := NewMockuserService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mock_entities.NewUserFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listUsers := make([]*entities.User, 0, count)
	for i := uint64(0); i < count; i++ {
		listUsers = append(listUsers, mock_entities.NewUser(t))
	}
	type fields struct {
		userService userService
		logger      logger
	}
	type args struct {
		ctx    context.Context
		filter *entities.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*entities.User
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserService.EXPECT().
					List(ctx, filter).
					Return(listUsers, count, nil)
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
				mockUserService.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				userService: mockUserService,
				logger:      mockLogger,
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
			i := &UserUseCase{
				userService: tt.fields.userService,
				logger:      tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}
