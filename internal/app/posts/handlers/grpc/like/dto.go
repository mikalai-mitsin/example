package handlers

import (
	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func encodeLikeCreate(input *examplepb.LikeCreate) entities.LikeCreate {
	create := entities.LikeCreate{
		PostId: uuid.MustParse(input.GetPostId()),
		Value:  input.GetValue(),
		UserId: uuid.MustParse(input.GetUserId()),
	}
	return create
}
func encodeLikeFilter(input *examplepb.LikeFilter) entities.LikeFilter {
	filter := entities.LikeFilter{
		PageSize:   nil,
		PageNumber: nil,
		IsDeleted:  nil,
		OrderBy:    []entities.LikeOrdering{},
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
		filter.OrderBy = append(filter.OrderBy, entities.LikeOrdering(orderBy))
	}
	return filter
}
func encodeLikeUpdate(input *examplepb.LikeUpdate) entities.LikeUpdate {
	update := entities.LikeUpdate{ID: uuid.MustParse(input.GetId())}
	if input.GetPostId() != nil {
		update.PostId = pointer.Of(uuid.MustParse(input.GetPostId().GetValue()))
	}
	if input.GetValue() != nil {
		update.Value = pointer.Of(string(input.GetValue().GetValue()))
	}
	if input.GetUserId() != nil {
		update.UserId = pointer.Of(uuid.MustParse(input.GetUserId().GetValue()))
	}
	return update
}
func encodeLikeDelete(input *examplepb.LikeDelete) entities.LikeDelete {
	del := entities.LikeDelete{ID: uuid.MustParse(input.GetId())}
	return del
}
func decodeLike(like entities.Like) *examplepb.Like {
	response := &examplepb.Like{
		Id:        like.ID.String(),
		CreatedAt: timestamppb.New(like.CreatedAt),
		UpdatedAt: timestamppb.New(like.UpdatedAt),
		DeletedAt: nil,
		PostId:    like.PostId.String(),
		Value:     like.Value,
		UserId:    like.UserId.String(),
	}
	if like.DeletedAt != nil {
		response.DeletedAt = timestamppb.New(*like.DeletedAt)
	}
	return response
}
func decodeListLike(items []entities.Like, count uint64) *examplepb.ListLike {
	response := &examplepb.ListLike{Items: make([]*examplepb.Like, 0, len(items)), Count: count}
	for _, like := range items {
		response.Items = append(response.Items, decodeLike(like))
	}
	return response
}
func decodeLikeUpdate(update entities.LikeUpdate) *examplepb.LikeUpdate {
	result := &examplepb.LikeUpdate{
		Id:     string(update.ID.String()),
		PostId: wrapperspb.String(update.PostId.String()),
		Value:  wrapperspb.String(*update.Value),
		UserId: wrapperspb.String(update.UserId.String()),
	}
	return result
}
