package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArticleUseCase struct {
	articleService      articleService
	articleEventService articleEventService
	dtxManager          dtxManager
	logger              logger
}

func NewArticleUseCase(
	articleService articleService,
	articleEventService articleEventService,
	dtxManager dtxManager,
	logger logger,
) *ArticleUseCase {
	return &ArticleUseCase{
		articleService:      articleService,
		articleEventService: articleEventService,
		dtxManager:          dtxManager,
		logger:              logger,
	}
}

func (u *ArticleUseCase) Create(
	ctx context.Context,
	create entities.ArticleCreate,
) (entities.Article, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	article, err := u.articleService.Create(ctx, tx, create)
	if err != nil {
		return entities.Article{}, err
	}
	if err := u.articleEventService.Send(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (u *ArticleUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Article, error) {
	article, err := u.articleService.Get(ctx, id)
	if err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (u *ArticleUseCase) List(
	ctx context.Context,
	filter entities.ArticleFilter,
) ([]entities.Article, uint64, error) {
	articles, count, err := u.articleService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return articles, count, nil
}

func (u *ArticleUseCase) Update(
	ctx context.Context,
	update entities.ArticleUpdate,
) (entities.Article, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	article, err := u.articleService.Update(ctx, tx, update)
	if err != nil {
		return entities.Article{}, err
	}
	if err := u.articleEventService.Send(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (u *ArticleUseCase) Delete(
	ctx context.Context,
	del entities.ArticleDelete,
) (entities.Article, error) {
	logger := u.logger.WithContext(ctx)
	tx := u.dtxManager.NewTx()
	defer func(tx dtx.TX) {
		if err := tx.Rollback(); err != nil {
			logger.Error("cant rollback transaction", log.Error(err))
		}
	}(tx)
	article, err := u.articleService.Delete(ctx, tx, del)
	if err != nil {
		return entities.Article{}, err
	}
	if err := u.articleEventService.Send(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	if err := tx.Commit(); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
