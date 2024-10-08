package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/mikalai-mitsin/example/internal/app/auth/interceptors/mock"
	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	mock_models "github.com/mikalai-mitsin/example/internal/app/auth/models/mock"
	user_models "github.com/mikalai-mitsin/example/internal/app/user/models"
	mock_user_models "github.com/mikalai-mitsin/example/internal/app/user/models/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"go.uber.org/mock/gomock"
)

func TestAuthInterceptor_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	token := mock_models.NewToken(t)
	user := mock_user_models.NewUser(t)
	type fields struct {
		authUseCase AuthUseCase
		logger      Logger
	}
	type args struct {
		ctx    context.Context
		access models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *user_models.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().Auth(ctx, token).Return(user, nil).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				access: token,
			},
			want:    user,
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authUseCase.EXPECT().
					Auth(ctx, token).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:    ctx,
				access: token,
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := AuthInterceptor{
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			got, err := i.Auth(tt.args.ctx, tt.args.access)
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

func TestAuthInterceptor_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	login := mock_models.NewLogin(t)
	pair := mock_models.NewTokenPair(t)
	clockmock := mock_interceptors.NewMockClock(ctrl)
	type fields struct {
		authUseCase AuthUseCase
		logger      Logger
		clock       Clock
	}
	type args struct {
		ctx   context.Context
		login *models.Login
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().CreateToken(ctx, login).Return(pair, nil).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "create requestUser error",
			setup: func() {
				authUseCase.EXPECT().
					CreateToken(ctx, login).
					Return(nil, errs.NewInvalidParameter("email or password")).
					Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			args: args{
				ctx:   ctx,
				login: login,
			},
			want:    nil,
			wantErr: errs.NewInvalidParameter("email or password"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := AuthInterceptor{
				authUseCase: tt.fields.authUseCase,
				clock:       tt.fields.clock,
				logger:      tt.fields.logger,
			}
			got, err := i.CreateToken(tt.args.ctx, tt.args.login)
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

func TestAuthInterceptor_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	pair := mock_models.NewTokenPair(t)
	refresh := mock_models.NewToken(t)
	clockmock := mock_interceptors.NewMockClock(ctrl)
	type fields struct {
		authUseCase AuthUseCase
		logger      Logger
		clock       Clock
	}
	type args struct {
		ctx     context.Context
		refresh models.Token
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().RefreshToken(ctx, refresh).Return(pair, nil).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			args: args{
				ctx:     ctx,
				refresh: refresh,
			},
			want:    pair,
			wantErr: nil,
		},
		{
			name: "bad requestUser",
			setup: func() {
				authUseCase.EXPECT().
					RefreshToken(ctx, refresh).
					Return(nil, errs.NewBadTokenError()).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			args: args{
				ctx:     ctx,
				refresh: refresh,
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := AuthInterceptor{
				authUseCase: tt.fields.authUseCase,
				clock:       tt.fields.clock,
				logger:      tt.fields.logger,
			}
			got, err := i.RefreshToken(tt.args.ctx, tt.args.refresh)
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

func TestNewAuthInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	clockmock := mock_interceptors.NewMockClock(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	type args struct {
		authUseCase AuthUseCase
		logger      Logger
		clock       Clock
	}
	tests := []struct {
		name string
		args args
		want *AuthInterceptor
	}{
		{
			name: "ok",
			args: args{
				authUseCase: authUseCase,
				logger:      logger,
				clock:       clockmock,
			},
			want: &AuthInterceptor{
				authUseCase: authUseCase,
				clock:       clockmock,
				logger:      logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthInterceptor(tt.args.authUseCase, tt.args.clock, tt.args.logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthInterceptor_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	logger := mock_interceptors.NewMockLogger(ctrl)
	ctx := context.Background()
	token := models.Token("this_is_valid_token")
	type fields struct {
		authUseCase AuthUseCase
		logger      Logger
	}
	type args struct {
		ctx   context.Context
		token models.Token
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
				authUseCase.EXPECT().ValidateToken(ctx, token).Return(nil).Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:   ctx,
				token: token,
			},
			wantErr: nil,
		},
		{
			name: "repository error",
			setup: func() {
				authUseCase.EXPECT().
					ValidateToken(ctx, token).
					Return(errs.NewUnexpectedBehaviorError("35124345")).
					Times(1)
			},
			fields: fields{
				authUseCase: authUseCase,
				logger:      logger,
			},
			args: args{
				ctx:   ctx,
				token: token,
			},
			wantErr: &errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "35124345"}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := AuthInterceptor{
				authUseCase: tt.fields.authUseCase,
				logger:      tt.fields.logger,
			}
			if err := i.ValidateToken(tt.args.ctx, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
