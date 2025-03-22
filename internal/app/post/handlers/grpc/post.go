package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/post/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
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
	item, err := s.postUseCase.Get(ctx, uuid.UUID(input.GetId()))
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
	if err := s.postUseCase.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodePostCreate(input *examplepb.PostCreate) entities.PostCreate {
	create := entities.PostCreate{
		Title:       input.GetTitle(),
		Body:        input.GetBody(),
		IsPrivate:   input.GetIsPrivate(),
		Tags:        input.GetTags(),
		PublishedAt: input.GetPublishedAt().AsTime(),
		AuthorId:    uuid.UUID(input.GetAuthorId()),
	}
	return create
}
func encodePostFilter(input *examplepb.PostFilter) entities.PostFilter {
	filter := entities.PostFilter{
		IDs:        nil,
		PageSize:   nil,
		PageNumber: nil,
		OrderBy:    input.GetOrderBy(),
		Search:     nil,
	}
	if input.GetPageSize() != nil {
		filter.PageSize = pointer.Pointer(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = pointer.Pointer(input.GetPageNumber().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, uuid.UUID(id))
	}
	return filter
}
func encodePostUpdate(input *examplepb.PostUpdate) entities.PostUpdate {
	update := entities.PostUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetTitle() != nil {
		update.Title = pointer.Pointer(string(input.GetTitle().GetValue()))
	}
	if input.GetBody() != nil {
		update.Body = pointer.Pointer(string(input.GetBody().GetValue()))
	}
	if input.GetIsPrivate() != nil {
		update.IsPrivate = pointer.Pointer(bool(input.GetIsPrivate().GetValue()))
	}
	if input.GetTags() != nil {
		var params []string
		for _, item := range input.GetTags().GetValues() {
			params = append(params, string(item.GetStringValue()))
		}
		update.Tags = pointer.Pointer(params)
	}
	if input.GetPublishedAt() != nil {
		update.PublishedAt = pointer.Pointer(input.GetPublishedAt().AsTime())
	}
	if input.GetAuthorId() != nil {
		update.AuthorId = pointer.Pointer(uuid.UUID(input.GetAuthorId().GetValue()))
	}
	return update
}
func decodePost(item entities.Post) *examplepb.Post {
	response := &examplepb.Post{
		Id:          string(item.ID),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Title:       item.Title,
		Body:        item.Body,
		IsPrivate:   item.IsPrivate,
		Tags:        item.Tags,
		PublishedAt: timestamppb.New(item.PublishedAt),
		AuthorId:    string(item.AuthorId),
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
		Id:          string(string(update.ID)),
		Title:       wrapperspb.String(*update.Title),
		Body:        wrapperspb.String(*update.Body),
		IsPrivate:   wrapperspb.Bool(*update.IsPrivate),
		Tags:        nil,
		PublishedAt: timestamppb.New(*update.PublishedAt),
		AuthorId:    wrapperspb.String(string(*update.AuthorId)),
	}
	if update.Tags != nil {
		params, err := structpb.NewList(pointer.ToAnySlice(*update.Tags))
		if err != nil {
			return nil
		}
		result.Tags = params
	}
	return result
}
