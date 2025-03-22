package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/comment/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type CommentServiceServer struct {
	examplepb.UnimplementedCommentServiceServer
	commentUseCase commentUseCase
	logger         logger
}

func NewCommentServiceServer(commentUseCase commentUseCase, logger logger) *CommentServiceServer {
	return &CommentServiceServer{commentUseCase: commentUseCase, logger: logger}
}

func (s *CommentServiceServer) Create(
	ctx context.Context,
	input *examplepb.CommentCreate,
) (*examplepb.Comment, error) {
	item, err := s.commentUseCase.Create(ctx, encodeCommentCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeComment(item), nil
}

func (s *CommentServiceServer) Get(
	ctx context.Context,
	input *examplepb.CommentGet,
) (*examplepb.Comment, error) {
	item, err := s.commentUseCase.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeComment(item), nil
}

func (s *CommentServiceServer) List(
	ctx context.Context,
	filter *examplepb.CommentFilter,
) (*examplepb.ListComment, error) {
	items, count, err := s.commentUseCase.List(ctx, encodeCommentFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListComment(items, count), nil
}

func (s *CommentServiceServer) Update(
	ctx context.Context,
	input *examplepb.CommentUpdate,
) (*examplepb.Comment, error) {
	item, err := s.commentUseCase.Update(ctx, encodeCommentUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeComment(item), nil
}

func (s *CommentServiceServer) Delete(
	ctx context.Context,
	input *examplepb.CommentDelete,
) (*emptypb.Empty, error) {
	if err := s.commentUseCase.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodeCommentCreate(input *examplepb.CommentCreate) entities.CommentCreate {
	create := entities.CommentCreate{
		Text:     input.GetText(),
		AuthorId: uuid.UUID(input.GetAuthorId()),
		PostId:   uuid.UUID(input.GetPostId()),
	}
	return create
}
func encodeCommentFilter(input *examplepb.CommentFilter) entities.CommentFilter {
	filter := entities.CommentFilter{
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
func encodeCommentUpdate(input *examplepb.CommentUpdate) entities.CommentUpdate {
	update := entities.CommentUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetText() != nil {
		update.Text = pointer.Pointer(string(input.GetText().GetValue()))
	}
	if input.GetAuthorId() != nil {
		update.AuthorId = pointer.Pointer(uuid.UUID(input.GetAuthorId().GetValue()))
	}
	if input.GetPostId() != nil {
		update.PostId = pointer.Pointer(uuid.UUID(input.GetPostId().GetValue()))
	}
	return update
}
func decodeComment(item entities.Comment) *examplepb.Comment {
	response := &examplepb.Comment{
		Id:        string(item.ID),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		Text:      item.Text,
		AuthorId:  string(item.AuthorId),
		PostId:    string(item.PostId),
	}
	return response
}
func decodeListComment(items []entities.Comment, count uint64) *examplepb.ListComment {
	response := &examplepb.ListComment{
		Items: make([]*examplepb.Comment, 0, len(items)),
		Count: count,
	}
	for _, item := range items {
		response.Items = append(response.Items, decodeComment(item))
	}
	return response
}
func decodeCommentUpdate(update entities.CommentUpdate) *examplepb.CommentUpdate {
	result := &examplepb.CommentUpdate{
		Id:       string(string(update.ID)),
		Text:     wrapperspb.String(*update.Text),
		AuthorId: wrapperspb.String(string(*update.AuthorId)),
		PostId:   wrapperspb.String(string(*update.PostId)),
	}
	return result
}
