package services

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	mockEntities "github.com/mikalai-mitsin/example/internal/app/auth/entities/mock"
	userEntities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	mockUserEntities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"go.uber.org/mock/gomock"
)

func TestAuthService_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mockUserEntities.NewUser(t)
	type fields struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	type args struct {
		ctx    context.Context
		access entities.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *userEntities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().
					GetSubject(ctx, entities.Token("mytoken")).
					Return(string(user.ID), nil).
					Times(1)
				userRepository.EXPECT().Get(ctx, user.ID).Return(user, nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "bad user",
			setup: func() {
				authRepository.EXPECT().
					GetSubject(ctx, entities.Token("mytoken")).
					Return("", errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
		{
			name: "user not found",
			setup: func() {
				authRepository.EXPECT().
					GetSubject(ctx, entities.Token("mytoken")).
					Return(string(user.ID), nil).
					Times(1)
				userRepository.EXPECT().
					Get(ctx, user.ID).
					Return(nil, errs.NewEntityNotFoundError()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				access: "mytoken",
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthService{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Auth(tt.args.ctx, tt.args.access)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("Auth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Auth() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mockUserEntities.NewUser(t)
	login := mockEntities.NewLogin(t)
	user.Email = login.Email
	pair := mockEntities.NewTokenPair(t)
	user.SetPassword(login.Password)
	type fields struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	type args struct {
		ctx   context.Context
		login *entities.Login
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				userRepository.EXPECT().GetByEmail(ctx, user.Email).Return(user, nil).Times(1)
				authRepository.EXPECT().Create(ctx, user).Return(pair, nil).Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "user not found",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(nil, errs.NewEntityNotFoundError()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
		{
			name: "bad password",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(user, nil).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx: ctx,
				login: &entities.Login{
					Email:    login.Email,
					Password: "mojParol'",
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidParameter("email or password"),
		},
		{
			name: "bad password",
			setup: func() {
				userRepository.EXPECT().
					GetByEmail(ctx, user.Email).Return(user, nil).
					Times(1)
				authRepository.EXPECT().Create(ctx, user).
					Return(nil, errs.NewUnexpectedBehaviorError("system errpr")).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("system errpr"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthService{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.CreateToken(tt.args.ctx, tt.args.login)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	pair := mockEntities.NewTokenPair(t)
	type fields struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	type args struct {
		ctx     context.Context
		refresh entities.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().
					RefreshToken(ctx, entities.Token("my_r_token")).
					Return(pair, nil).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:     ctx,
				refresh: "my_r_token",
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authRepository.EXPECT().
					RefreshToken(ctx, entities.Token("my_r_token")).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:     ctx,
				refresh: "my_r_token",
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthService{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.RefreshToken(tt.args.ctx, tt.args.refresh)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("RefreshToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RefreshToken() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	type fields struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	type args struct {
		ctx    context.Context
		access entities.Token
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
		setup   func()
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().
					Validate(ctx, entities.Token("my_token")).
					Return(nil).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				access: "my_token",
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authRepository.EXPECT().
					Validate(ctx, entities.Token("my_token")).
					Return(errs.NewUnexpectedBehaviorError("error 345")).
					Times(1)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:    ctx,
				access: "my_token",
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "error 345"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := AuthService{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			if err := u.ValidateToken(tt.args.ctx, tt.args.access); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewAuthService(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	tests := []struct {
		name string
		args args
		want *AuthService
	}{
		{
			name: "ok",
			args: args{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			want: &AuthService{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthService(tt.args.authRepository, tt.args.userRepository, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthService_CreateTokenByUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	userRepository := NewMockUserRepository(ctrl)
	authRepository := NewMockAuthRepository(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	user := mockUserEntities.NewUser(t)
	tokenPair := mockEntities.NewTokenPair(t)
	type fields struct {
		authRepository AuthRepository
		userRepository UserRepository
		logger         logger
	}
	type args struct {
		ctx  context.Context
		user *userEntities.User
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *entities.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authRepository.EXPECT().Create(ctx, user).Return(tokenPair, nil)
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:  ctx,
				user: user,
			},
			want:    tokenPair,
			wantErr: nil,
		},
		{
			name: "error",
			setup: func() {
				authRepository.EXPECT().
					Create(ctx, user).
					Return(nil, errs.NewUnexpectedBehaviorError("asd"))
			},
			fields: fields{
				authRepository: authRepository,
				userRepository: userRepository,
				logger:         mockLogger,
			},
			args: args{
				ctx:  ctx,
				user: user,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("asd"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := AuthService{
				authRepository: tt.fields.authRepository,
				userRepository: tt.fields.userRepository,
				logger:         tt.fields.logger,
			}
			tt.setup()
			got, err := u.CreateTokenByUser(tt.args.ctx, tt.args.user)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CreateTokenByUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateTokenByUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}
