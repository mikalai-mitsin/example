package handlers

import (
	"context"

	"github.com/mikalai-mitsin/example/internal/app/session/models"
	"github.com/mikalai-mitsin/example/internal/pkg/pointer"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
	examplepb "github.com/mikalai-mitsin/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type SessionServiceServer struct {
	examplepb.UnimplementedSessionServiceServer
	sessionInterceptor SessionInterceptor
	logger             Logger
}

func NewSessionServiceServer(
	sessionInterceptor SessionInterceptor,
	logger Logger,
) *SessionServiceServer {
	return &SessionServiceServer{sessionInterceptor: sessionInterceptor, logger: logger}
}

func (s *SessionServiceServer) Create(
	ctx context.Context,
	input *examplepb.SessionCreate,
) (*examplepb.Session, error) {
	item, err := s.sessionInterceptor.Create(ctx, encodeSessionCreate(input))
	if err != nil {
		return nil, err
	}
	return decodeSession(item), nil
}

func (s *SessionServiceServer) Get(
	ctx context.Context,
	input *examplepb.SessionGet,
) (*examplepb.Session, error) {
	item, err := s.sessionInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, err
	}
	return decodeSession(item), nil
}

func (s *SessionServiceServer) List(
	ctx context.Context,
	filter *examplepb.SessionFilter,
) (*examplepb.ListSession, error) {
	items, count, err := s.sessionInterceptor.List(ctx, encodeSessionFilter(filter))
	if err != nil {
		return nil, err
	}
	return decodeListSession(items, count), nil
}

func (s *SessionServiceServer) Update(
	ctx context.Context,
	input *examplepb.SessionUpdate,
) (*examplepb.Session, error) {
	item, err := s.sessionInterceptor.Update(ctx, encodeSessionUpdate(input))
	if err != nil {
		return nil, err
	}
	return decodeSession(item), nil
}

func (s *SessionServiceServer) Delete(
	ctx context.Context,
	input *examplepb.SessionDelete,
) (*emptypb.Empty, error) {
	if err := s.sessionInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
func encodeSessionCreate(input *examplepb.SessionCreate) *models.SessionCreate {
	create := &models.SessionCreate{Title: input.GetTitle(), Description: input.GetDescription()}
	return create
}
func encodeSessionFilter(input *examplepb.SessionFilter) *models.SessionFilter {
	filter := &models.SessionFilter{
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
func encodeSessionUpdate(input *examplepb.SessionUpdate) *models.SessionUpdate {
	update := &models.SessionUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetTitle() != nil {
		update.Title = pointer.Pointer(string(input.GetTitle().GetValue()))
	}
	if input.GetDescription() != nil {
		update.Description = pointer.Pointer(string(input.GetDescription().GetValue()))
	}
	return update
}
func decodeSession(item *models.Session) *examplepb.Session {
	response := &examplepb.Session{
		Id:          string(item.ID),
		CreatedAt:   timestamppb.New(item.CreatedAt),
		UpdatedAt:   timestamppb.New(item.UpdatedAt),
		Title:       item.Title,
		Description: item.Description,
	}
	return response
}
func decodeListSession(items []*models.Session, count uint64) *examplepb.ListSession {
	response := &examplepb.ListSession{
		Items: make([]*examplepb.Session, 0, len(items)),
		Count: count,
	}
	for _, item := range items {
		response.Items = append(response.Items, decodeSession(item))
	}
	return response
}
func decodeSessionUpdate(update *models.SessionUpdate) *examplepb.SessionUpdate {
	result := &examplepb.SessionUpdate{
		Id:          string(string(update.ID)),
		Title:       wrapperspb.String(*update.Title),
		Description: wrapperspb.String(*update.Description),
	}
	return result
}
