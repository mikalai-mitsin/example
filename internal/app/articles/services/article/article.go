package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

type ArticleService struct {
	articleRepository articleRepository
	clock             clock
	logger            logger
	uuid              uuidGenerator
}

func NewArticleService(
	articleRepository articleRepository,
	clock clock,
	logger logger,
	uuid uuidGenerator,
) *ArticleService {
	return &ArticleService{
		articleRepository: articleRepository,
		clock:             clock,
		logger:            logger,
		uuid:              uuid,
	}
}

func (u *ArticleService) Create(
	ctx context.Context,
	create entities.ArticleCreate,
) (entities.Article, error) {
	if err := create.Validate(); err != nil {
		return entities.Article{}, err
	}
	now := u.clock.Now().UTC()
	article := entities.Article{
		ID:          u.uuid.NewUUID(),
		UpdatedAt:   now,
		CreatedAt:   now,
		Title:       create.Title,
		Subtitle:    create.Subtitle,
		Body:        create.Body,
		IsPublished: create.IsPublished,
	}
	if err := u.articleRepository.Create(ctx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (u *ArticleService) Get(ctx context.Context, id uuid.UUID) (entities.Article, error) {
	article, err := u.articleRepository.Get(ctx, id)
	if err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (u *ArticleService) List(
	ctx context.Context,
	filter entities.ArticleFilter,
) ([]entities.Article, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	article, err := u.articleRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := u.articleRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return article, count, nil
}

func (u *ArticleService) Update(
	ctx context.Context,
	update entities.ArticleUpdate,
) (entities.Article, error) {
	if err := update.Validate(); err != nil {
		return entities.Article{}, err
	}
	article, err := u.articleRepository.Get(ctx, update.ID)
	if err != nil {
		return entities.Article{}, err
	}
	{
		if update.Title != nil {
			article.Title = *update.Title
		}
		if update.Subtitle != nil {
			article.Subtitle = *update.Subtitle
		}
		if update.Body != nil {
			article.Body = *update.Body
		}
		if update.IsPublished != nil {
			article.IsPublished = *update.IsPublished
		}
	}
	article.UpdatedAt = u.clock.Now().UTC()
	if err := u.articleRepository.Update(ctx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (u *ArticleService) Delete(ctx context.Context, id uuid.UUID) error {
	if err := u.articleRepository.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}
