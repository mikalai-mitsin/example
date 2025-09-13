package services

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
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

func (s *ArticleService) Create(
	ctx context.Context,
	tx dtx.TX,
	create entities.ArticleCreate,
) (entities.Article, error) {
	if err := create.Validate(); err != nil {
		return entities.Article{}, err
	}
	now := s.clock.Now().UTC()
	article := entities.Article{
		ID:          s.uuid.NewUUID(),
		UpdatedAt:   now,
		CreatedAt:   now,
		Title:       create.Title,
		Subtitle:    create.Subtitle,
		Body:        create.Body,
		IsPublished: create.IsPublished,
	}
	if err := s.articleRepository.Create(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
func (s *ArticleService) Get(ctx context.Context, id uuid.UUID) (entities.Article, error) {
	article, err := s.articleRepository.Get(ctx, id)
	if err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (s *ArticleService) List(
	ctx context.Context,
	filter entities.ArticleFilter,
) ([]entities.Article, uint64, error) {
	if err := filter.Validate(); err != nil {
		return nil, 0, err
	}
	article, err := s.articleRepository.List(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	count, err := s.articleRepository.Count(ctx, filter)
	if err != nil {
		return nil, 0, err
	}
	return article, count, nil
}

func (s *ArticleService) Update(
	ctx context.Context,
	tx dtx.TX,
	update entities.ArticleUpdate,
) (entities.Article, error) {
	if err := update.Validate(); err != nil {
		return entities.Article{}, err
	}
	article, err := s.articleRepository.Get(ctx, update.ID)
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
	article.UpdatedAt = s.clock.Now().UTC()
	if err := s.articleRepository.Update(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}

func (s *ArticleService) Delete(
	ctx context.Context,
	tx dtx.TX,
	id uuid.UUID,
) (entities.Article, error) {
	article, err := s.articleRepository.Get(ctx, id)
	if err != nil {
		return entities.Article{}, err
	}
	article.DeletedAt = pointer.Of(s.clock.Now().UTC())
	if err := s.articleRepository.Update(ctx, tx, article); err != nil {
		return entities.Article{}, err
	}
	return article, nil
}
