// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.15.8
// source: entity.proto

package collectionxservice

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type OrderTypeProto int32

const (
	OrderTypeProto_ORDER_TYPE_NONE OrderTypeProto = 0
	OrderTypeProto_ORDER_TYPE_ASC  OrderTypeProto = 1
	OrderTypeProto_ORDER_TYPE_DESC OrderTypeProto = 2
)

// Enum value maps for OrderTypeProto.
var (
	OrderTypeProto_name = map[int32]string{
		0: "ORDER_TYPE_NONE",
		1: "ORDER_TYPE_ASC",
		2: "ORDER_TYPE_DESC",
	}
	OrderTypeProto_value = map[string]int32{
		"ORDER_TYPE_NONE": 0,
		"ORDER_TYPE_ASC":  1,
		"ORDER_TYPE_DESC": 2,
	}
)

func (x OrderTypeProto) Enum() *OrderTypeProto {
	p := new(OrderTypeProto)
	*p = x
	return p
}

func (x OrderTypeProto) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (OrderTypeProto) Descriptor() protoreflect.EnumDescriptor {
	return file_entity_proto_enumTypes[0].Descriptor()
}

func (OrderTypeProto) Type() protoreflect.EnumType {
	return &file_entity_proto_enumTypes[0]
}

func (x OrderTypeProto) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use OrderTypeProto.Descriptor instead.
func (OrderTypeProto) EnumDescriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{0}
}

type DocumentChangeKind int32

const (
	DocumentChangeKind_DOCUMENT_KIND_ADDED     DocumentChangeKind = 0
	DocumentChangeKind_DOCUMENT_KIND_REMOVED   DocumentChangeKind = 1
	DocumentChangeKind_DOCUMENT_KIND_MODIFIED  DocumentChangeKind = 2
	DocumentChangeKind_DOCUMENT_KIND_SNAPSHOTS DocumentChangeKind = 3
)

// Enum value maps for DocumentChangeKind.
var (
	DocumentChangeKind_name = map[int32]string{
		0: "DOCUMENT_KIND_ADDED",
		1: "DOCUMENT_KIND_REMOVED",
		2: "DOCUMENT_KIND_MODIFIED",
		3: "DOCUMENT_KIND_SNAPSHOTS",
	}
	DocumentChangeKind_value = map[string]int32{
		"DOCUMENT_KIND_ADDED":     0,
		"DOCUMENT_KIND_REMOVED":   1,
		"DOCUMENT_KIND_MODIFIED":  2,
		"DOCUMENT_KIND_SNAPSHOTS": 3,
	}
)

func (x DocumentChangeKind) Enum() *DocumentChangeKind {
	p := new(DocumentChangeKind)
	*p = x
	return p
}

func (x DocumentChangeKind) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DocumentChangeKind) Descriptor() protoreflect.EnumDescriptor {
	return file_entity_proto_enumTypes[1].Descriptor()
}

func (DocumentChangeKind) Type() protoreflect.EnumType {
	return &file_entity_proto_enumTypes[1]
}

func (x DocumentChangeKind) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DocumentChangeKind.Descriptor instead.
func (DocumentChangeKind) EnumDescriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{1}
}

type StandardAPIProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string     `protobuf:"bytes,4,opt,name=message,proto3" json:"message,omitempty"`
	Meta    *MetaProto `protobuf:"bytes,5,opt,name=meta,proto3,oneof" json:"meta,omitempty"`
}

func (x *StandardAPIProto) Reset() {
	*x = StandardAPIProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StandardAPIProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StandardAPIProto) ProtoMessage() {}

func (x *StandardAPIProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StandardAPIProto.ProtoReflect.Descriptor instead.
func (*StandardAPIProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{0}
}

func (x *StandardAPIProto) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *StandardAPIProto) GetMeta() *MetaProto {
	if x != nil {
		return x.Meta
	}
	return nil
}

type DateRangeProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field string                 `protobuf:"bytes,1,opt,name=field,proto3" json:"field,omitempty"`
	Start *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
	End   *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=end,proto3" json:"end,omitempty"`
}

func (x *DateRangeProto) Reset() {
	*x = DateRangeProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DateRangeProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DateRangeProto) ProtoMessage() {}

