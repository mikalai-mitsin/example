package handlers

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type TagServiceServer struct {
	examplepb.UnimplementedTagServiceServer
	tagUseCase tagUseCase
	logger     logger
}

func NewTagServiceServer(tagUseCase tagUseCase, logger logger) *TagServiceServer {
	return &TagServiceServer{tagUseCase: tagUseCase, logger: logger}
}

func (s *TagServiceServer) Create(
	ctx context.Context,
	input *examplepb.TagCreate,
) (*examplepb.Tag, error) {
	item, err := s.tagUseCase.Create(ctx, encodeTagCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeTag(item), nil
}

func (s *TagServiceServer) Get(
	ctx context.Context,
	input *examplepb.TagGet,
) (*examplepb.Tag, error) {
	item, err := s.tagUseCase.Get(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeTag(item), nil
}

func (s *TagServiceServer) List(
	ctx context.Context,
	filter *examplepb.TagFilter,
) (*examplepb.ListTag, error) {
	items, count, err := s.tagUseCase.List(ctx, encodeTagFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListTag(items, count), nil
}

func (s *TagServiceServer) Update(
	ctx context.Context,
	input *examplepb.TagUpdate,
) (*examplepb.Tag, error) {
	item, err := s.tagUseCase.Update(ctx, encodeTagUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeTag(item), nil
}

func (s *TagServiceServer) Delete(
	ctx context.Context,
	input *examplepb.TagDelete,
) (*emptypb.Empty, error) {
	if err := s.tagUseCase.Delete(ctx, uuid.MustParse(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *TagServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.TagService_ServiceDesc, s)
	return nil
}
func encodeTagCreate(input *examplepb.TagCreate) entities.TagCreate {
	create := entities.TagCreate{PostId: uuid.MustParse(input.GetPostId()), Value: input.GetValue()}
	return create
}
func encodeTagFilter(input *examplepb.TagFilter) entities.TagFilter {
	filter := entities.TagFilter{
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    []entities.TagOrdering{},
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
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
func decodeTag(item entities.Tag) *examplepb.Tag {
	response := &examplepb.Tag{
		Id:        item.ID.String(),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		PostId:    item.PostId.String(),
		Value:     item.Value,
	}
	return response
}
func decodeListTag(items []entities.Tag, count uint64) *examplepb.ListTag {
	response := &examplepb.ListTag{Items: make([]*examplepb.Tag, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeTag(item))
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
