package interceptors

import (
	"context"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	"github.com/018bf/example/internal/domain/usecases"

	"github.com/018bf/example/pkg/log"
)

type PostInterceptor struct {
	postUseCase usecases.PostUseCase
	authUseCase usecases.AuthUseCase
	logger      log.Logger
}

func NewPostInterceptor(
	postUseCase usecases.PostUseCase,
	authUseCase usecases.AuthUseCase,
	logger log.Logger,
) interceptors.PostInterceptor {
	return &PostInterceptor{
		postUseCase: postUseCase,
		authUseCase: authUseCase,
		logger:      logger,
	}
}

func (i *PostInterceptor) Get(
	ctx context.Context,
	id string,
	requestUser *models.User,
) (*models.Post, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDPostDetail,
	); err != nil {
		return nil, err
	}
	post, err := i.postUseCase.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDPostDetail,
		post,
	)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (i *PostInterceptor) List(
	ctx context.Context,
	filter *models.PostFilter,
	requestUser *models.User,
) ([]*models.Post, uint64, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDPostList,
	); err != nil {
		return nil, 0, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDPostList,
		filter,
	); err != nil {
		return nil, 0, err
	}
	posts, count, err := i.postUseCase.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return posts, count, nil
}

func (i *PostInterceptor) Create(
	ctx context.Context,
	create *models.PostCreate,
	requestUser *models.User,
) (*models.Post, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDPostCreate,
	); err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDPostCreate,
		create,
	); err != nil {
		return nil, err
	}
	post, err := i.postUseCase.Create(ctx, create)
	if err != nil {
		return nil, err
	}
	return post, nil
}

func (i *PostInterceptor) Update(
	ctx context.Context,
	update *models.PostUpdate,
	requestUser *models.User,
) (*models.Post, error) {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDPostUpdate,
	); err != nil {
		return nil, err
	}
	post, err := i.postUseCase.Get(ctx, update.ID)
	if err != nil {
		return nil, err
	}
	if err := i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDPostUpdate,
		post,
	); err != nil {
		return nil, err
	}
	updatedPost, err := i.postUseCase.Update(ctx, update)
	if err != nil {
		return nil, err
	}
	return updatedPost, nil
}

func (i *PostInterceptor) Delete(
	ctx context.Context,
	id string,
	requestUser *models.User,
) error {
	if err := i.authUseCase.HasPermission(
		ctx,
		requestUser,
		models.PermissionIDPostDelete,
	); err != nil {
		return err
	}
	post, err := i.postUseCase.Get(ctx, id)
	if err != nil {
		return err
	}
	err = i.authUseCase.HasObjectPermission(
		ctx,
		requestUser,
		models.PermissionIDPostDelete,
		post,
	)
	if err != nil {
		return err
	}
	if err := i.postUseCase.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
