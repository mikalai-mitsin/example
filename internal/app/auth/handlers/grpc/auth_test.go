package handlers

import (
	"context"
	"errors"

	"github.com/mikalai-mitsin/example/internal/app/auth/entities"
	mock_entities "github.com/mikalai-mitsin/example/internal/app/auth/entities/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"reflect"
	"testing"

	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"go.uber.org/mock/gomock"
)

func TestAuthServiceServer_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := NewMockAuthUseCase(ctrl)
	ctx := context.Background()
	login := mock_entities.NewLogin(t)
	pair := mock_entities.NewTokenPair(t)
	type fields struct {
		UnimplementedAuthServiceServer examplepb.UnimplementedAuthServiceServer
		authUseCase                    AuthUseCase
	}
	type args struct {
		ctx   context.Context
		input *examplepb.CreateToken
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().CreateToken(ctx, login).Return(pair, nil).Times(1)

			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authUseCase:                    authUseCase,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CreateToken{
					Email:    login.Email,
					Password: login.Password,
				},
			},
			want:    decodeTokenPair(pair),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				authUseCase.EXPECT().
					CreateToken(ctx, login).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authUseCase:                    authUseCase,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.CreateToken{
					Email:    login.Email,
					Password: login.Password,
				},
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := AuthServiceServer{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				authUseCase:                    tt.fields.authUseCase,
			}
			got, err := s.CreateToken(tt.args.ctx, tt.args.input)
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

func TestAuthServiceServer_RefreshToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := NewMockAuthUseCase(ctrl)
	ctx := context.Background()
	token := mock_entities.NewToken(t)
	pair := mock_entities.NewTokenPair(t)
	type fields struct {
		UnimplementedAuthServiceServer examplepb.UnimplementedAuthServiceServer
		authUseCase                    AuthUseCase
	}
	type args struct {
		ctx   context.Context
		input *examplepb.RefreshToken
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.TokenPair
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().RefreshToken(ctx, token).Return(pair, nil).Times(1)

			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authUseCase:                    authUseCase,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.RefreshToken{
					Token: token.String(),
				},
			},
			want:    decodeTokenPair(pair),
			wantErr: nil,
		},
		{
			name: "usecase error",
			setup: func() {
				authUseCase.EXPECT().
					RefreshToken(ctx, token).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authUseCase:                    authUseCase,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.RefreshToken{
					Token: token.String(),
				},
			},
			want:    nil,
			wantErr: errs.NewBadTokenError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := AuthServiceServer{
				UnimplementedAuthServiceServer: tt.fields.UnimplementedAuthServiceServer,
				authUseCase:                    tt.fields.authUseCase,
			}
			got, err := s.RefreshToken(tt.args.ctx, tt.args.input)
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

func TestNewAuthServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := NewMockAuthUseCase(ctrl)
	type args struct {
		authUseCase AuthUseCase
	}
	tests := []struct {
		name string
		args args
		want examplepb.AuthServiceServer
	}{
		{
			name: "ok",
			args: args{
				authUseCase: authUseCase,
			},
			want: &AuthServiceServer{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authUseCase:                    authUseCase,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthServiceServer(tt.args.authUseCase); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeTokenPair(t *testing.T) {
	type args struct {
		pair *entities.TokenPair
	}
	tests := []struct {
		name string
		args args
		want *examplepb.TokenPair
	}{
		{
			name: "ok",
			args: args{
				pair: &entities.TokenPair{
					Access:  "dasasdasd",
					Refresh: "asdartge245",
				},
			},
			want: &examplepb.TokenPair{
				Access:  "dasasdasd",
				Refresh: "asdartge245",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeTokenPair(tt.args.pair); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeTokenPair() = %v, want %v", got, tt.want)
			}
		})
	}
}
