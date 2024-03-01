package interceptors

import (
	"context"
	"errors"
	"reflect"
	"testing"

	mock_interceptors "github.com/018bf/example/internal/app/session/interceptors/mock"
	"github.com/018bf/example/internal/app/session/models"
	mock_models "github.com/018bf/example/internal/app/session/models/mock"
	userModels "github.com/018bf/example/internal/app/user/models"
	userMockModels "github.com/018bf/example/internal/app/user/models/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"

	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"github.com/jaswdr/faker"

	"github.com/018bf/example/internal/pkg/uuid"
)

func TestNewSessionInterceptor(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		authUseCase    AuthUseCase
		sessionUseCase SessionUseCase
		logger         log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *SessionInterceptor
	}{
		{
			name:  "ok",
			setup: func() {},
			args: args{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			want: &SessionInterceptor{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewSessionInterceptor(tt.args.sessionUseCase, tt.args.logger, tt.args.authUseCase); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewSessionInterceptor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		authUseCase    AuthUseCase
		sessionUseCase SessionUseCase
		logger         log.Logger
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
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDetail).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDetail, session).
					Return(nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  uuid.UUID(session.ID),
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "object permission error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDetail).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDetail, session).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDetail).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "Session not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDetail).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	create := mock_models.NewSessionCreate(t)
	type fields struct {
		sessionUseCase SessionUseCase
		authUseCase    AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.SessionCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionCreate, create).
					Return(nil)
				sessionUseCase.EXPECT().Create(ctx, create).Return(session, nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionCreate, create).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionCreate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "create error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionCreate).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionCreate, create).
					Return(nil)
				sessionUseCase.EXPECT().
					Create(ctx, create).
					Return(nil, errs.NewUnexpectedBehaviorError("c u"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
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
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	update := mock_models.NewSessionUpdate(t)
	type fields struct {
		sessionUseCase SessionUseCase
		authUseCase    AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.SessionUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Session
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate, session).
					Return(nil)
				sessionUseCase.EXPECT().Update(ctx, update).Return(session, nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    session,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate, session).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "update error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, update.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate, session).
					Return(nil)
				sessionUseCase.EXPECT().
					Update(ctx, update).
					Return(nil, errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionUpdate).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, err := i.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSessionInterceptor_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	session := mock_models.NewSession(t)
	type fields struct {
		sessionUseCase SessionUseCase
		authUseCase    AuthUseCase
		logger         log.Logger
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
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDelete, session).
					Return(nil)
				sessionUseCase.EXPECT().
					Delete(ctx, session.ID).
					Return(nil)
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: nil,
		},
		{
			name: "Session not found",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, errs.NewEntityNotFound())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDelete, session).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "delete error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete).
					Return(nil)
				sessionUseCase.EXPECT().
					Get(ctx, session.ID).
					Return(session, nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionDelete, session).
					Return(nil)
				sessionUseCase.EXPECT().
					Delete(ctx, session.ID).
					Return(errs.NewUnexpectedBehaviorError("d 2"))
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewUnexpectedBehaviorError("d 2"),
		},
		{
			name: "permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionDelete).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				authUseCase:    authUseCase,
				sessionUseCase: sessionUseCase,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  session.ID,
			},
			wantErr: errs.NewPermissionDeniedError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			if err := i.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSessionInterceptor_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	authUseCase := mock_interceptors.NewMockAuthUseCase(ctrl)
	requestUser := userMockModels.NewUser(t)
	sessionUseCase := mock_interceptors.NewMockSessionUseCase(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	filter := mock_models.NewSessionFilter(t)
	count := faker.New().UInt64Between(2, 20)
	listSessions := make([]*models.Session, 0, count)
	for i := uint64(0); i < count; i++ {
		listSessions = append(listSessions, mock_models.NewSession(t))
	}
	type fields struct {
		sessionUseCase SessionUseCase
		authUseCase    AuthUseCase
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.SessionFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Session
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionList, filter).
					Return(nil)
				sessionUseCase.EXPECT().
					List(ctx, filter).
					Return(listSessions, count, nil)
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listSessions,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "object permission denied",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionList, filter).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "permission error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionList).
					Return(errs.NewPermissionDeniedError())
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    nil,
			want1:   0,
			wantErr: errs.NewPermissionDeniedError(),
		},
		{
			name: "list error",
			setup: func() {
				authUseCase.EXPECT().GetUser(ctx).Return(requestUser, nil)
				authUseCase.EXPECT().
					HasPermission(ctx, requestUser, userModels.PermissionIDSessionList).
					Return(nil)
				authUseCase.EXPECT().
					HasObjectPermission(ctx, requestUser, userModels.PermissionIDSessionList, filter).
					Return(nil)
				sessionUseCase.EXPECT().
					List(ctx, filter).
					Return(nil, uint64(0), errs.NewUnexpectedBehaviorError("l e"))
			},
			fields: fields{
				sessionUseCase: sessionUseCase,
				authUseCase:    authUseCase,
				logger:         logger,
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
			i := &SessionInterceptor{
				sessionUseCase: tt.fields.sessionUseCase,
				authUseCase:    tt.fields.authUseCase,
				logger:         tt.fields.logger,
			}
			got, got1, err := i.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("SessionInterceptor.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SessionInterceptor.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("SessionInterceptor.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
