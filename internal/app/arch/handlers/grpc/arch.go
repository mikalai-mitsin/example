package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/arch/models"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ArchServiceServer struct {
	examplepb.UnimplementedArchServiceServer
	archInterceptor ArchInterceptor
	logger          log.Logger
}

func NewArchServiceServer(
	archInterceptor ArchInterceptor,
	logger log.Logger,
) examplepb.ArchServiceServer {
	return &ArchServiceServer{archInterceptor: archInterceptor, logger: logger}
}

func (s *ArchServiceServer) Create(
	ctx context.Context,
	input *examplepb.ArchCreate,
) (*examplepb.Arch, error) {
	item, err := s.archInterceptor.Create(ctx, encodeArchCreate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeArch(item), nil
}

func (s *ArchServiceServer) Get(
	ctx context.Context,
	input *examplepb.ArchGet,
) (*examplepb.Arch, error) {
	item, err := s.archInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeArch(item), nil
}

func (s *ArchServiceServer) List(
	ctx context.Context,
	filter *examplepb.ArchFilter,
) (*examplepb.ListArch, error) {
	items, count, err := s.archInterceptor.List(ctx, encodeArchFilter(filter))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeListArch(items, count), nil
}

func (s *ArchServiceServer) Update(
	ctx context.Context,
	input *examplepb.ArchUpdate,
) (*examplepb.Arch, error) {
	item, err := s.archInterceptor.Update(ctx, encodeArchUpdate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeArch(item), nil
}

func (s *ArchServiceServer) Delete(
	ctx context.Context,
	input *examplepb.ArchDelete,
) (*emptypb.Empty, error) {
	if err := s.archInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, grpc.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeArchCreate(input *examplepb.ArchCreate) *models.ArchCreate {
	create := &models.ArchCreate{
		Name:        input.GetName(),
		Title:       input.GetTitle(),
		Subtitle:    input.GetSubtitle(),
		Tags:        input.GetTags(),
		Versions:    nil,
		OldVersions: input.GetOldVersions(),
		Release:     input.GetRelease().AsTime(),
		Tested:      input.GetTested().AsTime(),
		Mark:        input.GetMark(),
		Submarine:   input.GetSubmarine(),
		Numb:        input.GetNumb(),
	}
	return create
}
func encodeArchFilter(input *examplepb.ArchFilter) *models.ArchFilter {
	filter := &models.ArchFilter{
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
	if input.GetSearch() != nil {
		filter.Search = pointer.Pointer(input.GetSearch().GetValue())
	}
	return filter
}
func encodeArchUpdate(input *examplepb.ArchUpdate) *models.ArchUpdate {
	update := &models.ArchUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = pointer.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetTitle() != nil {
		update.Title = pointer.Pointer(string(input.GetTitle().GetValue()))
	}
	if input.GetSubtitle() != nil {
		update.Subtitle = pointer.Pointer(string(input.GetSubtitle().GetValue()))
	}
	if input.GetTags() != nil {
		var params []string
		for _, item := range input.GetTags().GetValues() {
			params = append(params, string(item.GetStringValue()))
		}
		update.Tags = pointer.Pointer(params)
	}
	if input.GetVersions() != nil {
		var params []uint
		for _, item := range input.GetVersions().GetValues() {
			params = append(params, uint(item.GetNumberValue()))
		}
		update.Versions = pointer.Pointer(params)
	}
	if input.GetOldVersions() != nil {
		var params []uint64
		for _, item := range input.GetOldVersions().GetValues() {
			params = append(params, uint64(item.GetNumberValue()))
		}
		update.OldVersions = pointer.Pointer(params)
	}
	if input.GetRelease() != nil {
		update.Release = pointer.Pointer(input.GetRelease().AsTime())
	}
	if input.GetTested() != nil {
		update.Tested = pointer.Pointer(input.GetTested().AsTime())
	}
	if input.GetMark() != nil {
		update.Mark = pointer.Pointer(string(input.GetMark().GetValue()))
	}
	if input.GetSubmarine() != nil {
		update.Submarine = pointer.Pointer(string(input.GetSubmarine().GetValue()))
	}
	if input.GetNumb() != nil {
		update.Numb = pointer.Pointer(uint64(input.GetNumb().GetValue()))
	}
	return update
}
func decodeArch(item *models.Arch) *examplepb.Arch {
	response := &examplepb.Arch{
		Id:          string(item.ID),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Name:        item.Name,
		Title:       item.Title,
		Subtitle:    item.Subtitle,
		Tags:        item.Tags,
		Versions:    pointer.ChangeType[uint32, uint](item.Versions),
		OldVersions: item.OldVersions,
		Release:     timestamppb.New(item.Release),
		Tested:      timestamppb.New(item.Tested),
		Mark:        item.Mark,
		Submarine:   item.Submarine,
		Numb:        item.Numb,
	}
	return response
}
func decodeListArch(items []*models.Arch, count uint64) *examplepb.ListArch {
	response := &examplepb.ListArch{Items: make([]*examplepb.Arch, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeArch(item))
	}
	return response
}
func decodeArchUpdate(update *models.ArchUpdate) *examplepb.ArchUpdate {
	result := &examplepb.ArchUpdate{
		Id:          string(string(update.ID)),
		Name:        wrapperspb.String(*update.Name),
		Title:       wrapperspb.String(*update.Title),
		Subtitle:    wrapperspb.String(*update.Subtitle),
		Tags:        nil,
		Versions:    nil,
		OldVersions: nil,
		Release:     timestamppb.New(*update.Release),
		Tested:      timestamppb.New(*update.Tested),
		Mark:        wrapperspb.String(*update.Mark),
		Submarine:   wrapperspb.String(*update.Submarine),
		Numb:        wrapperspb.UInt64(*update.Numb),
	}
	if update.Tags != nil {
		params, err := structpb.NewList(pointer.ToAnySlice(*update.Tags))
		if err != nil {
			return nil
		}
		result.Tags = params
	}
	if update.Versions != nil {
		params, err := structpb.NewList(pointer.ToAnySlice(*update.Versions))
		if err != nil {
			return nil
		}
		result.Versions = params
	}
	if update.OldVersions != nil {
		params, err := structpb.NewList(pointer.ToAnySlice(*update.OldVersions))
		if err != nil {
			return nil
		}
		result.OldVersions = params
	}
	return result
}
