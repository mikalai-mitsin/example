package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/app/arch/models"
	mock_models "github.com/018bf/example/internal/app/arch/models/mock"
	mock_usecases "github.com/018bf/example/internal/app/arch/usecases/mock"
	"github.com/018bf/example/internal/pkg/clock"
	mock_clock "github.com/018bf/example/internal/pkg/clock/mock"
	"github.com/018bf/example/internal/pkg/errs"
	"github.com/018bf/example/internal/pkg/log"
	mock_log "github.com/018bf/example/internal/pkg/log/mock"
	"github.com/018bf/example/internal/pkg/uuid"
	mock_uuid "github.com/018bf/example/internal/pkg/uuid/mock"
	"github.com/jaswdr/faker"
	"go.uber.org/mock/gomock"
)

func TestNewArchUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	mockLogger := mock_log.NewMockLogger(ctrl)
	mockUUID := mock_uuid.NewMockGenerator(ctrl)
	type args struct {
		archRepository ArchRepository
		clock          clock.Clock
		logger         log.Logger
		uuid           uuid.Generator
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  *ArchUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			want: &ArchUseCase{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewArchUseCase(tt.args.archRepository, tt.args.clock, tt.args.logger, tt.args.uuid); !reflect.DeepEqual(
				got,
				tt.want,
			) {
				t.Errorf("NewArchUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	type fields struct {
		archRepository ArchRepository
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
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archRepository.EXPECT().Get(ctx, arch.ID).Return(arch, nil)
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "Arch not found",
			setup: func() {
				archRepository.EXPECT().Get(ctx, arch.ID).Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArchUseCase{
				archRepository: tt.fields.archRepository,
				logger:         tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var listArches []*models.Arch
	count := faker.New().UInt64Between(2, 20)
	for i := uint64(0); i < count; i++ {
		listArches = append(listArches, mock_models.NewArch(t))
	}
	filter := mock_models.NewArchFilter(t)
	type fields struct {
		archRepository ArchRepository
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.ArchFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Arch
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				archRepository.EXPECT().List(ctx, filter).Return(listArches, nil)
				archRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    listArches,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				archRepository.EXPECT().
					List(ctx, filter).
					Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
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
				archRepository.EXPECT().List(ctx, filter).Return(listArches, nil)
				archRepository.EXPECT().
					Count(ctx, filter).
					Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
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
			u := &ArchUseCase{
				archRepository: tt.fields.archRepository,
				logger:         tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("ArchUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestArchUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	mockClock := mock_clock.NewMockClock(ctrl)
	mockLogger := mock_log.NewMockLogger(ctrl)
	mockUUID := mock_uuid.NewMockGenerator(ctrl)
	ctx := context.Background()
	create := mock_models.NewArchCreate(t)
	now := time.Now().UTC()
	type fields struct {
		archRepository ArchRepository
		clock          clock.Clock
		logger         log.Logger
		uuid           uuid.Generator
	}
	type args struct {
		ctx    context.Context
		create *models.ArchCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test"))
				archRepository.EXPECT().
					Create(
						ctx,
						&models.Arch{
							ID:          uuid.UUID("test"),
							Name:        create.Name,
							Title:       create.Title,
							Subtitle:    create.Subtitle,
							Tags:        create.Tags,
							Versions:    create.Versions,
							OldVersions: create.OldVersions,
							Release:     create.Release,
							Tested:      create.Tested,
							Mark:        create.Mark,
							Submarine:   create.Submarine,
							Numb:        create.Numb,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(nil)
			},
			fields: fields{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         mockLogger,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Arch{
				ID:          uuid.UUID("test"),
				Name:        create.Name,
				Title:       create.Title,
				Subtitle:    create.Subtitle,
				Tags:        create.Tags,
				Versions:    create.Versions,
				OldVersions: create.OldVersions,
				Release:     create.Release,
				Tested:      create.Tested,
				Mark:        create.Mark,
				Submarine:   create.Submarine,
				Numb:        create.Numb,
				UpdatedAt:   now,
				CreatedAt:   now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				mockUUID.EXPECT().NewUUID().Return(uuid.UUID("test 2"))
				archRepository.EXPECT().
					Create(
						ctx,
						&models.Arch{
							ID:          uuid.UUID("test 2"),
							Name:        create.Name,
							Title:       create.Title,
							Subtitle:    create.Subtitle,
							Tags:        create.Tags,
							Versions:    create.Versions,
							OldVersions: create.OldVersions,
							Release:     create.Release,
							Tested:      create.Tested,
							Mark:        create.Mark,
							Submarine:   create.Submarine,
							Numb:        create.Numb,
							UpdatedAt:   now,
							CreatedAt:   now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				archRepository: archRepository,
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
				archRepository: archRepository,
				logger:         mockLogger,
				clock:          mockClock,
				uuid:           mockUUID,
			},
			args: args{
				ctx:    ctx,
				create: &models.ArchCreate{},
			},
			want: nil,
			wantErr: errs.NewInvalidFormError().WithParams(
				errs.Param{Key: "name", Value: "cannot be blank"},
				errs.Param{Key: "title", Value: "cannot be blank"},
				errs.Param{Key: "subtitle", Value: "cannot be blank"},
				errs.Param{Key: "tags", Value: "cannot be blank"},
				errs.Param{Key: "versions", Value: "cannot be blank"},
				errs.Param{Key: "old_versions", Value: "cannot be blank"},
				errs.Param{Key: "release", Value: "cannot be blank"},
				errs.Param{Key: "tested", Value: "cannot be blank"},
				errs.Param{Key: "mark", Value: "cannot be blank"},
				errs.Param{Key: "submarine", Value: "cannot be blank"},
				errs.Param{Key: "numb", Value: "cannot be blank"},
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArchUseCase{
				archRepository: tt.fields.archRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
				uuid:           tt.fields.uuid,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	mockClock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewArchUpdate(t)
	now := arch.UpdatedAt
	type fields struct {
		archRepository ArchRepository
		clock          clock.Clock
		logger         log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.ArchUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Arch
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				archRepository.EXPECT().
					Get(ctx, update.ID).Return(arch, nil)
				archRepository.EXPECT().
					Update(ctx, arch).Return(nil)
			},
			fields: fields{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    arch,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				mockClock.EXPECT().Now().Return(now)
				archRepository.EXPECT().
					Get(ctx, update.ID).
					Return(arch, nil)
				archRepository.EXPECT().
					Update(ctx, arch).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Arch not found",
			setup: func() {
				archRepository.EXPECT().
					Get(ctx, update.ID).
					Return(nil, errs.NewEntityNotFoundError())
			},
			fields: fields{
				archRepository: archRepository,
				clock:          mockClock,
				logger:         logger,
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
				archRepository: archRepository,
				clock:          mockClock,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				update: &models.ArchUpdate{
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
			u := &ArchUseCase{
				archRepository: tt.fields.archRepository,
				clock:          tt.fields.clock,
				logger:         tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ArchUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestArchUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	archRepository := mock_usecases.NewMockArchRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	arch := mock_models.NewArch(t)
	type fields struct {
		archRepository ArchRepository
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
				archRepository.EXPECT().
					Delete(ctx, arch.ID).
					Return(nil)
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			wantErr: nil,
		},
		{
			name: "Arch not found",
			setup: func() {
				archRepository.EXPECT().
					Delete(ctx, arch.ID).
					Return(errs.NewEntityNotFoundError())
			},
			fields: fields{
				archRepository: archRepository,
				logger:         logger,
			},
			args: args{
				ctx: ctx,
				id:  arch.ID,
			},
			wantErr: errs.NewEntityNotFoundError(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &ArchUseCase{
				archRepository: tt.fields.archRepository,
				logger:         tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("ArchUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
