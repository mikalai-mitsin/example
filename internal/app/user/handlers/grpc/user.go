package grpc

import (
	"context"

	"github.com/018bf/example/internal/app/user/models"
	"github.com/018bf/example/internal/pkg/grpc"
	"github.com/018bf/example/internal/pkg/log"
	"github.com/018bf/example/internal/pkg/pointer"
	"github.com/018bf/example/internal/pkg/uuid"
	examplepb "github.com/018bf/example/pkg/examplepb/v1"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

type UserServiceServer struct {
	examplepb.UnimplementedUserServiceServer
	userInterceptor UserInterceptor
	logger          log.Logger
}

func NewUserServiceServer(
	userInterceptor UserInterceptor,
	logger log.Logger,
) examplepb.UserServiceServer {
	return &UserServiceServer{userInterceptor: userInterceptor, logger: logger}
}

func (s *UserServiceServer) Create(
	ctx context.Context,
	input *examplepb.UserCreate,
) (*examplepb.User, error) {
	item, err := s.userInterceptor.Create(ctx, encodeUserCreate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeUser(item), nil
}

func (s *UserServiceServer) Get(
	ctx context.Context,
	input *examplepb.UserGet,
) (*examplepb.User, error) {
	item, err := s.userInterceptor.Get(ctx, uuid.UUID(input.GetId()))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeUser(item), nil
}

func (s *UserServiceServer) List(
	ctx context.Context,
	filter *examplepb.UserFilter,
) (*examplepb.ListUser, error) {
	items, count, err := s.userInterceptor.List(ctx, encodeUserFilter(filter))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeListUser(items, count), nil
}

func (s *UserServiceServer) Update(
	ctx context.Context,
	input *examplepb.UserUpdate,
) (*examplepb.User, error) {
	item, err := s.userInterceptor.Update(ctx, encodeUserUpdate(input))
	if err != nil {
		return nil, grpc.DecodeError(err)
	}
	return decodeUser(item), nil
}

func (s *UserServiceServer) Delete(
	ctx context.Context,
	input *examplepb.UserDelete,
) (*emptypb.Empty, error) {
	if err := s.userInterceptor.Delete(ctx, uuid.UUID(input.GetId())); err != nil {
		return nil, grpc.DecodeError(err)
	}
	return &emptypb.Empty{}, nil
}
func encodeUserCreate(input *examplepb.UserCreate) *models.UserCreate {
	create := &models.UserCreate{
		FirstName: input.GetFirstName(),
		LastName:  input.GetLastName(),
		Password:  input.GetPassword(),
		Email:     input.GetEmail(),
		GroupID:   models.GroupID(input.GetGroupId()),
	}
	return create
}
func encodeUserFilter(input *examplepb.UserFilter) *models.UserFilter {
	filter := &models.UserFilter{
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
func encodeUserUpdate(input *examplepb.UserUpdate) *models.UserUpdate {
	update := &models.UserUpdate{ID: uuid.UUID(input.GetId())}
	if input.GetFirstName() != nil {
		update.FirstName = pointer.Pointer(string(input.GetFirstName().GetValue()))
	}
	if input.GetLastName() != nil {
		update.LastName = pointer.Pointer(string(input.GetLastName().GetValue()))
	}
	if input.GetPassword() != nil {
		update.Password = pointer.Pointer(string(input.GetPassword().GetValue()))
	}
	if input.GetEmail() != nil {
		update.Email = pointer.Pointer(string(input.GetEmail().GetValue()))
	}
	if input.GetGroupId() != nil {
		update.GroupID = pointer.Pointer(models.GroupID(input.GetGroupId().GetValue()))
	}
	return update
}
func decodeUser(item *models.User) *examplepb.User {
	response := &examplepb.User{
		Id:        string(item.ID),
		CreatedAt: timestamppb.New(item.CreatedAt),
		UpdatedAt: timestamppb.New(item.UpdatedAt),
		FirstName: item.FirstName,
		LastName:  item.LastName,
		Password:  item.Password,
		Email:     item.Email,
		GroupId:   string(item.GroupID),
	}
	return response
}
func decodeListUser(items []*models.User, count uint64) *examplepb.ListUser {
	response := &examplepb.ListUser{Items: make([]*examplepb.User, 0, len(items)), Count: count}
	for _, item := range items {
		response.Items = append(response.Items, decodeUser(item))
	}
	return response
}
func decodeUserUpdate(update *models.UserUpdate) *examplepb.UserUpdate {
	result := &examplepb.UserUpdate{
		Id:        string(string(update.ID)),
		FirstName: wrapperspb.String(*update.FirstName),
		LastName:  wrapperspb.String(*update.LastName),
		Password:  wrapperspb.String(*update.Password),
		Email:     wrapperspb.String(*update.Email),
		GroupId:   wrapperspb.String(string(*update.GroupID)),
	}
	return result
}
