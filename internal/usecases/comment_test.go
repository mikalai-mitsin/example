package usecases

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/018bf/example/internal/domain/errs"
	"github.com/018bf/example/internal/domain/models"
	mock_models "github.com/018bf/example/internal/domain/models/mock"
	"github.com/018bf/example/internal/domain/repositories"
	mock_repositories "github.com/018bf/example/internal/domain/repositories/mock"
	"github.com/018bf/example/internal/domain/usecases"
	"github.com/018bf/example/pkg/clock"
	mock_clock "github.com/018bf/example/pkg/clock/mock"
	"github.com/018bf/example/pkg/log"
	mock_log "github.com/018bf/example/pkg/log/mock"
	"github.com/golang/mock/gomock"
	"syreclabs.com/go/faker"
)

func TestNewCommentUseCase(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	type args struct {
		commentRepository repositories.CommentRepository
		clock             clock.Clock
		logger            log.Logger
	}
	tests := []struct {
		name  string
		setup func()
		args  args
		want  usecases.CommentUseCase
	}{
		{
			name: "ok",
			setup: func() {
			},
			args: args{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			want: &CommentUseCase{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			if got := NewCommentUseCase(tt.args.commentRepository, tt.args.clock, tt.args.logger); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCommentUseCase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentUseCase_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	type fields struct {
		commentRepository repositories.CommentRepository
		logger            log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				commentRepository.EXPECT().Get(ctx, comment.ID).Return(comment, nil)
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				commentRepository.EXPECT().Get(ctx, comment.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentUseCase{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			got, err := u.Get(tt.args.ctx, tt.args.id)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentUseCase.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentUseCase.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentUseCase_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	var comments []*models.Comment
	count := uint64(faker.Number().NumberInt(2))
	for i := uint64(0); i < count; i++ {
		comments = append(comments, mock_models.NewComment(t))
	}
	filter := mock_models.NewCommentFilter(t)
	type fields struct {
		commentRepository repositories.CommentRepository
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		filter *models.CommentFilter
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    []*models.Comment
		want1   uint64
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				commentRepository.EXPECT().List(ctx, filter).Return(comments, nil)
				commentRepository.EXPECT().Count(ctx, filter).Return(count, nil)
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				filter: filter,
			},
			want:    comments,
			want1:   count,
			wantErr: nil,
		},
		{
			name: "list error",
			setup: func() {
				commentRepository.EXPECT().List(ctx, filter).Return(nil, errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
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
				commentRepository.EXPECT().List(ctx, filter).Return(comments, nil)
				commentRepository.EXPECT().Count(ctx, filter).Return(uint64(0), errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
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
			u := &CommentUseCase{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			got, got1, err := u.List(tt.args.ctx, tt.args.filter)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentUseCase.List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentUseCase.List() = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("CommentUseCase.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestCommentUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	clockMock := mock_clock.NewMockClock(ctrl)
	ctx := context.Background()
	create := mock_models.NewCommentCreate(t)
	now := time.Now().UTC()
	type fields struct {
		commentRepository repositories.CommentRepository
		clock             clock.Clock
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		create *models.CommentCreate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				commentRepository.EXPECT().
					Create(
						ctx,
						&models.Comment{
							Body:      create.Body,
							PostId:    create.PostId,
							UserId:    create.UserId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(nil)
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want: &models.Comment{
				ID:        "",
				Body:      create.Body,
				PostId:    create.PostId,
				UserId:    create.UserId,
				UpdatedAt: now,
				CreatedAt: now,
			},
			wantErr: nil,
		},
		{
			name: "unexpected behavior",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				commentRepository.EXPECT().
					Create(
						ctx,
						&models.Comment{
							ID:        "",
							Body:      create.Body,
							PostId:    create.PostId,
							UserId:    create.UserId,
							UpdatedAt: now,
							CreatedAt: now,
						},
					).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				create: create,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		// TODO: Add validation rules or delete this case
		//{
		//	name: "invalid",
		//	setup: func() {
		//	},
		//	fields: fields{
		//		commentRepository: commentRepository,
		//		logger:           logger,
		//	},
		//	args: args{
		//		ctx: ctx,
		//		create: &models.CommentCreate{},
		//	},
		//	want: nil,
		//	wantErr: errs.NewInvalidFormError().WithParam("set", "it"),
		// },
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentUseCase{
				commentRepository: tt.fields.commentRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Create(tt.args.ctx, tt.args.create)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentUseCase.Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentUseCase.Create() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentUseCase_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	clockMock := mock_clock.NewMockClock(ctrl)
	update := mock_models.NewCommentUpdate(t)
	now := comment.UpdatedAt
	type fields struct {
		commentRepository repositories.CommentRepository
		clock             clock.Clock
		logger            log.Logger
	}
	type args struct {
		ctx    context.Context
		update *models.CommentUpdate
	}
	tests := []struct {
		name    string
		setup   func()
		fields  fields
		args    args
		want    *models.Comment
		wantErr error
	}{
		{
			name: "ok",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				commentRepository.EXPECT().
					Get(ctx, update.ID).Return(comment, nil)
				commentRepository.EXPECT().
					Update(ctx, comment).Return(nil)
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    comment,
			wantErr: nil,
		},
		{
			name: "update error",
			setup: func() {
				clockMock.EXPECT().Now().Return(now)
				commentRepository.EXPECT().
					Get(ctx, update.ID).
					Return(comment, nil)
				commentRepository.EXPECT().
					Update(ctx, comment).
					Return(errs.NewUnexpectedBehaviorError("test error"))
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewUnexpectedBehaviorError("test error"),
		},
		{
			name: "Comment not found",
			setup: func() {
				commentRepository.EXPECT().Get(ctx, update.ID).Return(nil, errs.NewEntityNotFound())
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx:    ctx,
				update: update,
			},
			want:    nil,
			wantErr: errs.NewEntityNotFound(),
		},
		{
			name: "invalid",
			setup: func() {
			},
			fields: fields{
				commentRepository: commentRepository,
				clock:             clockMock,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				update: &models.CommentUpdate{
					ID: faker.Number().Number(1),
				},
			},
			want:    nil,
			wantErr: errs.NewInvalidFormError().WithParam("id", "must be a valid UUID"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentUseCase{
				commentRepository: tt.fields.commentRepository,
				clock:             tt.fields.clock,
				logger:            tt.fields.logger,
			}
			got, err := u.Update(tt.args.ctx, tt.args.update)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentUseCase.Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CommentUseCase.Update() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommentUseCase_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	commentRepository := mock_repositories.NewMockCommentRepository(ctrl)
	logger := mock_log.NewMockLogger(ctrl)
	ctx := context.Background()
	comment := mock_models.NewComment(t)
	type fields struct {
		commentRepository repositories.CommentRepository
		logger            log.Logger
	}
	type args struct {
		ctx context.Context
		id  string
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
				commentRepository.EXPECT().
					Delete(ctx, comment.ID).
					Return(nil)
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: nil,
		},
		{
			name: "Comment not found",
			setup: func() {
				commentRepository.EXPECT().
					Delete(ctx, comment.ID).
					Return(errs.NewEntityNotFound())
			},
			fields: fields{
				commentRepository: commentRepository,
				logger:            logger,
			},
			args: args{
				ctx: ctx,
				id:  comment.ID,
			},
			wantErr: errs.NewEntityNotFound(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()
			u := &CommentUseCase{
				commentRepository: tt.fields.commentRepository,
				logger:            tt.fields.logger,
			}
			if err := u.Delete(tt.args.ctx, tt.args.id); !errors.Is(err, tt.wantErr) {
				t.Errorf("CommentUseCase.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
