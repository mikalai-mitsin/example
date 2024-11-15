// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: examplepb/v1/widget.proto

package v1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	WidgetService_Create_FullMethodName = "/examplepb.v1.WidgetService/Create"
	WidgetService_Get_FullMethodName    = "/examplepb.v1.WidgetService/Get"
	WidgetService_Update_FullMethodName = "/examplepb.v1.WidgetService/Update"
	WidgetService_Delete_FullMethodName = "/examplepb.v1.WidgetService/Delete"
	WidgetService_List_FullMethodName   = "/examplepb.v1.WidgetService/List"
)

// WidgetServiceClient is the client API for WidgetService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type WidgetServiceClient interface {
	Create(ctx context.Context, in *WidgetCreate, opts ...grpc.CallOption) (*Widget, error)
	Get(ctx context.Context, in *WidgetGet, opts ...grpc.CallOption) (*Widget, error)
	Update(ctx context.Context, in *WidgetUpdate, opts ...grpc.CallOption) (*Widget, error)
	Delete(ctx context.Context, in *WidgetDelete, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *WidgetFilter, opts ...grpc.CallOption) (*ListWidget, error)
}

type widgetServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewWidgetServiceClient(cc grpc.ClientConnInterface) WidgetServiceClient {
	return &widgetServiceClient{cc}
}

func (c *widgetServiceClient) Create(ctx context.Context, in *WidgetCreate, opts ...grpc.CallOption) (*Widget, error) {
	out := new(Widget)
	err := c.cc.Invoke(ctx, WidgetService_Create_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *widgetServiceClient) Get(ctx context.Context, in *WidgetGet, opts ...grpc.CallOption) (*Widget, error) {
	out := new(Widget)
	err := c.cc.Invoke(ctx, WidgetService_Get_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *widgetServiceClient) Update(ctx context.Context, in *WidgetUpdate, opts ...grpc.CallOption) (*Widget, error) {
	out := new(Widget)
	err := c.cc.Invoke(ctx, WidgetService_Update_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *widgetServiceClient) Delete(ctx context.Context, in *WidgetDelete, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, WidgetService_Delete_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *widgetServiceClient) List(ctx context.Context, in *WidgetFilter, opts ...grpc.CallOption) (*ListWidget, error) {
	out := new(ListWidget)
	err := c.cc.Invoke(ctx, WidgetService_List_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// WidgetServiceServer is the server API for WidgetService service.
// All implementations should embed UnimplementedWidgetServiceServer
// for forward compatibility
type WidgetServiceServer interface {
	Create(context.Context, *WidgetCreate) (*Widget, error)
	Get(context.Context, *WidgetGet) (*Widget, error)
	Update(context.Context, *WidgetUpdate) (*Widget, error)
	Delete(context.Context, *WidgetDelete) (*emptypb.Empty, error)
	List(context.Context, *WidgetFilter) (*ListWidget, error)
}

// UnimplementedWidgetServiceServer should be embedded to have forward compatible implementations.
type UnimplementedWidgetServiceServer struct {
}

func (UnimplementedWidgetServiceServer) Create(context.Context, *WidgetCreate) (*Widget, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedWidgetServiceServer) Get(context.Context, *WidgetGet) (*Widget, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedWidgetServiceServer) Update(context.Context, *WidgetUpdate) (*Widget, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedWidgetServiceServer) Delete(context.Context, *WidgetDelete) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedWidgetServiceServer) List(context.Context, *WidgetFilter) (*ListWidget, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}

// UnsafeWidgetServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to WidgetServiceServer will
// result in compilation errors.
type UnsafeWidgetServiceServer interface {
	mustEmbedUnimplementedWidgetServiceServer()
}

func RegisterWidgetServiceServer(s grpc.ServiceRegistrar, srv WidgetServiceServer) {
	s.RegisterService(&WidgetService_ServiceDesc, srv)
}

func _WidgetService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WidgetCreate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WidgetServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WidgetService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WidgetServiceServer).Create(ctx, req.(*WidgetCreate))
	}
	return interceptor(ctx, in, info, handler)
}

func _WidgetService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WidgetGet)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WidgetServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WidgetService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WidgetServiceServer).Get(ctx, req.(*WidgetGet))
	}
	return interceptor(ctx, in, info, handler)
}

func _WidgetService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WidgetUpdate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WidgetServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WidgetService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WidgetServiceServer).Update(ctx, req.(*WidgetUpdate))
	}
	return interceptor(ctx, in, info, handler)
}

func _WidgetService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WidgetDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WidgetServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WidgetService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WidgetServiceServer).Delete(ctx, req.(*WidgetDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _WidgetService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WidgetFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WidgetServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: WidgetService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WidgetServiceServer).List(ctx, req.(*WidgetFilter))
	}
	return interceptor(ctx, in, info, handler)
}

// WidgetService_ServiceDesc is the grpc.ServiceDesc for WidgetService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var WidgetService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "examplepb.v1.WidgetService",
	HandlerType: (*WidgetServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _WidgetService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _WidgetService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _WidgetService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _WidgetService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _WidgetService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "examplepb/v1/widget.proto",
}
