// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        (unknown)
// source: examplepb/v1/arch.proto

package v1

import (
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	structpb "google.golang.org/protobuf/types/known/structpb"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ArchCreate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name        string                 `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Title       string                 `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Subtitle    string                 `protobuf:"bytes,3,opt,name=subtitle,proto3" json:"subtitle,omitempty"`
	Tags        []string               `protobuf:"bytes,4,rep,name=tags,proto3" json:"tags,omitempty"`
	Versions    []uint32               `protobuf:"varint,5,rep,packed,name=versions,proto3" json:"versions,omitempty"`
	OldVersions []uint64               `protobuf:"varint,6,rep,packed,name=old_versions,json=oldVersions,proto3" json:"old_versions,omitempty"`
	Release     *timestamppb.Timestamp `protobuf:"bytes,7,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp `protobuf:"bytes,8,opt,name=tested,proto3" json:"tested,omitempty"`
	Mark        string                 `protobuf:"bytes,9,opt,name=mark,proto3" json:"mark,omitempty"`
	Submarine   string                 `protobuf:"bytes,10,opt,name=submarine,proto3" json:"submarine,omitempty"`
	Numb        uint64                 `protobuf:"varint,11,opt,name=numb,proto3" json:"numb,omitempty"`
}

func (x *ArchCreate) Reset() {
	*x = ArchCreate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArchCreate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchCreate) ProtoMessage() {}

func (x *ArchCreate) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchCreate.ProtoReflect.Descriptor instead.
func (*ArchCreate) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{0}
}

func (x *ArchCreate) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ArchCreate) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *ArchCreate) GetSubtitle() string {
	if x != nil {
		return x.Subtitle
	}
	return ""
}

func (x *ArchCreate) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *ArchCreate) GetVersions() []uint32 {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *ArchCreate) GetOldVersions() []uint64 {
	if x != nil {
		return x.OldVersions
	}
	return nil
}

func (x *ArchCreate) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *ArchCreate) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

func (x *ArchCreate) GetMark() string {
	if x != nil {
		return x.Mark
	}
	return ""
}

func (x *ArchCreate) GetSubmarine() string {
	if x != nil {
		return x.Submarine
	}
	return ""
}

func (x *ArchCreate) GetNumb() uint64 {
	if x != nil {
		return x.Numb
	}
	return 0
}

type ArchGet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ArchGet) Reset() {
	*x = ArchGet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArchGet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchGet) ProtoMessage() {}

func (x *ArchGet) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchGet.ProtoReflect.Descriptor instead.
func (*ArchGet) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{1}
}

func (x *ArchGet) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ArchUpdate struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                  `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name        *wrapperspb.StringValue `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Title       *wrapperspb.StringValue `protobuf:"bytes,3,opt,name=title,proto3" json:"title,omitempty"`
	Subtitle    *wrapperspb.StringValue `protobuf:"bytes,4,opt,name=subtitle,proto3" json:"subtitle,omitempty"`
	Tags        *structpb.ListValue     `protobuf:"bytes,5,opt,name=tags,proto3" json:"tags,omitempty"`
	Versions    *structpb.ListValue     `protobuf:"bytes,6,opt,name=versions,proto3" json:"versions,omitempty"`
	OldVersions *structpb.ListValue     `protobuf:"bytes,7,opt,name=old_versions,json=oldVersions,proto3" json:"old_versions,omitempty"`
	Release     *timestamppb.Timestamp  `protobuf:"bytes,8,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp  `protobuf:"bytes,9,opt,name=tested,proto3" json:"tested,omitempty"`
	Mark        *wrapperspb.StringValue `protobuf:"bytes,10,opt,name=mark,proto3" json:"mark,omitempty"`
	Submarine   *wrapperspb.StringValue `protobuf:"bytes,11,opt,name=submarine,proto3" json:"submarine,omitempty"`
	Numb        *wrapperspb.UInt64Value `protobuf:"bytes,12,opt,name=numb,proto3" json:"numb,omitempty"`
}

func (x *ArchUpdate) Reset() {
	*x = ArchUpdate{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArchUpdate) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchUpdate) ProtoMessage() {}

func (x *ArchUpdate) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchUpdate.ProtoReflect.Descriptor instead.
func (*ArchUpdate) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{2}
}

func (x *ArchUpdate) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *ArchUpdate) GetName() *wrapperspb.StringValue {
	if x != nil {
		return x.Name
	}
	return nil
}

func (x *ArchUpdate) GetTitle() *wrapperspb.StringValue {
	if x != nil {
		return x.Title
	}
	return nil
}

