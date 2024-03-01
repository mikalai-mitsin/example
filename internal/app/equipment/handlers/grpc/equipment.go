package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/equipment/models"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type EquipmentServiceServer struct {
	examplepb.UnimplementedEquipmentServiceServer
	equipmentInterceptor EquipmentInterceptor
	logger               log.Logger
}

func NewEquipmentServiceServer(
	equipmentInterceptor EquipmentInterceptor,
	logger log.Logger,
) examplepb.EquipmentServiceServer {
	return &EquipmentServiceServer{equipmentInterceptor: equipmentInterceptor, logger: logger}
}

func (s *EquipmentServiceServer) Create(
	ctx context.Context,
	input *examplepb.EquipmentCreate,
) (*examplepb.Equipment, error) {
	item, err := s.equipmentInterceptor.Create(ctx, encodeEquipmentCreate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeEquipment(item), nil
}

func (s *EquipmentServiceServer) Get(
	ctx context.Context,
	input *examplepb.EquipmentGet,
) (*examplepb.Equipment, error) {
	item, err := s.equipmentInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeEquipment(item), nil
}

func (s *EquipmentServiceServer) List(
	ctx context.Context,
	filter *examplepb.EquipmentFilter,
) (*examplepb.ListEquipment, error) {
	items, count, err := s.equipmentInterceptor.List(ctx, encodeEquipmentFilter(filter))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeListEquipment(items, count), nil
}

func (s *EquipmentServiceServer) Update(
	ctx context.Context,
	input *examplepb.EquipmentUpdate,
) (*examplepb.Equipment, error) {
	item, err := s.equipmentInterceptor.Update(ctx, encodeEquipmentUpdate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeEquipment(item), nil
}

func (s *EquipmentServiceServer) Delete(
	ctx context.Context,
	input *examplepb.EquipmentDelete,
) (*emptypb.Empty, error) {
	if err := s.equipmentInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, grpc.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeEquipmentCreate(input *examplepb.EquipmentCreate) *models.EquipmentCreate {
	create := &models.EquipmentCreate{
		Name:   input.GetName(),
		Repeat: int(input.GetRepeat()),
		Weight: int(input.GetWeight()),
	}
	return create
}
func encodeEquipmentFilter(input *examplepb.EquipmentFilter) *models.EquipmentFilter {
	filter := &models.EquipmentFilter{
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
func encodeEquipmentUpdate(input *examplepb.EquipmentUpdate) *models.EquipmentUpdate {
	update := &models.EquipmentUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = pointer.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetRepeat() != nil {
		update.Repeat = pointer.Pointer(int(input.GetRepeat().GetValue()))
	}
	if input.GetWeight() != nil {
		update.Weight = pointer.Pointer(int(input.GetWeight().GetValue()))
	}
	return update
}
func decodeEquipment(item *models.Equipment) *examplepb.Equipment {
	response := &examplepb.Equipment{
		Id:        string(item.ID),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		Name:      item.Name,
		Repeat:    int32(item.Repeat),
		Weight:    int32(item.Weight),
	}
	return response
}
func decodeListEquipment(items []*models.Equipment, count uint64) *examplepb.ListEquipment {
	response := &examplepb.ListEquipment{
		Items: make([]*examplepb.Equipment, 0, len(items)),
		Count: count,
	}
	for _, item := range items {
		response.Items = append(response.Items, decodeEquipment(item))
	}
	return response
}
func decodeEquipmentUpdate(update *models.EquipmentUpdate) *examplepb.EquipmentUpdate {
	result := &examplepb.EquipmentUpdate{
		Id:     string(string(update.ID)),
		Name:   wrapperspb.String(*update.Name),
		Repeat: wrapperspb.Int32(int32(*update.Repeat)),
		Weight: wrapperspb.Int32(int32(*update.Weight)),
	}
	return result
}
