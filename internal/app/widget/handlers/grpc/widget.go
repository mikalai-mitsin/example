package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/widget/entities"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type WidgetServiceServer struct {
	examplepb.UnimplementedWidgetServiceServer
	widgetUseCase widgetUseCase
	logger        logger
}

func NewWidgetServiceServer(widgetUseCase widgetUseCase, logger logger) *WidgetServiceServer {
	return &WidgetServiceServer{widgetUseCase: widgetUseCase, logger: logger}
}

func (s *WidgetServiceServer) Create(
	ctx context.Context,
	input *examplepb.WidgetCreate,
) (*examplepb.Widget, error) {
	item, err := s.widgetUseCase.Create(ctx, encodeWidgetCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeWidget(item), nil
}

func (s *WidgetServiceServer) Get(
	ctx context.Context,
	input *examplepb.WidgetGet,
) (*examplepb.Widget, error) {
	item, err := s.widgetUseCase.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeWidget(item), nil
}

func (s *WidgetServiceServer) List(
	ctx context.Context,
	filter *examplepb.WidgetFilter,
) (*examplepb.ListWidget, error) {
	items, count, err := s.widgetUseCase.List(ctx, encodeWidgetFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListWidget(items, count), nil
}

func (s *WidgetServiceServer) Update(
	ctx context.Context,
	input *examplepb.WidgetUpdate,
) (*examplepb.Widget, error) {
	item, err := s.widgetUseCase.Update(ctx, encodeWidgetUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeWidget(item), nil
}

func (s *WidgetServiceServer) Delete(
	ctx context.Context,
	input *examplepb.WidgetDelete,
) (*emptypb.Empty, error) {
	if err := s.widgetUseCase.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodeWidgetCreate(input *examplepb.WidgetCreate) *entities.WidgetCreate {
	create := &entities.WidgetCreate{
		FormScreenId: input.GetFormScreenId(),
		Name:         input.GetName(),
		Ordering:     input.GetOrdering(),
		IsOptional:   input.GetIsOptional(),
		UiSettings:   input.GetUiSettings(),
		DeletedAt:    input.GetDeletedAt().AsTime(),
	}
	return create
}
func encodeWidgetFilter(input *examplepb.WidgetFilter) *entities.WidgetFilter {
	filter := &entities.WidgetFilter{
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
func encodeWidgetUpdate(input *examplepb.WidgetUpdate) *entities.WidgetUpdate {
	update := &entities.WidgetUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetFormScreenId() != nil {
		update.FormScreenId = pointer.Pointer(string(input.GetFormScreenId().GetValue()))
	}
	if input.GetName() != nil {
		update.Name = pointer.Pointer(string(input.GetName().GetValue()))
	}
	if input.GetOrdering() != nil {
		update.Ordering = pointer.Pointer(int64(input.GetOrdering().GetValue()))
	}
	if input.GetIsOptional() != nil {
		update.IsOptional = pointer.Pointer(bool(input.GetIsOptional().GetValue()))
	}
	if input.GetUiSettings() != nil {
		update.UiSettings = pointer.Pointer(string(input.GetUiSettings().GetValue()))
	}
	if input.GetDeletedAt() != nil {
		update.DeletedAt = pointer.Pointer(input.GetDeletedAt().AsTime())
	}
	return update
}
func decodeWidget(item *entities.Widget) *examplepb.Widget {
	response := &examplepb.Widget{
		Id:           string(item.ID),
		CreatedAt:    timestamppb.New(item.CreatedAt),
		UpdatedAt:    timestamppb.New(item.UpdatedAt),
		FormScreenId: item.FormScreenId,
		Name:         item.Name,
		Ordering:     item.Ordering,
		IsOptional:   item.IsOptional,
		UiSettings:   item.UiSettings,
		DeletedAt:    timestamppb.New(item.DeletedAt),
	}
	return response
}
func decodeListWidget(items []*entities.Widget, count uint64) *examplepb.ListWidget {
	response := &examplepb.ListWidget{Items: make([]*examplepb.Widget, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeWidget(item))
	}
	return response
}
func decodeWidgetUpdate(update *entities.WidgetUpdate) *examplepb.WidgetUpdate {
	result := &examplepb.WidgetUpdate{
		Id:           string(string(update.ID)),
		FormScreenId: wrapperspb.String(*update.FormScreenId),
		Name:         wrapperspb.String(*update.Name),
		Ordering:     wrapperspb.Int64(*update.Ordering),
		IsOptional:   wrapperspb.Bool(*update.IsOptional),
		UiSettings:   wrapperspb.String(*update.UiSettings),
		DeletedAt:    timestamppb.New(*update.DeletedAt),
	}
	return result
}