func (x *DateRangeProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DateRangeProto.ProtoReflect.Descriptor instead.
func (*DateRangeProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{1}
}

func (x *DateRangeProto) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *DateRangeProto) GetStart() *timestamppb.Timestamp {
	if x != nil {
		return x.Start
	}
	return nil
}

func (x *DateRangeProto) GetEnd() *timestamppb.Timestamp {
	if x != nil {
		return x.End
	}
	return nil
}

type FilterProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	By  string `protobuf:"bytes,1,opt,name=by,proto3" json:"by,omitempty"`
	Op  string `protobuf:"bytes,2,opt,name=op,proto3" json:"op,omitempty"`
	Val []byte `protobuf:"bytes,3,opt,name=val,proto3" json:"val,omitempty"`
}

func (x *FilterProto) Reset() {
	*x = FilterProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilterProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilterProto) ProtoMessage() {}

func (x *FilterProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilterProto.ProtoReflect.Descriptor instead.
func (*FilterProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{2}
}

func (x *FilterProto) GetBy() string {
	if x != nil {
		return x.By
	}
	return ""
}

func (x *FilterProto) GetOp() string {
	if x != nil {
		return x.Op
	}
	return ""
}

func (x *FilterProto) GetVal() []byte {
	if x != nil {
		return x.Val
	}
	return nil
}

type SortProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	OrderBy   string         `protobuf:"bytes,1,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	OrderType OrderTypeProto `protobuf:"varint,2,opt,name=order_type,json=orderType,proto3,enum=collectionx.OrderTypeProto" json:"order_type,omitempty"`
}

func (x *SortProto) Reset() {
	*x = SortProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SortProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SortProto) ProtoMessage() {}

func (x *SortProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SortProto.ProtoReflect.Descriptor instead.
func (*SortProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{3}
}

func (x *SortProto) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *SortProto) GetOrderType() OrderTypeProto {
	if x != nil {
		return x.OrderType
	}
	return OrderTypeProto_ORDER_TYPE_NONE
}

type FilteringProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sort      []*SortProto    `protobuf:"bytes,1,rep,name=sort,proto3" json:"sort,omitempty"`
	Filter    []*FilterProto  `protobuf:"bytes,2,rep,name=filter,proto3" json:"filter,omitempty"`
	DateRange *DateRangeProto `protobuf:"bytes,3,opt,name=date_range,json=dateRange,proto3,oneof" json:"date_range,omitempty"`
}

func (x *FilteringProto) Reset() {
	*x = FilteringProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FilteringProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FilteringProto) ProtoMessage() {}

func (x *FilteringProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FilteringProto.ProtoReflect.Descriptor instead.
func (*FilteringProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{4}
}

func (x *FilteringProto) GetSort() []*SortProto {
	if x != nil {
		return x.Sort
	}
	return nil
}

func (x *FilteringProto) GetFilter() []*FilterProto {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *FilteringProto) GetDateRange() *DateRangeProto {
	if x != nil {
		return x.DateRange
	}
	return nil
}

type MetaProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page    int32 `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	PerPage int32 `protobuf:"varint,2,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
	Total   int32 `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *MetaProto) Reset() {
	*x = MetaProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MetaProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MetaProto) ProtoMessage() {}

func (x *MetaProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MetaProto.ProtoReflect.Descriptor instead.
func (*MetaProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{5}
}

func (x *MetaProto) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *MetaProto) GetPerPage() int32 {
	if x != nil {
		return x.PerPage
	}
	return 0
}

