package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/pkg/grpc"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
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
) (*examplepb.Tag, error) {
	tag, err := s.tagUseCase.Delete(ctx, uuid.MustParse(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeTag(tag), nil
}
func (s *TagServiceServer) RegisterGRPC(grpcServer *grpc.Server) error {
	grpcServer.AddHandler(&examplepb.TagService_ServiceDesc, s)
	return nil
}
