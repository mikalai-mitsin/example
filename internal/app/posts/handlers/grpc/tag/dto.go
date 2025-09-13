package handlers

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func encodeTagCreate(input *examplepb.TagCreate) entities.TagCreate {
	create := entities.TagCreate{PostId: uuid.MustParse(input.GetPostId()), Value: input.GetValue()}
	return create
}
func encodeTagFilter(input *examplepb.TagFilter) entities.TagFilter {
	filter := entities.TagFilter{
		PageSize:   nil,
		PageNumber: nil,
		IsDeleted:  nil,
		OrderBy:    []entities.TagOrdering{},
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
		filter.OrderBy = append(filter.OrderBy, entities.TagOrdering(orderBy))
	}
	return filter
}
func encodeTagUpdate(input *examplepb.TagUpdate) entities.TagUpdate {
	update := entities.TagUpdate{ID: uuid.MustParse(input.GetId())}
	if input.GetPostId() != nil {
		update.PostId = pointer.Of(uuid.MustParse(input.GetPostId().GetValue()))
	}
	if input.GetValue() != nil {
		update.Value = pointer.Of(string(input.GetValue().GetValue()))
	}
	return update
}
func decodeTag(tag entities.Tag) *examplepb.Tag {
	response := &examplepb.Tag{
		Id:        tag.ID.String(),
		CreatedAt: timestamppb.New(tag.CreatedAt),
		UpdatedAt: timestamppb.New(tag.UpdatedAt),
		DeletedAt: nil,
		PostId:    tag.PostId.String(),
		Value:     tag.Value,
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	if tag.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*tag.DeletedAt)
	}
	return response
}
func decodeListTag(items []entities.Tag, count uint64) *examplepb.ListTag {
	response := &examplepb.ListTag{Items: make([]*examplepb.Tag, 0, len(items)), Count: count}
	for _, tag := range items {
		response.Items = append(response.Items, decodeTag(tag))
	}
	return response
}
func decodeTagUpdate(update entities.TagUpdate) *examplepb.TagUpdate {
	result := &examplepb.TagUpdate{
		Id:     string(update.ID.String()),
		PostId: wrapperspb.String(update.PostId.String()),
		Value:  wrapperspb.String(*update.Value),
	}
	return result
}
