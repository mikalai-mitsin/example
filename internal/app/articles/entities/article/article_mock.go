package entities

import (
	"testing"
	"time"

	"github.com/jaswdr/faker"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

func NewMockArticle(t *testing.T) Article {
	t.Helper()
	return Article{
		ID:          uuid.NewUUID(),
		CreatedAt:   faker.New().Time().Time(time.Now()),
		UpdatedAt:   faker.New().Time().Time(time.Now()),
		Title:       faker.New().Lorem().Sentence(15),
		Subtitle:    faker.New().Lorem().Sentence(15),
		Body:        faker.New().Lorem().Sentence(15),
		IsPublished: faker.New().Bool(),
	}
}
func NewMockArticleFilter(t *testing.T) ArticleFilter {
	t.Helper()
	return ArticleFilter{
		PageSize:   pointer.Of(faker.New().UInt64()),
		PageNumber: pointer.Of(faker.New().UInt64()),
		Search:     pointer.Of(faker.New().Lorem().Sentence(15)),
		OrderBy: []ArticleOrdering{
			ArticleOrderingSubtitleDESC,
			ArticleOrderingBodyASC,
			ArticleOrderingIsPublishedDESC,
			ArticleOrderingIdASC,
			ArticleOrderingUpdatedAtDESC,
			ArticleOrderingTitleASC,
			ArticleOrderingTitleDESC,
			ArticleOrderingBodyDESC,
			ArticleOrderingIsPublishedASC,
			ArticleOrderingIdDESC,
			ArticleOrderingCreatedAtASC,
			ArticleOrderingCreatedAtDESC,
			ArticleOrderingUpdatedAtASC,
			ArticleOrderingSubtitleASC,
		},
	}
}
func NewMockArticleCreate(t *testing.T) ArticleCreate {
	t.Helper()
	return ArticleCreate{
		Title:       faker.New().Lorem().Sentence(15),
		Subtitle:    faker.New().Lorem().Sentence(15),
		Body:        faker.New().Lorem().Sentence(15),
		IsPublished: faker.New().Bool(),
	}
}
func NewMockArticleUpdate(t *testing.T) ArticleUpdate {
	t.Helper()
	return ArticleUpdate{
		ID:          uuid.NewUUID(),
		Title:       pointer.Of(faker.New().Lorem().Sentence(15)),
		Subtitle:    pointer.Of(faker.New().Lorem().Sentence(15)),
		Body:        pointer.Of(faker.New().Lorem().Sentence(15)),
		IsPublished: pointer.Of(faker.New().Bool()),
	}
}
