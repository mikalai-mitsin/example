package handlers

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func encodePostCreate(input *examplepb.PostCreate) entities.PostCreate {
	create := entities.PostCreate{Body: input.GetBody()}
	return create
}
func encodePostFilter(input *examplepb.PostFilter) entities.PostFilter {
	filter := entities.PostFilter{
		PageSize:   nil,
		PageNumber: nil,
		IsDeleted:  nil,
		OrderBy:    []entities.PostOrdering{},
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
	}
	if input.GetIsDeleted() != nil {
		filter.IsDeleted = pointer.Of(input.GetIsDeleted().GetValue())
	}
	for _, orderBy := range input.GetOrderBy() {
		filter.OrderBy = append(filter.OrderBy, entities.PostOrdering(orderBy))
	}
	return filter
}
func encodePostUpdate(input *examplepb.PostUpdate) entities.PostUpdate {
	update := entities.PostUpdate{ID: uuid.MustParse(input.GetId())}
	if input.GetBody() != nil {
		update.Body = pointer.Of(string(input.GetBody().GetValue()))
	}
	return update
}
func decodePost(post entities.Post) *examplepb.Post {
	response := &examplepb.Post{
		Id:        post.ID.String(),
		CreatedAt: timestamppb.New(post.CreatedAt),
		UpdatedAt: timestamppb.New(post.UpdatedAt),
		DeletedAt: nil,
		Body:      post.Body,
	}
	if post.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*post.DeletedAt)
	}
	if post.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*post.DeletedAt)
	}
	if post.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*post.DeletedAt)
	}
	if post.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*post.DeletedAt)
	}
	return response
}
func decodeListPost(items []entities.Post, count uint64) *examplepb.ListPost {
	response := &examplepb.ListPost{Items: make([]*examplepb.Post, 0, len(items)), Count: count}
	for _, post := range items {
		response.Items = append(response.Items, decodePost(post))
	}
	return response
}
func decodePostUpdate(update entities.PostUpdate) *examplepb.PostUpdate {
	result := &examplepb.PostUpdate{
		Id:   string(update.ID.String()),
		Body: wrapperspb.String(*update.Body),
	}
	return result
}
