package grpc

import (
	"context"
	"errors"

	mock_grpc "github.com/018bf/example/internal/app/session/handlers/grpc/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"

	"reflect"
	"testing"

	"github.com/018bf/example/internal/app/session/models"
	mock_models "github.com/018bf/example/internal/app/session/models/mock"
	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func TestNewSessionServiceServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		sessionInterceptor SessionInterceptor
		logger             log.Logger
	}
	tests := []struct {
		name string
		args args
		want examplepb.SessionServiceServer
	}{
		{
			name: "ok",
			args: args{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
			want: &SessionServiceServer{
				sessionInterceptor: sessionInterceptor,
				logger:             logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionServiceServer(tt.args.sessionInterceptor, tt.args.logger); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSessionServiceServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionServiceServer_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	// create := mock_models.NewSessionCreate(t)
	session := mock_models.NewSession(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(session, nil)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.SessionCreate{},
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.
					EXPECT().
					Create(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("interceptor error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: &examplepb.SessionCreate{},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("interceptor error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	id := uuid.NewUUID()
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionDelete
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
				sessionInterceptor.EXPECT().Delete(ctx, id).Return(nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionDelete{
					Id: id.String(),
				},
			},
			want:    &emptypb.Empty{},
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Delete(ctx, id).
					Return(errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionDelete{
					Id: id.String(),
				},
			},
			want: nil,
			wantErr: grpc.DecodeError(&errs.Error{
				Code:    13,
				Message: "Unexpected behavior.",
				Params:  errs.Params{{Key: "details", Value: "i error"}},
			}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionGet
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().Get(ctx, session.ID).Return(session, nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionGet{
					Id: string(session.ID),
				},
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Get(ctx, session.ID).
					Return(nil, errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionGet{
					Id: string(session.ID),
				},
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		tt.setup()
		t.Run(tt.name, func(t *testing.T) {
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewSessionFilter(t)
	var ids []uuid.UUID
	var stringIDs []string
	count := faker.New().UInt64Between(2, 20)
	response := &examplepb.ListSession{
		Items: make([]*examplepb.Session, 0, int(count)),
		Count: count,
	}
	listSessions := make([]*models.Session, 0, int(count))
	for i := 0; i < int(count); i++ {
		a := mock_models.NewSession(t)
		ids = append(ids, a.ID)
		stringIDs = append(stringIDs, string(a.ID))
		listSessions = append(listSessions, a)
		response.Items = append(response.Items, decodeSession(a))
	}
	filter.IDs = ids
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.ListSession
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().
					List(ctx, filter).
					Return(listSessions, count, nil).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionFilter{
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
				sessionInterceptor.
					EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("i error")).
					Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx: ctx,
				input: &examplepb.SessionFilter{
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
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func TestSessionServiceServer_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	sessionInterceptor := mock_grpc.NewMockSessionInterceptor(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	update := mock_models.NewSessionUpdate(t)
	type fields struct {
		UnimplementedSessionServiceServer examplepb.UnimplementedSessionServiceServer
		sessionInterceptor                SessionInterceptor
		logger                            log.Logger
	}
	type args struct {
		ctx   context.Context
		input *examplepb.SessionUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *examplepb.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				sessionInterceptor.EXPECT().Update(ctx, gomock.Any()).Return(session, nil).Times(1)
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeSessionUpdate(update),
			},
			want:    decodeSession(session),
			wantErr: nil,
		},
		{
			name: "interceptor error",
			setup: func() {
				sessionInterceptor.EXPECT().Update(ctx, gomock.Any()).
					Return(nil, errs.NewUnexpectedBehaviorError("i error"))
			},
			fields: fields{
				UnimplementedSessionServiceServer: examplepb.UnimplementedSessionServiceServer{},
				sessionInterceptor:                sessionInterceptor,
				logger:                            logger,
			},
			args: args{
				ctx:   ctx,
				input: decodeSessionUpdate(update),
			},
			want:    nil,
			wantErr: grpc.DecodeError(errs.NewUnexpectedBehaviorError("i error")),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			s := SessionServiceServer{
				UnimplementedSessionServiceServer: tt.fields.UnimplementedSessionServiceServer,
				sessionInterceptor:                tt.fields.sessionInterceptor,
				logger:                            tt.fields.logger,
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

func Test_decodeSession(t *testing.T) {
	session := mock_models.NewSession(t)
	result := &examplepb.Session{
		Id:          string(session.ID),
		UpdatedAt:   timestamppb.New(session.UpdatedAt),
		CreatedAt:   timestamppb.New(session.CreatedAt),
		Title:       string(session.Title),
		Description: string(session.Description),
	}
	type args struct {
		session *models.Session
	}
	tests := []struct {
		name string
		args args
		want *examplepb.Session
	}{
		{
			name: "ok",
			args: args{
				session: session,
			},
			want: result,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := decodeSession(tt.args.session); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("decodeSession() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_encodeSessionFilter(t *testing.T) {
	id := uuid.UUID(uuid.NewUUID())
	type args struct {
		input *examplepb.SessionFilter
	}
	tests := []struct {
		name string
		args args
		want *models.SessionFilter
	}{
		{
			name: "ok",
			args: args{
				input: &examplepb.SessionFilter{
					PageNumber: wrapperspb.UInt64(2),
					PageSize:   wrapperspb.UInt64(5),
					Search:     wrapperspb.String("my name is"),
					OrderBy:    []string{"created_at", "id"},
					Ids:        []string{string(id)},
				},
			},
			want: &models.SessionFilter{
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
			if got := encodeSessionFilter(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("encodeUserFilter() = %v, want %v", got, tt.want)
			}
		})
	}
}
