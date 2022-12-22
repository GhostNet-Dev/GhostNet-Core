// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.6.1
// source: masternet_packet.proto

package packets

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MasterNodeUserInfoSq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	User   *GhostUser    `protobuf:"bytes,2,opt,name=User,proto3" json:"User,omitempty"`
}

func (x *MasterNodeUserInfoSq) Reset() {
	*x = MasterNodeUserInfoSq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MasterNodeUserInfoSq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasterNodeUserInfoSq) ProtoMessage() {}

func (x *MasterNodeUserInfoSq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasterNodeUserInfoSq.ProtoReflect.Descriptor instead.
func (*MasterNodeUserInfoSq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{0}
}

func (x *MasterNodeUserInfoSq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *MasterNodeUserInfoSq) GetUser() *GhostUser {
	if x != nil {
		return x.User
	}
	return nil
}

type MasterNodeUserInfoCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	User   *GhostUser    `protobuf:"bytes,2,opt,name=User,proto3" json:"User,omitempty"`
}

func (x *MasterNodeUserInfoCq) Reset() {
	*x = MasterNodeUserInfoCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MasterNodeUserInfoCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasterNodeUserInfoCq) ProtoMessage() {}

func (x *MasterNodeUserInfoCq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasterNodeUserInfoCq.ProtoReflect.Descriptor instead.
func (*MasterNodeUserInfoCq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{1}
}

func (x *MasterNodeUserInfoCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *MasterNodeUserInfoCq) GetUser() *GhostUser {
	if x != nil {
		return x.User
	}
	return nil
}

type ReqeustMasterNodeListSq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master     *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	StartIndex uint32        `protobuf:"varint,2,opt,name=StartIndex,proto3" json:"StartIndex,omitempty"`
}

func (x *ReqeustMasterNodeListSq) Reset() {
	*x = ReqeustMasterNodeListSq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReqeustMasterNodeListSq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReqeustMasterNodeListSq) ProtoMessage() {}

func (x *ReqeustMasterNodeListSq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReqeustMasterNodeListSq.ProtoReflect.Descriptor instead.
func (*ReqeustMasterNodeListSq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{2}
}

func (x *ReqeustMasterNodeListSq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *ReqeustMasterNodeListSq) GetStartIndex() uint32 {
	if x != nil {
		return x.StartIndex
	}
	return 0
}

type RequestMasterNodeListCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
}

func (x *RequestMasterNodeListCq) Reset() {
	*x = RequestMasterNodeListCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestMasterNodeListCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestMasterNodeListCq) ProtoMessage() {}