func (x *MetaProto) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type ParameterRequestProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PerPage   int32           `protobuf:"varint,1,opt,name=per_page,json=perPage,proto3" json:"per_page,omitempty"`
	Page      int32           `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	OrderBy   string          `protobuf:"bytes,3,opt,name=order_by,json=orderBy,proto3" json:"order_by,omitempty"`
	OrderType OrderTypeProto  `protobuf:"varint,4,opt,name=order_type,json=orderType,proto3,enum=collectionx.OrderTypeProto" json:"order_type,omitempty"`
	Search    string          `protobuf:"bytes,5,opt,name=search,proto3" json:"search,omitempty"`
	Keyword   string          `protobuf:"bytes,6,opt,name=keyword,proto3" json:"keyword,omitempty"`
	SearchBy  string          `protobuf:"bytes,7,opt,name=search_by,json=searchBy,proto3" json:"search_by,omitempty"`
	DateRange *DateRangeProto `protobuf:"bytes,8,opt,name=date_range,json=dateRange,proto3,oneof" json:"date_range,omitempty"`
}

func (x *ParameterRequestProto) Reset() {
	*x = ParameterRequestProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ParameterRequestProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ParameterRequestProto) ProtoMessage() {}

func (x *ParameterRequestProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ParameterRequestProto.ProtoReflect.Descriptor instead.
func (*ParameterRequestProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{6}
}

func (x *ParameterRequestProto) GetPerPage() int32 {
	if x != nil {
		return x.PerPage
	}
	return 0
}

func (x *ParameterRequestProto) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *ParameterRequestProto) GetOrderBy() string {
	if x != nil {
		return x.OrderBy
	}
	return ""
}

func (x *ParameterRequestProto) GetOrderType() OrderTypeProto {
	if x != nil {
		return x.OrderType
	}
	return OrderTypeProto_ORDER_TYPE_NONE
}

func (x *ParameterRequestProto) GetSearch() string {
	if x != nil {
		return x.Search
	}
	return ""
}

func (x *ParameterRequestProto) GetKeyword() string {
	if x != nil {
		return x.Keyword
	}
	return ""
}

func (x *ParameterRequestProto) GetSearchBy() string {
	if x != nil {
		return x.SearchBy
	}
	return ""
}

func (x *ParameterRequestProto) GetDateRange() *DateRangeProto {
	if x != nil {
		return x.DateRange
	}
	return nil
}

type TimestampProto struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	CreatedTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=created_time,json=createdTime,proto3" json:"created_time,omitempty"`
	ReadTime    *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=read_time,json=readTime,proto3" json:"read_time,omitempty"`
	UpdateTime  *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *TimestampProto) Reset() {
	*x = TimestampProto{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimestampProto) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimestampProto) ProtoMessage() {}

func (x *TimestampProto) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimestampProto.ProtoReflect.Descriptor instead.
func (*TimestampProto) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{7}
}

func (x *TimestampProto) GetCreatedTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreatedTime
	}
	return nil
}

func (x *TimestampProto) GetReadTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ReadTime
	}
	return nil
}

func (x *TimestampProto) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

type DocumentChange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Kind      DocumentChangeKind `protobuf:"varint,1,opt,name=kind,proto3,enum=collectionx.DocumentChangeKind" json:"kind,omitempty"`
	Data      []byte             `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Timestamp *TimestampProto    `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *DocumentChange) Reset() {
	*x = DocumentChange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_entity_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentChange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentChange) ProtoMessage() {}

