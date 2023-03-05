package grpc

import (
	"context"
	"fmt"

	"github.com/018bf/example/internal/domain/interceptors"
	"github.com/018bf/example/internal/domain/models"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"github.com/018bf/example/pkg/log"
	"github.com/018bf/example/pkg/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type ArchServiceServer struct {
	examplepb.UnimplementedArchServiceServer
	archInterceptor interceptors.ArchInterceptor
	logger          log.Logger
}

func NewArchServiceServer(
	archInterceptor interceptors.ArchInterceptor,
	logger log.Logger,
) examplepb.ArchServiceServer {
	return &ArchServiceServer{archInterceptor: archInterceptor, logger: logger}
}

func (s *ArchServiceServer) Create(
	ctx context.Context,
	input *examplepb.ArchCreate,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Create(
		ctx,
		encodeArchCreate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) Get(
	ctx context.Context,
	input *examplepb.ArchGet,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Get(
		ctx,
		models.UUID(input.GetId()),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) List(
	ctx context.Context,
	filter *examplepb.ArchFilter,
) (*examplepb.ListArch, error) {
	listArches, count, err := s.archInterceptor.List(
		ctx,
		encodeArchFilter(filter),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	header := metadata.Pairs("count", fmt.Sprint(count))
	_ = grpc.SendHeader(ctx, header)
	return decodeListArch(listArches, count), nil
}

func (s *ArchServiceServer) Update(
	ctx context.Context,
	input *examplepb.ArchUpdate,
) (*examplepb.Arch, error) {
	arch, err := s.archInterceptor.Update(
		ctx,
		encodeArchUpdate(input),
		ctx.Value(UserKey).(*models.User),
	)
	if err != nil {
		return nil, decodeError(err)
	}
	return decodeArch(arch), nil
}

func (s *ArchServiceServer) Delete(
	ctx context.Context,
	input *examplepb.ArchDelete,
) (*emptypb.Empty, error) {
	if err := s.archInterceptor.Delete(ctx, models.UUID(input.GetId()), ctx.Value(UserKey).(*models.User)); err != nil {
		return nil, decodeError(err)
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
		filter.PageSize = utils.Pointer(input.GetPageSize().GetValue())
	}
	if input.GetPageNumber() != nil {
		filter.PageNumber = utils.Pointer(input.GetPageNumber().GetValue())
	}
	for _, id := range input.GetIds() {
		filter.IDs = append(filter.IDs, models.UUID(id))
	}
	if input.GetSearch() != nil {
		filter.Search = utils.Pointer(input.GetSearch().GetValue())
	}
	return filter
}
func encodeArchUpdate(input *examplepb.ArchUpdate) *models.ArchUpdate {
	update := &models.ArchUpdate{ID: models.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = utils.Pointer(input.GetName().GetValue())
	}
	if input.GetTitle() != nil {
		update.Title = utils.Pointer(input.GetTitle().GetValue())
	}
	if input.GetSubtitle() != nil {
		update.Subtitle = utils.Pointer(input.GetSubtitle().GetValue())
	}
	if input.GetTags() != nil {
		var params []string
		for _, item := range input.GetTags().GetValues() {
			params = append(params, string(item.GetStringValue()))
		}
		update.Tags = utils.Pointer(params)
	}
	if input.GetVersions() != nil {
		var params []uint
		for _, item := range input.GetVersions().GetValues() {
			params = append(params, uint(item.GetNumberValue()))
		}
		update.Versions = utils.Pointer(params)
	}
	if input.GetOldVersions() != nil {
		var params []uint64
		for _, item := range input.GetOldVersions().GetValues() {
			params = append(params, uint64(item.GetNumberValue()))
		}
		update.OldVersions = utils.Pointer(params)
	}
	if input.GetRelease() != nil {
		update.Release = utils.Pointer(input.GetRelease().AsTime())
	}
	if input.GetTested() != nil {
		update.Tested = utils.Pointer(input.GetTested().AsTime())
	}
	if input.GetMark() != nil {
		update.Mark = utils.Pointer(input.GetMark().GetValue())
	}
	if input.GetSubmarine() != nil {
		update.Submarine = utils.Pointer(input.GetSubmarine().GetValue())
	}
	if input.GetNumb() != nil {
		update.Numb = utils.Pointer(input.GetNumb().GetValue())
	}
	return update
}
func decodeArch(arch *models.Arch) *examplepb.Arch {
	response := &examplepb.Arch{
		Id:          string(arch.ID),
		UpdatedAt:   timestamppb.New(arch.UpdatedAt),
		CreatedAt:   timestamppb.New(arch.CreatedAt),
		Name:        arch.Name,
		Title:       arch.Title,
		Subtitle:    arch.Subtitle,
		Tags:        arch.Tags,
		Versions:    utils.ChangeType[uint32, uint](arch.Versions),
		OldVersions: arch.OldVersions,
		Release:     timestamppb.New(arch.Release),
		Tested:      timestamppb.New(arch.Tested),
		Mark:        arch.Mark,
		Submarine:   arch.Submarine,
		Numb:        arch.Numb,
	}
	return response
}
func decodeListArch(listArches []*models.Arch, count uint64) *examplepb.ListArch {
	response := &examplepb.ListArch{
		Items: make([]*examplepb.Arch, 0, len(listArches)),
		Count: count,
	}
	for _, arch := range listArches {
		response.Items = append(response.Items, decodeArch(arch))
	}
	return response
}
func decodeArchUpdate(update *models.ArchUpdate) *examplepb.ArchUpdate {
	result := &examplepb.ArchUpdate{
		Id:          string(update.ID),
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
		params, err := structpb.NewList(utils.ToAnySlice(*update.Tags))
		if err != nil {
			return nil
		}
		result.Tags = params
	}
	if update.Versions != nil {
		params, err := structpb.NewList(utils.ToAnySlice(*update.Versions))
		if err != nil {
			return nil
		}
		result.Versions = params
	}
	if update.OldVersions != nil {
		params, err := structpb.NewList(utils.ToAnySlice(*update.OldVersions))
		if err != nil {
			return nil
		}
		result.OldVersions = params
	}
	return result
}
