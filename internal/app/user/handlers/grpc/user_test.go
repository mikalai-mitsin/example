package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"testing"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/app/user/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewUserServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		userUseCase userUseCase
		logger      logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.UserServiceServer
	}{
		{
			name: "ok",
			args: args{
				userUseCase: mockUserUseCase,
				logger:      mockLogger,
			},
			want: &UserServiceServer{
				userUseCase: mockUserUseCase,
				logger:      mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewUserServiceServer(tt.args.userUseCase, tt.args.logger)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	// create := mock_entities.NewUserCreate(t)
	user := mock_entities.NewUser(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userUseCase                    userUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(user, nil)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.UserCreate{},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockUserUseCase.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("usecase error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.UserCreate{},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("usecase error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userUseCase:                    tt.fields.userUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userUseCase                    userUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserDelete
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *emptypb.Empty
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserUseCase.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockUserUseCase.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserDelete{
					Id: id.String(),
				},
			},
			want: nil,
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "i error"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userUseCase:                    tt.fields.userUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userUseCase                    userUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserUseCase.EXPECT().Get(ctx, user.ID).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserGet{
					Id: string(user.ID),
				},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockUserUseCase.EXPECT().Get(ctx, user.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserGet{
					Id: string(user.ID),
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userUseCase:                    tt.fields.userUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	filter := mock_entities.NewUserFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListUser{
		Items: make([]*examplepb.User, 0, int(count)),
		Count: count,
	}
	listUsers := make([]*entities.User, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_entities.NewUser(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listUsers = append(listUsers, a)
		response.Items = append(response.Items, decodeUser(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userUseCase                    userUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListUser
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserUseCase.EXPECT().
					List(ctx, gomock.Any()).
					Return(listUsers, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    response,
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockUserUseCase.
					EXPECT().
					List(ctx, gomock.Any()).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserFilter{
					PageNumber: wrapperspb.UInt64(*filter.PageNumber),
					PageSize:   wrapperspb.UInt64(*filter.PageSize),
					Search:     wrapperspb.String(*filter.Search),
					OrderBy:    filter.OrderBy,
					Ids:        stringIDs,
				},
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userUseCase:                    tt.fields.userUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestUserServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserUseCase := NewMockuserUseCase(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mock_entities.NewUser(t)
	update := mock_entities.NewUserUpdate(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userUseCase                    userUseCase
		logger                         logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.UserUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockUserUseCase.EXPECT().Update(ctx, gomock.Any()).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeUserUpdate(update),
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				mockUserUseCase.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userUseCase:                    mockUserUseCase,
				logger:                         mockLogger,
			},
			args: args{
				ctx:   ctx,
				input: decodeUserUpdate(update),
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("i error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userUseCase:                    tt.fields.userUseCase,
				logger:                         tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			assert.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_decodeUser(t *testing.T) {
	user := mock_entities.NewUser(t)
	result := &examplepb.User{
		Id:        string(user.ID),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
		CreatedAt: timestamppb.New(user.CreatedAt),
		FirstName: string(user.FirstName),
		LastName:  string(user.LastName),
		Password:  string(user.Password),
		Email:     string(user.Email),
		GroupId:   string(user.GroupID),
	}
	type args struct {
		user *entities.User
	}
	tests := []struct {
		name string
		args args
		want *examplepb.User
	}{
		{
			name: "ok",
			args: args{
				user: user,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := decodeUser(tt.args.user)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_encodeUserFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.UserFilter
	}
	tests := []struct {
		name string
		args args
		want *entities.UserFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.UserFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &entities.UserFilter{
				PageSize:   pointer.Pointer(uint64(5)),
				PageNumber: pointer.Pointer(uint64(2)),
				OrderBy:    []string{"created_at", "id"},
				Search:     pointer.Pointer("my name is"),
				IDs:        []uuid.UUID{id},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := encodeUserFilter(tt.args.input)
			assert.Equal(t, tt.want, got)
		})
	}
}