func (x *ArchUpdate) GetSubtitle() *wrapperspb.StringValue {
	if x != nil {
		return x.Subtitle
	}
	return nil
}

func (x *ArchUpdate) GetTags() *structpb.ListValue {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *ArchUpdate) GetVersions() *structpb.ListValue {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *ArchUpdate) GetOldVersions() *structpb.ListValue {
	if x != nil {
		return x.OldVersions
	}
	return nil
}

func (x *ArchUpdate) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *ArchUpdate) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

func (x *ArchUpdate) GetMark() *wrapperspb.StringValue {
	if x != nil {
		return x.Mark
	}
	return nil
}

func (x *ArchUpdate) GetSubmarine() *wrapperspb.StringValue {
	if x != nil {
		return x.Submarine
	}
	return nil
}

func (x *ArchUpdate) GetNumb() *wrapperspb.UInt64Value {
	if x != nil {
		return x.Numb
	}
	return nil
}

type Arch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id          string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	UpdatedAt   *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty"`
	CreatedAt   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty"`
	Name        string                 `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	Title       string                 `protobuf:"bytes,5,opt,name=title,proto3" json:"title,omitempty"`
	Subtitle    string                 `protobuf:"bytes,6,opt,name=subtitle,proto3" json:"subtitle,omitempty"`
	Tags        []string               `protobuf:"bytes,7,rep,name=tags,proto3" json:"tags,omitempty"`
	Versions    []uint32               `protobuf:"varint,8,rep,packed,name=versions,proto3" json:"versions,omitempty"`
	OldVersions []uint64               `protobuf:"varint,9,rep,packed,name=old_versions,json=oldVersions,proto3" json:"old_versions,omitempty"`
	Release     *timestamppb.Timestamp `protobuf:"bytes,10,opt,name=release,proto3" json:"release,omitempty"`
	Tested      *timestamppb.Timestamp `protobuf:"bytes,11,opt,name=tested,proto3" json:"tested,omitempty"`
	Mark        string                 `protobuf:"bytes,12,opt,name=mark,proto3" json:"mark,omitempty"`
	Submarine   string                 `protobuf:"bytes,13,opt,name=submarine,proto3" json:"submarine,omitempty"`
	Numb        uint64                 `protobuf:"varint,14,opt,name=numb,proto3" json:"numb,omitempty"`
}

func (x *Arch) Reset() {
	*x = Arch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Arch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Arch) ProtoMessage() {}

func (x *Arch) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Arch.ProtoReflect.Descriptor instead.
func (*Arch) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{3}
}

func (x *Arch) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Arch) GetUpdatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdatedAt
	}
	return nil
}

func (x *Arch) GetCreatedAt() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedAt
	}
	return nil
}

func (x *Arch) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Arch) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *Arch) GetSubtitle() string {
	if x != nil {
		return x.Subtitle
	}
	return ""
}

func (x *Arch) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Arch) GetVersions() []uint32 {
	if x != nil {
		return x.Versions
	}
	return nil
}

func (x *Arch) GetOldVersions() []uint64 {
	if x != nil {
		return x.OldVersions
	}
	return nil
}

func (x *Arch) GetRelease() *timestamppb.Timestamp {
	if x != nil {
		return x.Release
	}
	return nil
}

func (x *Arch) GetTested() *timestamppb.Timestamp {
	if x != nil {
		return x.Tested
	}
	return nil
}

func (x *Arch) GetMark() string {
	if x != nil {
		return x.Mark
	}
	return ""
}

func (x *Arch) GetSubmarine() string {
	if x != nil {
		return x.Submarine
	}
	return ""
}

func (x *Arch) GetNumb() uint64 {
	if x != nil {
		return x.Numb
	}
	return 0
}

type ListArch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items []*Arch `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	Count uint64  `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
}

func (x *ListArch) Reset() {
	*x = ListArch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListArch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListArch) ProtoMessage() {}

func (x *ListArch) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListArch.ProtoReflect.Descriptor instead.
func (*ListArch) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{4}
}

func (x *ListArch) GetItems() []*Arch {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ListArch) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

type ArchDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ArchDelete) Reset() {
	*x = ArchDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArchDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchDelete) ProtoMessage() {}

func (x *ArchDelete) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchDelete.ProtoReflect.Descriptor instead.
func (*ArchDelete) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{5}
}

func (x *ArchDelete) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ArchFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PageNumber *wrapperspb.UInt64Value `protobuf:"bytes,1,opt,name=page_number,json=pageNumber,proto3" json:"page_number,omitempty"`
	PageSize   *wrapperspb.UInt64Value `protobuf:"bytes,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
	OrderBy    []string                `protobuf:"bytes,3,rep,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	Ids        []string                `protobuf:"bytes,4,rep,name=ids,proto3" json:"ids,omitempty"`
	Search     *wrapperspb.StringValue `protobuf:"bytes,5,opt,name=search,proto3" json:"search,omitempty"`
}

func (x *ArchFilter) Reset() {
	*x = ArchFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_examplepb_v1_arch_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArchFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArchFilter) ProtoMessage() {}

func (x *ArchFilter) ProtoReflect() protoreflect.Message {
	mi := &file_examplepb_v1_arch_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArchFilter.ProtoReflect.Descriptor instead.
func (*ArchFilter) Descriptor() ([]byte, []int) {
	return file_examplepb_v1_arch_proto_rawDescGZIP(), []int{6}
}

func (x *ArchFilter) GetPageNumber() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageNumber
	}
	return nil
}

func (x *ArchFilter) GetPageSize() *wrapperspb.UInt64Value {
	if x != nil {
		return x.PageSize
	}
	return nil
}

func (x *ArchFilter) GetOrderBy() []string {
	if x != nil {
		return x.OrderBy
	}
	return nil
}

func (x *ArchFilter) GetIds() []string {
	if x != nil {
		return x.Ids
	}
	return nil
}

func (x *ArchFilter) GetSearch() *wrapperspb.StringValue {
	if x != nil {
		return x.Search
	}
	return nil
}

var File_examplepb_v1_arch_proto protoreflect.FileDescriptor

var file_examplepb_v1_arch_proto_rawDesc = []byte{
	0x0a, 0x17, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x72, 0x63, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xd5, 0x02, 0x0a, 0x0a, 0x41, 0x72, 0x63, 0x68, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x73, 0x75,
	0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x75,
	0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x04,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x08, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x6c, 0x64, 0x5f, 0x76, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0b, 0x6f, 0x6c,
	0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x72, 0x65, 0x6c,
	0x65, 0x61, 0x73, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x12,
	0x32, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x65, 0x73,
	0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x6d, 0x61,
	0x72, 0x69, 0x6e, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6d,
	0x61, 0x72, 0x69, 0x6e, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x75, 0x6d, 0x62, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x04, 0x6e, 0x75, 0x6d, 0x62, 0x22, 0x19, 0x0a, 0x07, 0x41, 0x72, 0x63,
	0x68, 0x47, 0x65, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x02, 0x69, 0x64, 0x22, 0xed, 0x04, 0x0a, 0x0a, 0x41, 0x72, 0x63, 0x68, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x32, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x38, 0x0a, 0x08, 0x73, 0x75, 0x62,
	0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74,
	0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x73, 0x75, 0x62, 0x74, 0x69,
	0x74, 0x6c, 0x65, 0x12, 0x2e, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04, 0x74,
	0x61, 0x67, 0x73, 0x12, 0x36, 0x0a, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18,
	0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75,
	0x65, 0x52, 0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x3d, 0x0a, 0x0c, 0x6f,
	0x6c, 0x64, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0b, 0x6f,
	0x6c, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x34, 0x0a, 0x07, 0x72, 0x65,
	0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65,
	0x12, 0x32, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06, 0x74, 0x65,
	0x73, 0x74, 0x65, 0x64, 0x12, 0x30, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x0a, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x3a, 0x0a, 0x09, 0x73, 0x75, 0x62, 0x6d, 0x61, 0x72,
	0x69, 0x6e, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x09, 0x73, 0x75, 0x62, 0x6d, 0x61, 0x72, 0x69,
	0x6e, 0x65, 0x12, 0x30, 0x0a, 0x04, 0x6e, 0x75, 0x6d, 0x62, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x04,
	0x6e, 0x75, 0x6d, 0x62, 0x22, 0xd5, 0x03, 0x0a, 0x04, 0x41, 0x72, 0x63, 0x68, 0x12, 0x0e, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x39, 0x0a,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x75,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x64, 0x41, 0x74, 0x12, 0x39, 0x0a, 0x0a, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x1a, 0x0a,
	0x08, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x73, 0x75, 0x62, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x61, 0x67,
	0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x0d, 0x52,
	0x08, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x6f, 0x6c, 0x64,
	0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x04, 0x52,
	0x0b, 0x6f, 0x6c, 0x64, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x34, 0x0a, 0x07,
	0x72, 0x65, 0x6c, 0x65, 0x61, 0x73, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x07, 0x72, 0x65, 0x6c, 0x65, 0x61,
	0x73, 0x65, 0x12, 0x32, 0x0a, 0x06, 0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x06,
	0x74, 0x65, 0x73, 0x74, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x18, 0x0c,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6d, 0x61, 0x72, 0x6b, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x75,
	0x62, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x65, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73,
	0x75, 0x62, 0x6d, 0x61, 0x72, 0x69, 0x6e, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x75, 0x6d, 0x62,
	0x18, 0x0e, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x6e, 0x75, 0x6d, 0x62, 0x22, 0x4a, 0x0a, 0x08,
	0x4c, 0x69, 0x73, 0x74, 0x41, 0x72, 0x63, 0x68, 0x12, 0x28, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68, 0x52, 0x05, 0x69, 0x74, 0x65,
	0x6d, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x1c, 0x0a, 0x0a, 0x41, 0x72, 0x63, 0x68,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xe9, 0x01, 0x0a, 0x0a, 0x41, 0x72, 0x63, 0x68, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0b, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e,
	0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x4e, 0x75,
	0x6d, 0x62, 0x65, 0x72, 0x12, 0x39, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x55, 0x49, 0x6e, 0x74, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x69, 0x64,
	0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x03, 0x69, 0x64, 0x73, 0x12, 0x34, 0x0a, 0x06,
	0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53,
	0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x32, 0xb2, 0x03, 0x0a, 0x0b, 0x41, 0x72, 0x63, 0x68, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x51, 0x0a, 0x06, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x1a, 0x12, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93,
	0x02, 0x13, 0x3a, 0x01, 0x2a, 0x22, 0x0e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61,
	0x72, 0x63, 0x68, 0x65, 0x73, 0x12, 0x4d, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x15, 0x2e, 0x65,
	0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68,
	0x47, 0x65, 0x74, 0x1a, 0x12, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e,
	0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15, 0x12,
	0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x72, 0x63, 0x68, 0x65, 0x73, 0x2f,
	0x7b, 0x69, 0x64, 0x7d, 0x12, 0x56, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x18,
	0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72,
	0x63, 0x68, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x1a, 0x12, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68, 0x22, 0x1e, 0x82, 0xd3,
	0xe4, 0x93, 0x02, 0x18, 0x3a, 0x01, 0x2a, 0x32, 0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x72, 0x63, 0x68, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x57, 0x0a, 0x06,
	0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65,
	0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63, 0x68, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65,
	0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x1b, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x15,
	0x2a, 0x13, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x61, 0x72, 0x63, 0x68, 0x65, 0x73,
	0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x50, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x18, 0x2e,
	0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x72, 0x63,
	0x68, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x1a, 0x16, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c,
	0x65, 0x70, 0x62, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x41, 0x72, 0x63, 0x68, 0x22,
	0x16, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x10, 0x12, 0x0e, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31,
	0x2f, 0x61, 0x72, 0x63, 0x68, 0x65, 0x73, 0x42, 0x34, 0x5a, 0x32, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x69, 0x6b, 0x61, 0x6c, 0x61, 0x69, 0x2d, 0x6d, 0x69,
	0x74, 0x73, 0x69, 0x6e, 0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_examplepb_v1_arch_proto_rawDescOnce sync.Once
	file_examplepb_v1_arch_proto_rawDescData = file_examplepb_v1_arch_proto_rawDesc
)

func file_examplepb_v1_arch_proto_rawDescGZIP() []byte {
	file_examplepb_v1_arch_proto_rawDescOnce.Do(func() {
		file_examplepb_v1_arch_proto_rawDescData = protoimpl.X.CompressGZIP(file_examplepb_v1_arch_proto_rawDescData)
	})
	return file_examplepb_v1_arch_proto_rawDescData
}

var file_examplepb_v1_arch_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_examplepb_v1_arch_proto_goTypes = []interface{}{
	(*ArchCreate)(nil),             // 0: examplepb.v1.ArchCreate
	(*ArchGet)(nil),                // 1: examplepb.v1.ArchGet
	(*ArchUpdate)(nil),             // 2: examplepb.v1.ArchUpdate
	(*Arch)(nil),                   // 3: examplepb.v1.Arch
	(*ListArch)(nil),               // 4: examplepb.v1.ListArch
	(*ArchDelete)(nil),             // 5: examplepb.v1.ArchDelete
	(*ArchFilter)(nil),             // 6: examplepb.v1.ArchFilter
	(*timestamppb.Timestamp)(nil),  // 7: google.protobuf.Timestamp
	(*wrapperspb.StringValue)(nil), // 8: google.protobuf.StringValue
	(*structpb.ListValue)(nil),     // 9: google.protobuf.ListValue
	(*wrapperspb.UInt64Value)(nil), // 10: google.protobuf.UInt64Value
	(*emptypb.Empty)(nil),          // 11: google.protobuf.Empty
}
var file_examplepb_v1_arch_proto_depIdxs = []int32{
	7,  // 0: examplepb.v1.ArchCreate.release:type_name -> google.protobuf.Timestamp
	7,  // 1: examplepb.v1.ArchCreate.tested:type_name -> google.protobuf.Timestamp
	8,  // 2: examplepb.v1.ArchUpdate.name:type_name -> google.protobuf.StringValue
	8,  // 3: examplepb.v1.ArchUpdate.title:type_name -> google.protobuf.StringValue
	8,  // 4: examplepb.v1.ArchUpdate.subtitle:type_name -> google.protobuf.StringValue
	9,  // 5: examplepb.v1.ArchUpdate.tags:type_name -> google.protobuf.ListValue
	9,  // 6: examplepb.v1.ArchUpdate.versions:type_name -> google.protobuf.ListValue
	9,  // 7: examplepb.v1.ArchUpdate.old_versions:type_name -> google.protobuf.ListValue
	7,  // 8: examplepb.v1.ArchUpdate.release:type_name -> google.protobuf.Timestamp
	7,  // 9: examplepb.v1.ArchUpdate.tested:type_name -> google.protobuf.Timestamp
	8,  // 10: examplepb.v1.ArchUpdate.mark:type_name -> google.protobuf.StringValue
	8,  // 11: examplepb.v1.ArchUpdate.submarine:type_name -> google.protobuf.StringValue
	10, // 12: examplepb.v1.ArchUpdate.numb:type_name -> google.protobuf.UInt64Value
	7,  // 13: examplepb.v1.Arch.updated_at:type_name -> google.protobuf.Timestamp
	7,  // 14: examplepb.v1.Arch.created_at:type_name -> google.protobuf.Timestamp
	7,  // 15: examplepb.v1.Arch.release:type_name -> google.protobuf.Timestamp
	7,  // 16: examplepb.v1.Arch.tested:type_name -> google.protobuf.Timestamp
	3,  // 17: examplepb.v1.ListArch.items:type_name -> examplepb.v1.Arch
	10, // 18: examplepb.v1.ArchFilter.page_number:type_name -> google.protobuf.UInt64Value
	10, // 19: examplepb.v1.ArchFilter.page_size:type_name -> google.protobuf.UInt64Value
	8,  // 20: examplepb.v1.ArchFilter.search:type_name -> google.protobuf.StringValue
	0,  // 21: examplepb.v1.ArchService.Create:input_type -> examplepb.v1.ArchCreate
	1,  // 22: examplepb.v1.ArchService.Get:input_type -> examplepb.v1.ArchGet
	2,  // 23: examplepb.v1.ArchService.Update:input_type -> examplepb.v1.ArchUpdate
	5,  // 24: examplepb.v1.ArchService.Delete:input_type -> examplepb.v1.ArchDelete
	6,  // 25: examplepb.v1.ArchService.List:input_type -> examplepb.v1.ArchFilter
	3,  // 26: examplepb.v1.ArchService.Create:output_type -> examplepb.v1.Arch
	3,  // 27: examplepb.v1.ArchService.Get:output_type -> examplepb.v1.Arch
	3,  // 28: examplepb.v1.ArchService.Update:output_type -> examplepb.v1.Arch
	11, // 29: examplepb.v1.ArchService.Delete:output_type -> google.protobuf.Empty
	4,  // 30: examplepb.v1.ArchService.List:output_type -> examplepb.v1.ListArch
	26, // [26:31] is the sub-list for method output_type
	21, // [21:26] is the sub-list for method input_type
	21, // [21:21] is the sub-list for extension type_name
	21, // [21:21] is the sub-list for extension extendee
	0,  // [0:21] is the sub-list for field type_name
}

func init() { file_examplepb_v1_arch_proto_init() }
func file_examplepb_v1_arch_proto_init() {
	if File_examplepb_v1_arch_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_examplepb_v1_arch_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArchCreate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArchGet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArchUpdate); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Arch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListArch); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArchDelete); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_examplepb_v1_arch_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArchFilter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_examplepb_v1_arch_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_examplepb_v1_arch_proto_goTypes,
		DependencyIndexes: file_examplepb_v1_arch_proto_depIdxs,
		MessageInfos:      file_examplepb_v1_arch_proto_msgTypes,
	}.Build()
	File_examplepb_v1_arch_proto = out.File
	file_examplepb_v1_arch_proto_rawDesc = nil
	file_examplepb_v1_arch_proto_goTypes = nil
	file_examplepb_v1_arch_proto_depIdxs = nil
}
