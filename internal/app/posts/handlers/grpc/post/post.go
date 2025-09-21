package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
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
) (*examplepb.Post, error) {
	post, err := s.postUseCase.Delete(ctx, encodePostDelete(input))
	if err != nil {
		return nil, err
	}
	return decodePost(post), nil
}
func (s *PostServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.PostService_ServiceDesc, s)
	return nil
}
