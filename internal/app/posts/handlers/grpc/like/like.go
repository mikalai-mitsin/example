package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
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
) (*examplepb.Like, error) {
	like, err := s.likeUseCase.Delete(ctx, encodeLikeDelete(input))
	if err != nil {
		return nil, err
	}
	return decodeLike(like), nil
}
func (s *LikeServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.LikeService_ServiceDesc, s)
	return nil
}