func (x *RequestMasterNodeListCq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestMasterNodeListCq.ProtoReflect.Descriptor instead.
func (*RequestMasterNodeListCq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{3}
}

func (x *RequestMasterNodeListCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

type ResponseMasterNodeListSq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	User   []*GhostUser  `protobuf:"bytes,2,rep,name=User,proto3" json:"User,omitempty"`
	Index  uint32        `protobuf:"varint,3,opt,name=Index,proto3" json:"Index,omitempty"`
	Total  uint32        `protobuf:"varint,4,opt,name=Total,proto3" json:"Total,omitempty"`
}

func (x *ResponseMasterNodeListSq) Reset() {
	*x = ResponseMasterNodeListSq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseMasterNodeListSq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseMasterNodeListSq) ProtoMessage() {}

func (x *ResponseMasterNodeListSq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseMasterNodeListSq.ProtoReflect.Descriptor instead.
func (*ResponseMasterNodeListSq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{4}
}

func (x *ResponseMasterNodeListSq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *ResponseMasterNodeListSq) GetUser() []*GhostUser {
	if x != nil {
		return x.User
	}
	return nil
}

func (x *ResponseMasterNodeListSq) GetIndex() uint32 {
	if x != nil {
		return x.Index
	}
	return 0
}

func (x *ResponseMasterNodeListSq) GetTotal() uint32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type ResponseMasterNodeListCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
}

func (x *ResponseMasterNodeListCq) Reset() {
	*x = ResponseMasterNodeListCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseMasterNodeListCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseMasterNodeListCq) ProtoMessage() {}

func (x *ResponseMasterNodeListCq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseMasterNodeListCq.ProtoReflect.Descriptor instead.
func (*ResponseMasterNodeListCq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{5}
}

func (x *ResponseMasterNodeListCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

type SearchGhostPubKeySq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master   *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	Nickname string        `protobuf:"bytes,2,opt,name=Nickname,proto3" json:"Nickname,omitempty"`
}

func (x *SearchGhostPubKeySq) Reset() {
	*x = SearchGhostPubKeySq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchGhostPubKeySq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchGhostPubKeySq) ProtoMessage() {}

func (x *SearchGhostPubKeySq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchGhostPubKeySq.ProtoReflect.Descriptor instead.
func (*SearchGhostPubKeySq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{6}
}

func (x *SearchGhostPubKeySq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *SearchGhostPubKeySq) GetNickname() string {
	if x != nil {
		return x.Nickname
	}
	return ""
}

type SearchGhostPubKeyCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	User   []*GhostUser  `protobuf:"bytes,2,rep,name=User,proto3" json:"User,omitempty"`
}

func (x *SearchGhostPubKeyCq) Reset() {
	*x = SearchGhostPubKeyCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_masternet_packet_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SearchGhostPubKeyCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SearchGhostPubKeyCq) ProtoMessage() {}

func (x *SearchGhostPubKeyCq) ProtoReflect() protoreflect.Message {
	mi := &file_masternet_packet_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SearchGhostPubKeyCq.ProtoReflect.Descriptor instead.
func (*SearchGhostPubKeyCq) Descriptor() ([]byte, []int) {
	return file_masternet_packet_proto_rawDescGZIP(), []int{7}
}

func (x *SearchGhostPubKeyCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *SearchGhostPubKeyCq) GetUser() []*GhostUser {
	if x != nil {
		return x.User
	}
	return nil
}

var File_masternet_packet_proto protoreflect.FileDescriptor

var file_masternet_packet_proto_rawDesc = []byte{
	0x0a, 0x16, 0x6d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x5f, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e,
	0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x7f, 0x0a, 0x14, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x73, 0x65,
	0x72, 0x49, 0x6e, 0x66, 0x6f, 0x53, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e,
	0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65,
	0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12,
	0x2f, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b, 0x2e,
	0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73,
	0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72,
	0x22, 0x7f, 0x0a, 0x14, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x55, 0x73,
	0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x43, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x12, 0x2f, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1b,
	0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x55, 0x73, 0x65,
	0x72, 0x22, 0x71, 0x0a, 0x17, 0x52, 0x65, 0x71, 0x65, 0x75, 0x73, 0x74, 0x4d, 0x61, 0x73, 0x74,
	0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x53, 0x71, 0x12, 0x36, 0x0a, 0x06,
	0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e,
	0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61,
	0x73, 0x74, 0x65, 0x72, 0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x72, 0x74, 0x49, 0x6e, 0x64,
	0x65, 0x78, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x72, 0x74, 0x49,
	0x6e, 0x64, 0x65, 0x78, 0x22, 0x51, 0x0a, 0x17, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d,
	0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x43, 0x71, 0x12,
	0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52,
	0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x22, 0xaf, 0x01, 0x0a, 0x18, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69,
	0x73, 0x74, 0x53, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04,
	0x55, 0x73, 0x65, 0x72, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x68, 0x6f,
	0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x47, 0x68,
	0x6f, 0x73, 0x74, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72, 0x12, 0x14, 0x0a,
	0x05, 0x49, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x49, 0x6e,
	0x64, 0x65, 0x78, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x05, 0x54, 0x6f, 0x74, 0x61, 0x6c, 0x22, 0x52, 0x0a, 0x18, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x4c,
	0x69, 0x73, 0x74, 0x43, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74,
	0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x22, 0x69, 0x0a,
	0x13, 0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x75, 0x62, 0x4b,
	0x65, 0x79, 0x53, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x1a, 0x0a, 0x08,
	0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x22, 0x7e, 0x0a, 0x13, 0x53, 0x65, 0x61, 0x72,
	0x63, 0x68, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x43, 0x71, 0x12,
	0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52,
	0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x2f, 0x0a, 0x04, 0x55, 0x73, 0x65, 0x72, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1b, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74,
	0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x55, 0x73,
	0x65, 0x72, 0x52, 0x04, 0x55, 0x73, 0x65, 0x72, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d,
	0x44, 0x65, 0x76, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d, 0x43, 0x6f, 0x72,
	0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_masternet_packet_proto_rawDescOnce sync.Once
	file_masternet_packet_proto_rawDescData = file_masternet_packet_proto_rawDesc
)

func file_masternet_packet_proto_rawDescGZIP() []byte {
	file_masternet_packet_proto_rawDescOnce.Do(func() {
		file_masternet_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_masternet_packet_proto_rawDescData)
	})
	return file_masternet_packet_proto_rawDescData
}

var file_masternet_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_masternet_packet_proto_goTypes = []interface{}{
	(*MasterNodeUserInfoSq)(nil),     // 0: ghostnet.packets.MasterNodeUserInfoSq
	(*MasterNodeUserInfoCq)(nil),     // 1: ghostnet.packets.MasterNodeUserInfoCq
	(*ReqeustMasterNodeListSq)(nil),  // 2: ghostnet.packets.ReqeustMasterNodeListSq
	(*RequestMasterNodeListCq)(nil),  // 3: ghostnet.packets.RequestMasterNodeListCq
	(*ResponseMasterNodeListSq)(nil), // 4: ghostnet.packets.ResponseMasterNodeListSq
	(*ResponseMasterNodeListCq)(nil), // 5: ghostnet.packets.ResponseMasterNodeListCq
	(*SearchGhostPubKeySq)(nil),      // 6: ghostnet.packets.SearchGhostPubKeySq
	(*SearchGhostPubKeyCq)(nil),      // 7: ghostnet.packets.SearchGhostPubKeyCq
	(*MasterPacket)(nil),             // 8: ghostnet.packets.MasterPacket
	(*GhostUser)(nil),                // 9: ghostnet.packets.GhostUser
}
var file_masternet_packet_proto_depIdxs = []int32{
	8,  // 0: ghostnet.packets.MasterNodeUserInfoSq.Master:type_name -> ghostnet.packets.MasterPacket
	9,  // 1: ghostnet.packets.MasterNodeUserInfoSq.User:type_name -> ghostnet.packets.GhostUser
	8,  // 2: ghostnet.packets.MasterNodeUserInfoCq.Master:type_name -> ghostnet.packets.MasterPacket
	9,  // 3: ghostnet.packets.MasterNodeUserInfoCq.User:type_name -> ghostnet.packets.GhostUser
	8,  // 4: ghostnet.packets.ReqeustMasterNodeListSq.Master:type_name -> ghostnet.packets.MasterPacket
	8,  // 5: ghostnet.packets.RequestMasterNodeListCq.Master:type_name -> ghostnet.packets.MasterPacket
	8,  // 6: ghostnet.packets.ResponseMasterNodeListSq.Master:type_name -> ghostnet.packets.MasterPacket
	9,  // 7: ghostnet.packets.ResponseMasterNodeListSq.User:type_name -> ghostnet.packets.GhostUser
	8,  // 8: ghostnet.packets.ResponseMasterNodeListCq.Master:type_name -> ghostnet.packets.MasterPacket
	8,  // 9: ghostnet.packets.SearchGhostPubKeySq.Master:type_name -> ghostnet.packets.MasterPacket
	8,  // 10: ghostnet.packets.SearchGhostPubKeyCq.Master:type_name -> ghostnet.packets.MasterPacket
	9,  // 11: ghostnet.packets.SearchGhostPubKeyCq.User:type_name -> ghostnet.packets.GhostUser
	12, // [12:12] is the sub-list for method output_type
	12, // [12:12] is the sub-list for method input_type
	12, // [12:12] is the sub-list for extension type_name
	12, // [12:12] is the sub-list for extension extendee
	0,  // [0:12] is the sub-list for field type_name
}

func init() { file_masternet_packet_proto_init() }
func file_masternet_packet_proto_init() {
	if File_masternet_packet_proto != nil {
		return
	}
	file_common_packet_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_masternet_packet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MasterNodeUserInfoSq); i {
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
		file_masternet_packet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MasterNodeUserInfoCq); i {
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
		file_masternet_packet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReqeustMasterNodeListSq); i {
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
		file_masternet_packet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestMasterNodeListCq); i {
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
		file_masternet_packet_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseMasterNodeListSq); i {
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
		file_masternet_packet_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseMasterNodeListCq); i {
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
		file_masternet_packet_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchGhostPubKeySq); i {
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
		file_masternet_packet_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SearchGhostPubKeyCq); i {
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
			RawDescriptor: file_masternet_packet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_masternet_packet_proto_goTypes,
		DependencyIndexes: file_masternet_packet_proto_depIdxs,
		MessageInfos:      file_masternet_packet_proto_msgTypes,
	}.Build()
	File_masternet_packet_proto = out.File
	file_masternet_packet_proto_rawDesc = nil
	file_masternet_packet_proto_goTypes = nil
	file_masternet_packet_proto_depIdxs = nil
}