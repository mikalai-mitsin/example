package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type CommentInterceptor struct {
	commentUseCase usecases.CommentUseCase
	authUseCase    usecases.AuthUseCase
	logger         log.Logger
}

func NewCommentInterceptor(
	commentUseCase usecases.CommentUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.CommentInterceptor {
	return &CommentInterceptor{
		commentUseCase: commentUseCase,
		authUseCase:    authUseCase,
		logger:         logger,
	}
}

func (i *CommentInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.Comment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentDetail,
	); err != nil {
		return nil, err
	}
	comment, err := i.commentUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentDetail,
		comment,
	)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (i *CommentInterceptor) List(
	ctx context.Context,
	filter *models.CommentFilter,
	requestUser *models.User,
) ([]*models.Comment, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	comments, count, err := i.commentUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return comments, count, nil
}

func (i *CommentInterceptor) Create(
	ctx context.Context,
	create *models.CommentCreate,
	requestUser *models.User,
) (*models.Comment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentCreate,
		create,
	); err != nil {
		return nil, err
	}
	comment, err := i.commentUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return comment, nil
}

func (i *CommentInterceptor) Update(
	ctx context.Context,
	update *models.CommentUpdate,
	requestUser *models.User,
) (*models.Comment, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentUpdate,
	); err != nil {
		return nil, err
	}
	comment, err := i.commentUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentUpdate,
		comment,
	); err != nil {
		return nil, err
	}
	updatedComment, err := i.commentUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedComment, nil
}

func (i *CommentInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentDelete,
	); err != nil {
		return err
	}
	comment, err := i.commentUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDCommentDelete,
		comment,
	)
	if err != nil {
		return err
	}
	if err := i.commentUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
