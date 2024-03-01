package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/day/models"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type DayServiceServer struct {
	examplepb.UnimplementedDayServiceServer
	dayInterceptor DayInterceptor
	logger         log.Logger
}

func NewDayServiceServer(
	dayInterceptor DayInterceptor,
	logger log.Logger,
) examplepb.DayServiceServer {
	return &DayServiceServer{dayInterceptor: dayInterceptor, logger: logger}
}

func (s *DayServiceServer) Create(
	ctx context.Context,
	input *examplepb.DayCreate,
) (*examplepb.Day, error) {
	item, err := s.dayInterceptor.Create(ctx, encodeDayCreate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeDay(item), nil
}

func (s *DayServiceServer) Get(
	ctx context.Context,
	input *examplepb.DayGet,
) (*examplepb.Day, error) {
	item, err := s.dayInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeDay(item), nil
}

func (s *DayServiceServer) List(
	ctx context.Context,
	filter *examplepb.DayFilter,
) (*examplepb.ListDay, error) {
	items, count, err := s.dayInterceptor.List(ctx, encodeDayFilter(filter))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeListDay(items, count), nil
}

func (s *DayServiceServer) Update(
	ctx context.Context,
	input *examplepb.DayUpdate,
) (*examplepb.Day, error) {
	item, err := s.dayInterceptor.Update(ctx, encodeDayUpdate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeDay(item), nil
}

func (s *DayServiceServer) Delete(
	ctx context.Context,
	input *examplepb.DayDelete,
) (*emptypb.Empty, error) {
	if err := s.dayInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, grpc.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeDayCreate(input *examplepb.DayCreate) *models.DayCreate {
	create := &models.DayCreate{
		Name:        input.GetName(),
		Repeat:      int(input.GetRepeat()),
		EquipmentID: input.GetEquipmentId(),
	}
	return create
}
func encodeDayFilter(input *examplepb.DayFilter) *models.DayFilter {
	filter := &models.DayFilter{
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
func encodeDayUpdate(input *examplepb.DayUpdate) *models.DayUpdate {
	update := &models.DayUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = pointer.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetRepeat() != nil {
		update.Repeat = pointer.Pointer(int(input.GetRepeat().GetValue()))
	}
	if input.GetEquipmentId() != nil {
		update.EquipmentID = pointer.Pointer(string(input.GetEquipmentId().GetValue()))
	}
	return update
}
func decodeDay(item *models.Day) *examplepb.Day {
	response := &examplepb.Day{
		Id:          string(item.ID),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Name:        item.Name,
		Repeat:      int32(item.Repeat),
		EquipmentId: item.EquipmentID,
	}
	return response
}
func decodeListDay(items []*models.Day, count uint64) *examplepb.ListDay {
	response := &examplepb.ListDay{Items: make([]*examplepb.Day, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeDay(item))
	}
	return response
}
func decodeDayUpdate(update *models.DayUpdate) *examplepb.DayUpdate {
	result := &examplepb.DayUpdate{
		Id:          string(string(update.ID)),
		Name:        wrapperspb.String(*update.Name),
		Repeat:      wrapperspb.Int32(int32(*update.Repeat)),
		EquipmentId: wrapperspb.String(*update.EquipmentID),
	}
	return result
}
