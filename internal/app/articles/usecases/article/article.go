package usecases

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArticleUseCase struct {
	articleService       articleService
	articleEventProducer articleEventProducer
	logger               logger
}

func NewArticleUseCase(
	articleService articleService,
	articleEventProducer articleEventProducer,
	logger logger,
) *ArticleUseCase {
	return &ArticleUseCase{
		articleService:       articleService,
		articleEventProducer: articleEventProducer,
		logger:               logger,
	}
}

func (i *ArticleUseCase) Create(
	ctx context.Context,
	create entities.ArticleCreate,
) (entities.Article, error) {
	article, err := i.articleService.Create(ctx, create)
	if err != nil {
		return entities.Article{}, err
	}
	if err := i.articleEventProducer.Created(ctx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (i *ArticleUseCase) Get(ctx context.Context, id uuid.UUID) (entities.Article, error) {
	article, err := i.articleService.Get(ctx, id)
	if err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (i *ArticleUseCase) List(
	ctx context.Context,
	filter entities.ArticleFilter,
) ([]entities.Article, uint64, error) {
	articles, count, err := i.articleService.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return articles, count, nil
}

func (i *ArticleUseCase) Update(
	ctx context.Context,
	update entities.ArticleUpdate,
) (entities.Article, error) {
	article, err := i.articleService.Update(ctx, update)
	if err != nil {
		return entities.Article{}, err
	}
	if err := i.articleEventProducer.Updated(ctx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (i *ArticleUseCase) Delete(ctx context.Context, id uuid.UUID) error {
	if err := i.articleService.Delete(ctx, id); err != nil {
		return err
	}
	if err := i.articleEventProducer.Deleted(ctx, id); err != nil {
		return err
	}
	return nil
}