func (x *DocumentChange) ProtoReflect() protoreflect.Message {
	mi := &file_entity_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentChange.ProtoReflect.Descriptor instead.
func (*DocumentChange) Descriptor() ([]byte, []int) {
	return file_entity_proto_rawDescGZIP(), []int{8}
}

func (x *DocumentChange) GetKind() DocumentChangeKind {
	if x != nil {
		return x.Kind
	}
	return DocumentChangeKind_DOCUMENT_KIND_ADDED
}

func (x *DocumentChange) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *DocumentChange) GetTimestamp() *TimestampProto {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_entity_proto protoreflect.FileDescriptor

var file_entity_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0b,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x66, 0x0a, 0x10,
	0x53, 0x74, 0x61, 0x6e, 0x64, 0x61, 0x72, 0x64, 0x41, 0x50, 0x49, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x2f, 0x0a, 0x04, 0x6d, 0x65,
	0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x48, 0x00, 0x52, 0x04, 0x6d, 0x65, 0x74, 0x61, 0x88, 0x01, 0x01, 0x42, 0x07, 0x0a, 0x05, 0x5f,
	0x6d, 0x65, 0x74, 0x61, 0x22, 0x86, 0x01, 0x0a, 0x0e, 0x44, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x30, 0x0a,
	0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12,
	0x2c, 0x0a, 0x03, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x03, 0x65, 0x6e, 0x64, 0x22, 0x3f, 0x0a,
	0x0b, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x0a, 0x02,
	0x62, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x62, 0x79, 0x12, 0x0e, 0x0a, 0x02,
	0x6f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x6f, 0x70, 0x12, 0x10, 0x0a, 0x03,
	0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x76, 0x61, 0x6c, 0x22, 0x62,
	0x0a, 0x09, 0x53, 0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x0a, 0x08, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6f,
	0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x3a, 0x0a, 0x0a, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79,
	0x70, 0x65, 0x22, 0xbe, 0x01, 0x0a, 0x0e, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x69, 0x6e, 0x67,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x2a, 0x0a, 0x04, 0x73, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x78, 0x2e, 0x53, 0x6f, 0x72, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x04, 0x73, 0x6f, 0x72,
	0x74, 0x12, 0x30, 0x0a, 0x06, 0x66, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x18, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e,
	0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x06, 0x66, 0x69, 0x6c,
	0x74, 0x65, 0x72, 0x12, 0x3f, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x61, 0x6e, 0x67,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x48, 0x00, 0x52, 0x09, 0x64, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67,
	0x65, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x61,
	0x6e, 0x67, 0x65, 0x22, 0x50, 0x0a, 0x09, 0x4d, 0x65, 0x74, 0x61, 0x50, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x61, 0x67, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x67, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x70, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x12,
	0x14, 0x0a, 0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0xbc, 0x02, 0x0a, 0x15, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x65,
	0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x19, 0x0a, 0x08, 0x70, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x07, 0x70, 0x65, 0x72, 0x50, 0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x19,
	0x0a, 0x08, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x5f, 0x62, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x6f, 0x72, 0x64, 0x65, 0x72, 0x42, 0x79, 0x12, 0x3a, 0x0a, 0x0a, 0x6f, 0x72, 0x64,
	0x65, 0x72, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1b, 0x2e,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x4f, 0x72, 0x64, 0x65,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x09, 0x6f, 0x72, 0x64, 0x65,
	0x72, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x65, 0x61, 0x72, 0x63, 0x68, 0x12, 0x18, 0x0a,
	0x07, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x73, 0x65, 0x61, 0x72, 0x63,
	0x68, 0x5f, 0x62, 0x79, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x42, 0x79, 0x12, 0x3f, 0x0a, 0x0a, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x72, 0x61, 0x6e,
	0x67, 0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f, 0x6c, 0x6c, 0x65,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x44, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x48, 0x00, 0x52, 0x09, 0x64, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x88, 0x01, 0x01, 0x42, 0x0d, 0x0a, 0x0b, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x72,
	0x61, 0x6e, 0x67, 0x65, 0x22, 0xc5, 0x01, 0x0a, 0x0e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x3d, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x72, 0x65, 0x61, 0x64, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x08, 0x72, 0x65, 0x61, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x94, 0x01, 0x0a,
	0x0e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x12,
	0x33, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1f, 0x2e,
	0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4b, 0x69, 0x6e, 0x64, 0x52, 0x04,
	0x6b, 0x69, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x12, 0x39, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x63, 0x6f,
	0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74,
	0x61, 0x6d, 0x70, 0x2a, 0x4e, 0x0a, 0x0e, 0x4f, 0x72, 0x64, 0x65, 0x72, 0x54, 0x79, 0x70, 0x65,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x13, 0x0a, 0x0f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54,
	0x59, 0x50, 0x45, 0x5f, 0x4e, 0x4f, 0x4e, 0x45, 0x10, 0x00, 0x12, 0x12, 0x0a, 0x0e, 0x4f, 0x52,
	0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x41, 0x53, 0x43, 0x10, 0x01, 0x12, 0x13,
	0x0a, 0x0f, 0x4f, 0x52, 0x44, 0x45, 0x52, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x44, 0x45, 0x53,
	0x43, 0x10, 0x02, 0x2a, 0x81, 0x01, 0x0a, 0x12, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x4b, 0x69, 0x6e, 0x64, 0x12, 0x17, 0x0a, 0x13, 0x44, 0x4f,
	0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x41, 0x44, 0x44, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x19, 0x0a, 0x15, 0x44, 0x4f, 0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f,
	0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x52, 0x45, 0x4d, 0x4f, 0x56, 0x45, 0x44, 0x10, 0x01, 0x12, 0x1a,
	0x0a, 0x16, 0x44, 0x4f, 0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f,
	0x4d, 0x4f, 0x44, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x02, 0x12, 0x1b, 0x0a, 0x17, 0x44, 0x4f,
	0x43, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4b, 0x49, 0x4e, 0x44, 0x5f, 0x53, 0x4e, 0x41, 0x50,
	0x53, 0x48, 0x4f, 0x54, 0x53, 0x10, 0x03, 0x42, 0x26, 0x5a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x70, 0x72, 0x69, 0x76, 0x79, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6c,
	0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x78, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_entity_proto_rawDescOnce sync.Once
	file_entity_proto_rawDescData = file_entity_proto_rawDesc
)

