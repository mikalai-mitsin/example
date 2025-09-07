package handlers

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PostServiceServer struct {
	examplepb.UnimplementedPostServiceServer
	postUseCase postUseCase
	logger      logger
}

func NewPostServiceServer(postUseCase postUseCase, logger logger) *PostServiceServer {
	return &PostServiceServer{postUseCase: postUseCase, logger: logger}
}

func (s *PostServiceServer) Create(
	ctx context.Context,
	input *examplepb.PostCreate,
) (*examplepb.Post, error) {
	item, err := s.postUseCase.Create(ctx, encodePostCreate(input))
	if err != nil {
		return nil, err
	}
	return decodePost(item), nil
}

func (s *PostServiceServer) Get(
	ctx context.Context,
	input *examplepb.PostGet,
) (*examplepb.Post, error) {
	item, err := s.postUseCase.Get(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodePost(item), nil
}

func (s *PostServiceServer) List(
	ctx context.Context,
	filter *examplepb.PostFilter,
) (*examplepb.ListPost, error) {
	items, count, err := s.postUseCase.List(ctx, encodePostFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListPost(items, count), nil
}

func (s *PostServiceServer) Update(
	ctx context.Context,
	input *examplepb.PostUpdate,
) (*examplepb.Post, error) {
	item, err := s.postUseCase.Update(ctx, encodePostUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodePost(item), nil
}

func (s *PostServiceServer) Delete(
	ctx context.Context,
	input *examplepb.PostDelete,
) (*emptypb.Empty, error) {
	if err := s.postUseCase.Delete(ctx, uuid.MustParse(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func (s *PostServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.PostService_ServiceDesc, s)
	return nil
}
func encodePostCreate(input *examplepb.PostCreate) entities.PostCreate {
	create := entities.PostCreate{Body: input.GetBody()}
	return create
}
func encodePostFilter(input *examplepb.PostFilter) entities.PostFilter {
	filter := entities.PostFilter{
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    []entities.PostOrdering{},
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Of(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Of(input.GetPageNumber().GetValue())
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
func decodePost(item entities.Post) *examplepb.Post {
	response := &examplepb.Post{
		Id:        item.ID.String(),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		Body:      item.Body,
	}
	return response
}
func decodeListPost(items []entities.Post, count uint64) *examplepb.ListPost {
	response := &examplepb.ListPost{Items: make([]*examplepb.Post, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodePost(item))
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
