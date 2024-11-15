package services

import (
	"context"
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestNewUserService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	type args struct {
		userRepository userRepository
		clock          clock
		logger         logger
		uuid           UUIDGenerator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *UserService
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				userRepository: mockUserRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			want: &UserService{
				userRepository: mockUserRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			got := NewUserService(
				tt.args.userRepository,
				tt.args.clock,
				tt.args.logger,
				tt.args.uuid,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	type fields struct {
		userRepository userRepository
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
		want    *entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserRepository.EXPECT().Get(ctx, user.ID).Return(user, nil)
			},
			fields: fields{
				userRepository: mockUserRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				mockUserRepository.EXPECT().
					Get(ctx, user.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				userRepository: mockUserRepository,
				logger:         mockLogger,
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
			u := &UserService{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	var listUsers []*entities.User
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listUsers = append(listUsers, mock_entities.NewUser(t))
	}
	filter := mock_entities.NewUserFilter(t)
	type fields struct {
		userRepository userRepository
		logger         logger
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
				mockUserRepository.EXPECT().List(ctx, filter).Return(listUsers, nil)
				mockUserRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				userRepository: mockUserRepository,
				logger:         mockLogger,
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
				mockUserRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: mockUserRepository,
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
				mockUserRepository.EXPECT().List(ctx, filter).Return(listUsers, nil)
				mockUserRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: mockUserRepository,
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
			u := &UserService{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want1, got1)
		})
	}
}

func TestUserService_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockClock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	mockUUID := NewMockUUIDGenerator(ctrl)
	ctx := context.Background()
	create := mock_entities.NewUserCreate(t)
	now := time.Now().UTC()
	type fields struct {
		userRepository userRepository
		clock          clock
		logger         logger
		uuid           UUIDGenerator
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
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test"))
				mockUserRepository.EXPECT().
					Create(
						ctx,
						&entities.User{
							ID:        uuid.UUID("test"),
							FirstName: create.FirstName,
							LastName:  create.LastName,
							Password:  create.Password,
							Email:     create.Email,
							GroupID:   create.GroupID,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				userRepository: mockUserRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &entities.User{
				ID:        uuid.UUID("test"),
				FirstName: create.FirstName,
				LastName:  create.LastName,
				Password:  create.Password,
				Email:     create.Email,
				GroupID:   create.GroupID,
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
				mockUserRepository.EXPECT().
					Create(
						ctx,
						&entities.User{
							ID:        uuid.UUID("test 2"),
							FirstName: create.FirstName,
							LastName:  create.LastName,
							Password:  create.Password,
							Email:     create.Email,
							GroupID:   create.GroupID,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: mockUserRepository,
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
				userRepository: mockUserRepository,
				logger:         mockLogger,
				clock:          mockClock,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: &entities.UserCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "first_name", Value: "cannot be blank"},
				errs.Param{Key: "last_name", Value: "cannot be blank"},
				errs.Param{Key: "password", Value: "cannot be blank"},
				errs.Param{Key: "email", Value: "cannot be blank"},
				errs.Param{Key: "group_id", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserService{
				userRepository: tt.fields.userRepository,
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

func TestUserService_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	mockClock := NewMockclock(ctrl)
	update := mock_entities.NewUserUpdate(t)
	now := user.UpdatedAt
	type fields struct {
		userRepository userRepository
		clock          clock
		logger         logger
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
				mockClock.EXPECT().Now().Return(now)
				mockUserRepository.EXPECT().
					Get(ctx, update.ID).Return(user, nil)
				mockUserRepository.EXPECT().
					Update(ctx, user).Return(nil)
			},
			fields: fields{
				userRepository: mockUserRepository,
				clock:          mockClock,
				logger:         mockLogger,
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
				mockClock.EXPECT().Now().Return(now)
				mockUserRepository.EXPECT().
					Get(ctx, update.ID).
					Return(user, nil)
				mockUserRepository.EXPECT().
					Update(ctx, user).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				userRepository: mockUserRepository,
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
			name: "User not found",
			setup: func() {
				mockUserRepository.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				userRepository: mockUserRepository,
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
				userRepository: mockUserRepository,
				clock:          mockClock,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				update: &entities.UserUpdate{
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
			u := &UserService{
				userRepository: tt.fields.userRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserService_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepository := NewMockuserRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	type fields struct {
		userRepository userRepository
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
				mockUserRepository.EXPECT().
					Delete(ctx, user.ID).
					Return(nil)
			},
			fields: fields{
				userRepository: mockUserRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: nil,
		},
		{
			name: "User not found",
			setup: func() {
				mockUserRepository.EXPECT().
					Delete(ctx, user.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				userRepository: mockUserRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				id:  user.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &UserService{
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			err := u.Delete(tt.args.ctx, tt.args.id)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
