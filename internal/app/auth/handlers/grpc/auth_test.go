package handlers

import (
	"context"
	"errors"

	mock_grpc "github.com/mikalai-mitsin/example/internal/app/auth/handlers/grpc/mock"
	"github.com/mikalai-mitsin/example/internal/app/auth/models"
	mock_models "github.com/mikalai-mitsin/example/internal/app/auth/models/mock"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"

	"reflect"
	"testing"

	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"go.uber.org/mock/gomock"
)

func TestAuthServiceServer_CreateToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authInterceptor := mock_grpc.NewMockAuthInterceptor(ctrl)
	ctx := context.Background()
	login := mock_models.NewLogin(t)
	pair := mock_models.NewTokenPair(t)
	type fields struct {
		UnimplementedAuthServiceServer examplepb.UnimplementedAuthServiceServer
		authInterceptor                AuthInterceptor
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
				authInterceptor.EXPECT().CreateToken(ctx, login).Return(pair, nil).Times(1)

			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authInterceptor:                authInterceptor,
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
			name: "interceptor error",
			setup: func() {
				authInterceptor.EXPECT().
					CreateToken(ctx, login).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authInterceptor:                authInterceptor,
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
				authInterceptor:                tt.fields.authInterceptor,
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
	authInterceptor := mock_grpc.NewMockAuthInterceptor(ctrl)
	ctx := context.Background()
	token := mock_models.NewToken(t)
	pair := mock_models.NewTokenPair(t)
	type fields struct {
		UnimplementedAuthServiceServer examplepb.UnimplementedAuthServiceServer
		authInterceptor                AuthInterceptor
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
				authInterceptor.EXPECT().RefreshToken(ctx, token).Return(pair, nil).Times(1)

			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authInterceptor:                authInterceptor,
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
			name: "interceptor error",
			setup: func() {
				authInterceptor.EXPECT().
					RefreshToken(ctx, token).
					Return(nil, errs.NewBadTokenError()).
					Times(1)
			},
			fields: fields{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authInterceptor:                authInterceptor,
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
				authInterceptor:                tt.fields.authInterceptor,
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
	authInterceptor := mock_grpc.NewMockAuthInterceptor(ctrl)
	type args struct {
		authInterceptor AuthInterceptor
	}
	tests := []struct {
		name string
		args args
		want examplepb.AuthServiceServer
	}{
		{
			name: "ok",
			args: args{
				authInterceptor: authInterceptor,
			},
			want: &AuthServiceServer{
				UnimplementedAuthServiceServer: examplepb.UnimplementedAuthServiceServer{},
				authInterceptor:                authInterceptor,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthServiceServer(tt.args.authInterceptor); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewAuthServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_decodeTokenPair(t *testing.T) {
	type args struct {
		pair *models.TokenPair
	}
	tests := []struct {
		name string
		args args
		want *examplepb.TokenPair
	}{
		{
			name: "ok",
			args: args{
				pair: &models.TokenPair{
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
