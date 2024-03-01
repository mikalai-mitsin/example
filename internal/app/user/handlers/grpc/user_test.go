package grpc

import (
	"context"
	"errors"

	mock_grpc "github.com/018bf/example/internal/app/user/handlers/grpc/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"

	"reflect"
	"testing"

	"github.com/018bf/example/internal/app/user/models"
	mock_models "github.com/018bf/example/internal/app/user/models/mock"
	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewUserServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		userInterceptor UserInterceptor
		logger          log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.UserServiceServer
	}{
		{
			name: "ok",
			args: args{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
			want: &UserServiceServer{
				userInterceptor: userInterceptor,
				logger:          logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserServiceServer(tt.args.userInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewUserServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	// create := mock_models.NewUserCreate(t)
	user := mock_models.NewUser(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                UserInterceptor
		logger                         log.Logger
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
				userInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(user, nil)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.UserCreate{},
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.UserCreate{},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Create(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Create() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                UserInterceptor
		logger                         log.Logger
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
				userInterceptor.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
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
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserDelete{
					Id: id.String(),
				},
			},
			want: nil,
			wantErr: grpc.DecodeError(&errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params: map[string]string{
					"details": "i error",
				},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Delete(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Delete() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                UserInterceptor
		logger                         log.Logger
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
				userInterceptor.EXPECT().Get(ctx, user.ID).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
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
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Get(ctx, user.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.UserGet{
					Id: string(user.ID),
				},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Get(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewUserFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListUser{
		Items: make([]*examplepb.User, 0, int(count)),
		Count: count,
	}
	listUsers := make([]*models.User, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewUser(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listUsers = append(listUsers, a)
		response.Items = append(response.Items, decodeUser(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                UserInterceptor
		logger                         log.Logger
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
				userInterceptor.EXPECT().List(ctx, filter).Return(listUsers, count, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
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
			name: "interceptor error",
			setup: func() {
				userInterceptor.
					EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
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
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.List(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("List() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userInterceptor := mock_grpc.NewMockUserInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	user := mock_models.NewUser(t)
	update := mock_models.NewUserUpdate(t)
	type fields struct {
		UnimplementedUserServiceServer examplepb.UnimplementedUserServiceServer
		userInterceptor                UserInterceptor
		logger                         log.Logger
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
				userInterceptor.EXPECT().Update(ctx, gomock.Any()).Return(user, nil).Times(1)
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeUserUpdate(update),
			},
			want:    decodeUser(user),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				userInterceptor.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedUserServiceServer: examplepb.UnimplementedUserServiceServer{},
				userInterceptor:                userInterceptor,
				logger:                         logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeUserUpdate(update),
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := UserServiceServer{
				UnimplementedUserServiceServer: tt.fields.UnimplementedUserServiceServer,
				userInterceptor:                tt.fields.userInterceptor,
				logger:                         tt.fields.logger,
			}
			got, err := s.Update(tt.args.ctx, tt.args.input)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Update() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeUser(t *testing.T) {
	user := mock_models.NewUser(t)
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
		user *models.User
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
			if got := decodeUser(tt.args.user); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeUser() = %v, want %v", got, tt.want)
			}
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
		want *models.UserFilter
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
			want: &models.UserFilter{
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
			if got := encodeUserFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
