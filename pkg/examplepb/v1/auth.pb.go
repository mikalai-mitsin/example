// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        (unknown)
// source: examplepb/v1/auth.proto

package v1

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CreateToken struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Email         string                 `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	Password      string                 `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *CreateToken) Reset() {
	*x = CreateToken{}
	mi := &file_examplepb_v1_auth_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CreateToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateToken) ProtoMessage() {}

func (x *CreateToken) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_auth_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateToken.ProtoReflect.Descriptor instead.
func (*CreateToken) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_auth_proto_rawDescGZIP(), []int{0}
}

func (x *CreateToken) GetEmail() string {
	if x != nil {
		return x.Email
	}
	return ""
}

func (x *CreateToken) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type AccessToken struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *AccessToken) Reset() {
	*x = AccessToken{}
	mi := &file_examplepb_v1_auth_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AccessToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccessToken) ProtoMessage() {}

func (x *AccessToken) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_auth_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccessToken.ProtoReflect.Descriptor instead.
func (*AccessToken) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_auth_proto_rawDescGZIP(), []int{1}
}

func (x *AccessToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RefreshToken struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RefreshToken) Reset() {
	*x = RefreshToken{}
	mi := &file_examplepb_v1_auth_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RefreshToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RefreshToken) ProtoMessage() {}

func (x *RefreshToken) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_auth_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RefreshToken.ProtoReflect.Descriptor instead.
func (*RefreshToken) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_auth_proto_rawDescGZIP(), []int{2}
}

func (x *RefreshToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type RevokeToken struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Token         string                 `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *RevokeToken) Reset() {
	*x = RevokeToken{}
	mi := &file_examplepb_v1_auth_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RevokeToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RevokeToken) ProtoMessage() {}

func (x *RevokeToken) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_auth_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RevokeToken.ProtoReflect.Descriptor instead.
func (*RevokeToken) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_auth_proto_rawDescGZIP(), []int{3}
}

func (x *RevokeToken) GetToken() string {
	if x != nil {
		return x.Token
	}
	return ""
}

type TokenPair struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Access        string                 `protobuf:"bytes,1,opt,name=access,proto3" json:"access,omitempty"`
	Refresh       string                 `protobuf:"bytes,2,opt,name=refresh,proto3" json:"refresh,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TokenPair) Reset() {
	*x = TokenPair{}
	mi := &file_examplepb_v1_auth_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TokenPair) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenPair) ProtoMessage() {}

func (x *TokenPair) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_auth_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenPair.ProtoReflect.Descriptor instead.
func (*TokenPair) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_auth_proto_rawDescGZIP(), []int{4}
}

func (x *TokenPair) GetAccess() string {
	if x != nil {
		return x.Access
	}
	return ""
}

func (x *TokenPair) GetRefresh() string {
	if x != nil {
		return x.Refresh
	}
	return ""
}

var File_examplepb_v1_auth_proto protoreflect.FileDescriptor

var file_examplepb_v1_auth_proto_rawDesc = string([]byte{
	0x0a, 0x17, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x3f, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x23, 0x0a, 0x0b, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x24, 0x0a, 0x0c, 0x52,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x74,
	0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x6b, 0x65,
	0x6e, 0x22, 0x23, 0x0a, 0x0b, 0x52, 0x65, 0x76, 0x6f, 0x6b, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e,
	0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x22, 0x3d, 0x0a, 0x09, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x50,
	0x61, 0x69, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x72,
	0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65,
	0x66, 0x72, 0x65, 0x73, 0x68, 0x32, 0xc7, 0x01, 0x0a, 0x0b, 0x41, 0x75, 0x74, 0x68, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x5a, 0x0a, 0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x19, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62,
	0x2e, 0x76, 0x31, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a,
	0x17, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x54,
	0x6f, 0x6b, 0x65, 0x6e, 0x50, 0x61, 0x69, 0x72, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11,
	0x3a, 0x01, 0x2a, 0x22, 0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74,
	0x68, 0x12, 0x5c, 0x0a, 0x0c, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x1a, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x66, 0x72, 0x65, 0x73, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x1a, 0x17, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x50, 0x61, 0x69, 0x72, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x3a, 0x01,
	0x2a, 0x32, 0x0c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x42,
	0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x69,
	0x6b, 0x61, 0x6c, 0x61, 0x69, 0x2d, 0x6d, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x2f, 0x65, 0x78, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x70, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
})

var (
	file_examplepb_v1_auth_proto_rawDescOnce sync.Once
	file_examplepb_v1_auth_proto_rawDescData []byte
)

func file_examplepb_v1_auth_proto_rawDescGZIP() []byte {
	file_examplepb_v1_auth_proto_rawDescOnce.Do(func() {
		file_examplepb_v1_auth_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_examplepb_v1_auth_proto_rawDesc), len(file_examplepb_v1_auth_proto_rawDesc)))
	})
	return file_examplepb_v1_auth_proto_rawDescData
}

var file_examplepb_v1_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_examplepb_v1_auth_proto_goTypes = []any{
	(*CreateToken)(nil),  // 0: examplepb.v1.CreateToken
	(*AccessToken)(nil),  // 1: examplepb.v1.AccessToken
	(*RefreshToken)(nil), // 2: examplepb.v1.RefreshToken
	(*RevokeToken)(nil),  // 3: examplepb.v1.RevokeToken
	(*TokenPair)(nil),    // 4: examplepb.v1.TokenPair
}
var file_examplepb_v1_auth_proto_depIdxs = []int32{
	0, // 0: examplepb.v1.AuthService.CreateToken:input_type -> examplepb.v1.CreateToken
	2, // 1: examplepb.v1.AuthService.RefreshToken:input_type -> examplepb.v1.RefreshToken
	4, // 2: examplepb.v1.AuthService.CreateToken:output_type -> examplepb.v1.TokenPair
	4, // 3: examplepb.v1.AuthService.RefreshToken:output_type -> examplepb.v1.TokenPair
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_examplepb_v1_auth_proto_init() }
func file_examplepb_v1_auth_proto_init() {
	if File_examplepb_v1_auth_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_examplepb_v1_auth_proto_rawDesc), len(file_examplepb_v1_auth_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examplepb_v1_auth_proto_goTypes,
		DependencyIndexes: file_examplepb_v1_auth_proto_depIdxs,
		MessageInfos:      file_examplepb_v1_auth_proto_msgTypes,
	}.Build()
	File_examplepb_v1_auth_proto = out.File
	file_examplepb_v1_auth_proto_goTypes = nil
	file_examplepb_v1_auth_proto_depIdxs = nil
}
