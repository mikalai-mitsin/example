package handlers

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type LikeServiceServer struct {
	examplepb.UnimplementedLikeServiceServer
	likeUseCase likeUseCase
	logger      logger
}

func NewLikeServiceServer(likeUseCase likeUseCase, logger logger) *LikeServiceServer {
	return &LikeServiceServer{likeUseCase: likeUseCase, logger: logger}
}

func (s *LikeServiceServer) Create(
	ctx context.Context,
	input *examplepb.LikeCreate,
) (*examplepb.Like, error) {
	item, err := s.likeUseCase.Create(ctx, encodeLikeCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeLike(item), nil
}

func (s *LikeServiceServer) Get(
	ctx context.Context,
	input *examplepb.LikeGet,
) (*examplepb.Like, error) {
	item, err := s.likeUseCase.Get(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeLike(item), nil
}

func (s *LikeServiceServer) List(
	ctx context.Context,
	filter *examplepb.LikeFilter,
) (*examplepb.ListLike, error) {
	items, count, err := s.likeUseCase.List(ctx, encodeLikeFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListLike(items, count), nil
}

func (s *LikeServiceServer) Update(
	ctx context.Context,
	input *examplepb.LikeUpdate,
) (*examplepb.Like, error) {
	item, err := s.likeUseCase.Update(ctx, encodeLikeUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeLike(item), nil
}

func (s *LikeServiceServer) Delete(
	ctx context.Context,
	input *examplepb.LikeDelete,
) (*emptypb.Empty, error) {
	if err := s.likeUseCase.Delete(ctx, uuid.MustParse(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodeLikeCreate(input *examplepb.LikeCreate) entities.LikeCreate {
	create := entities.LikeCreate{
		PostID: uuid.MustParse(input.GetPostId()),
		Value:  input.GetValue(),
		UserId: uuid.MustParse(input.GetUserId()),
	}
	return create
}
func encodeLikeFilter(input *examplepb.LikeFilter) entities.LikeFilter {
	filter := entities.LikeFilter{
		IDs:        nil,
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    input.GetOrderBy(),
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, uuid.MustParse(id))
	}
	return filter
}
func encodeLikeUpdate(input *examplepb.LikeUpdate) entities.LikeUpdate {
	update := entities.LikeUpdate{ID: uuid.MustParse(input.GetId())}
	if input.GetPostId() != nil {
		update.PostID = pointer.Of(uuid.MustParse(input.GetPostId().GetValue()))
	}
	if input.GetValue() != nil {
		update.Value = pointer.Of(string(input.GetValue().GetValue()))
	}
	if input.GetUserId() != nil {
		update.UserId = pointer.Of(uuid.MustParse(input.GetUserId().GetValue()))
	}
	return update
}
func decodeLike(item entities.Like) *examplepb.Like {
	response := &examplepb.Like{
		Id:        item.ID.String(),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		PostId:    item.PostID.String(),
		Value:     item.Value,
		UserId:    item.UserId.String(),
	}
	return response
}
func decodeListLike(items []entities.Like, count uint64) *examplepb.ListLike {
	response := &examplepb.ListLike{Items: make([]*examplepb.Like, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeLike(item))
	}
	return response
}
func decodeLikeUpdate(update entities.LikeUpdate) *examplepb.LikeUpdate {
	result := &examplepb.LikeUpdate{
		Id:     string(update.ID.String()),
		PostId: wrapperspb.String(update.PostID.String()),
		Value:  wrapperspb.String(*update.Value),
		UserId: wrapperspb.String(update.UserId.String()),
	}
	return result
}