func file_entity_proto_rawDescGZIP() []byte {
	file_entity_proto_rawDescOnce.Do(func() {
		file_entity_proto_rawDescData = protoimpl.X.CompressGZIP(file_entity_proto_rawDescData)
	})
	return file_entity_proto_rawDescData
}

var file_entity_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_entity_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_entity_proto_goTypes = []interface{}{
	(OrderTypeProto)(0),           // 0: collectionx.OrderTypeProto
	(DocumentChangeKind)(0),       // 1: collectionx.DocumentChangeKind
	(*StandardAPIProto)(nil),      // 2: collectionx.StandardAPIProto
	(*DateRangeProto)(nil),        // 3: collectionx.DateRangeProto
	(*FilterProto)(nil),           // 4: collectionx.FilterProto
	(*SortProto)(nil),             // 5: collectionx.SortProto
	(*FilteringProto)(nil),        // 6: collectionx.FilteringProto
	(*MetaProto)(nil),             // 7: collectionx.MetaProto
	(*ParameterRequestProto)(nil), // 8: collectionx.ParameterRequestProto
	(*TimestampProto)(nil),        // 9: collectionx.TimestampProto
	(*DocumentChange)(nil),        // 10: collectionx.DocumentChange
	(*timestamppb.Timestamp)(nil), // 11: google.protobuf.Timestamp
}
var file_entity_proto_depIdxs = []int32{
	7,  // 0: collectionx.StandardAPIProto.meta:type_name -> collectionx.MetaProto
	11, // 1: collectionx.DateRangeProto.start:type_name -> google.protobuf.Timestamp
	11, // 2: collectionx.DateRangeProto.end:type_name -> google.protobuf.Timestamp
	0,  // 3: collectionx.SortProto.order_type:type_name -> collectionx.OrderTypeProto
	5,  // 4: collectionx.FilteringProto.sort:type_name -> collectionx.SortProto
	4,  // 5: collectionx.FilteringProto.filter:type_name -> collectionx.FilterProto
	3,  // 6: collectionx.FilteringProto.date_range:type_name -> collectionx.DateRangeProto
	0,  // 7: collectionx.ParameterRequestProto.order_type:type_name -> collectionx.OrderTypeProto
	3,  // 8: collectionx.ParameterRequestProto.date_range:type_name -> collectionx.DateRangeProto
	11, // 9: collectionx.TimestampProto.created_time:type_name -> google.protobuf.Timestamp
	11, // 10: collectionx.TimestampProto.read_time:type_name -> google.protobuf.Timestamp
	11, // 11: collectionx.TimestampProto.update_time:type_name -> google.protobuf.Timestamp
	1,  // 12: collectionx.DocumentChange.kind:type_name -> collectionx.DocumentChangeKind
	9,  // 13: collectionx.DocumentChange.timestamp:type_name -> collectionx.TimestampProto
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_entity_proto_init() }
func file_entity_proto_init() {
	if File_entity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_entity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*StandardAPIProto); i {
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
		file_entity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DateRangeProto); i {
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
		file_entity_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilterProto); i {
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
		file_entity_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SortProto); i {
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
		file_entity_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FilteringProto); i {
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
		file_entity_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MetaProto); i {
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
		file_entity_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ParameterRequestProto); i {
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
		file_entity_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimestampProto); i {
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
		file_entity_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocumentChange); i {
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
	file_entity_proto_msgTypes[0].OneofWrappers = []interface{}{}
	file_entity_proto_msgTypes[4].OneofWrappers = []interface{}{}
	file_entity_proto_msgTypes[6].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_entity_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_entity_proto_goTypes,
		DependencyIndexes: file_entity_proto_depIdxs,
		EnumInfos:         file_entity_proto_enumTypes,
		MessageInfos:      file_entity_proto_msgTypes,
	}.Build()
	File_entity_proto = out.File
	file_entity_proto_rawDesc = nil
	file_entity_proto_goTypes = nil
	file_entity_proto_depIdxs = nil
}
