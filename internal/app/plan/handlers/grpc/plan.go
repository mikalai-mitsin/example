package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/plan/models"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type PlanServiceServer struct {
	examplepb.UnimplementedPlanServiceServer
	planInterceptor PlanInterceptor
	logger          log.Logger
}

func NewPlanServiceServer(
	planInterceptor PlanInterceptor,
	logger log.Logger,
) examplepb.PlanServiceServer {
	return &PlanServiceServer{planInterceptor: planInterceptor, logger: logger}
}

func (s *PlanServiceServer) Create(
	ctx context.Context,
	input *examplepb.PlanCreate,
) (*examplepb.Plan, error) {
	item, err := s.planInterceptor.Create(ctx, encodePlanCreate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodePlan(item), nil
}

func (s *PlanServiceServer) Get(
	ctx context.Context,
	input *examplepb.PlanGet,
) (*examplepb.Plan, error) {
	item, err := s.planInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodePlan(item), nil
}

func (s *PlanServiceServer) List(
	ctx context.Context,
	filter *examplepb.PlanFilter,
) (*examplepb.ListPlan, error) {
	items, count, err := s.planInterceptor.List(ctx, encodePlanFilter(filter))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeListPlan(items, count), nil
}

func (s *PlanServiceServer) Update(
	ctx context.Context,
	input *examplepb.PlanUpdate,
) (*examplepb.Plan, error) {
	item, err := s.planInterceptor.Update(ctx, encodePlanUpdate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodePlan(item), nil
}

func (s *PlanServiceServer) Delete(
	ctx context.Context,
	input *examplepb.PlanDelete,
) (*emptypb.Empty, error) {
	if err := s.planInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, grpc.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodePlanCreate(input *examplepb.PlanCreate) *models.PlanCreate {
	create := &models.PlanCreate{
		Name:        input.GetName(),
		Repeat:      input.GetRepeat(),
		EquipmentID: input.GetEquipmentId(),
	}
	return create
}
func encodePlanFilter(input *examplepb.PlanFilter) *models.PlanFilter {
	filter := &models.PlanFilter{
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
func encodePlanUpdate(input *examplepb.PlanUpdate) *models.PlanUpdate {
	update := &models.PlanUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetName() != nil {
		update.Name = pointer.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetRepeat() != nil {
		update.Repeat = pointer.Pointer(uint64(input.GetRepeat().GetValue()))
	}
	if input.GetEquipmentId() != nil {
		update.EquipmentID = pointer.Pointer(string(input.GetEquipmentId().GetValue()))
	}
	return update
}
func decodePlan(item *models.Plan) *examplepb.Plan {
	response := &examplepb.Plan{
		Id:          string(item.ID),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Name:        item.Name,
		Repeat:      item.Repeat,
		EquipmentId: item.EquipmentID,
	}
	return response
}
func decodeListPlan(items []*models.Plan, count uint64) *examplepb.ListPlan {
	response := &examplepb.ListPlan{Items: make([]*examplepb.Plan, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodePlan(item))
	}
	return response
}
func decodePlanUpdate(update *models.PlanUpdate) *examplepb.PlanUpdate {
	result := &examplepb.PlanUpdate{
		Id:          string(string(update.ID)),
		Name:        wrapperspb.String(*update.Name),
		Repeat:      wrapperspb.UInt64(*update.Repeat),
		EquipmentId: wrapperspb.String(*update.EquipmentID),
	}
	return result
}
