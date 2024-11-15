package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/auth/entities/mock"
	user_entities "github.com/mikalai-mitsin/example/internal/app/user/entities"
	mock_user_entities "github.com/mikalai-mitsin/example/internal/app/user/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"go.uber.org/mock/gomock"
)

func TestAuthUseCase_Auth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authService := NewMockAuthService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	token := mock_entities.NewToken(t)
	user := mock_user_entities.NewUser(t)
	type fields struct {
		authService AuthService
		logger      logger
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
		want    *user_entities.User
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authService.EXPECT().Auth(ctx, token).Return(user, nil).Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
				authService.EXPECT().
					Auth(ctx, token).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
			i := AuthUseCase{
				authService: tt.fields.authService,
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

func TestAuthUseCase_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authService := NewMockAuthService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	login := mock_entities.NewLogin(t)
	pair := mock_entities.NewTokenPair(t)
	clockmock := NewMockclock(ctrl)
	type fields struct {
		authService AuthService
		logger      logger
		clock       clock
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
				authService.EXPECT().CreateToken(ctx, login).Return(pair, nil).Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
				authService.EXPECT().
					CreateToken(ctx, login).
					Return(nil, errs.NewInvalidParameter("email or password")).
					Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
			i := AuthUseCase{
				authService: tt.fields.authService,
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

func TestAuthUseCase_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authService := NewMockAuthService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	pair := mock_entities.NewTokenPair(t)
	refresh := mock_entities.NewToken(t)
	clockmock := NewMockclock(ctrl)
	type fields struct {
		authService AuthService
		logger      logger
		clock       clock
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
				authService.EXPECT().RefreshToken(ctx, refresh).Return(pair, nil).Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
				authService.EXPECT().
					RefreshToken(ctx, refresh).
					Return(nil, errs.NewBadTokenError()).Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
			i := AuthUseCase{
				authService: tt.fields.authService,
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

func TestNewAuthUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authService := NewMockAuthService(ctrl)
	clockmock := NewMockclock(ctrl)
	mockLogger := NewMocklogger(ctrl)
	type args struct {
		authService AuthService
		logger      logger
		clock       clock
	}
	tests := []struct {
		name string
		args args
		want *AuthUseCase
	}{
		{
			name: "ok",
			args: args{
				authService: authService,
				logger:      mockLogger,
				clock:       clockmock,
			},
			want: &AuthUseCase{
				authService: authService,
				clock:       clockmock,
				logger:      mockLogger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewAuthUseCase(tt.args.authService, tt.args.clock, tt.args.logger)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthUseCase_ValidateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authService := NewMockAuthService(ctrl)
	mockLogger := NewMocklogger(ctrl)
	ctx := context.Background()
	token := entities.Token("this_is_valid_token")
	type fields struct {
		authService AuthService
		logger      logger
	}
	type args struct {
		ctx   context.Context
		token entities.Token
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
				authService.EXPECT().ValidateToken(ctx, token).Return(nil).Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
				authService.EXPECT().
					ValidateToken(ctx, token).
					Return(errs.NewUnexpectedBehaviorError("35124345")).
					Times(1)
			},
			fields: fields{
				authService: authService,
				logger:      mockLogger,
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
			i := AuthUseCase{
				authService: tt.fields.authService,
				logger:      tt.fields.logger,
			}
			if err := i.ValidateToken(tt.args.ctx, tt.args.token); !errors.Is(err, tt.wantErr) {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
