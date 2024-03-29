// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.6.1
// source: file_packet.proto

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

type FileRequestType int32

const (
	FileRequestType_GetFileInfo FileRequestType = 0
	FileRequestType_GetFileData FileRequestType = 1
)

// Enum value maps for FileRequestType.
var (
	FileRequestType_name = map[int32]string{
		0: "GetFileInfo",
		1: "GetFileData",
	}
	FileRequestType_value = map[string]int32{
		"GetFileInfo": 0,
		"GetFileData": 1,
	}
)

func (x FileRequestType) Enum() *FileRequestType {
	p := new(FileRequestType)
	*p = x
	return p
}

func (x FileRequestType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (FileRequestType) Descriptor() protoreflect.EnumDescriptor {
	return file_file_packet_proto_enumTypes[0].Descriptor()
}

func (FileRequestType) Type() protoreflect.EnumType {
	return &file_file_packet_proto_enumTypes[0]
}

func (x FileRequestType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use FileRequestType.Descriptor instead.
func (FileRequestType) EnumDescriptor() ([]byte, []int) {
	return file_file_packet_proto_rawDescGZIP(), []int{0}
}

type RequestFilePacketSq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master      *MasterPacket   `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	RequestType FileRequestType `protobuf:"varint,2,opt,name=RequestType,proto3,enum=ghostnet.packets.FileRequestType" json:"RequestType,omitempty"`
	Filename    string          `protobuf:"bytes,3,opt,name=Filename,proto3" json:"Filename,omitempty"`
	StartOffset uint64          `protobuf:"varint,4,opt,name=StartOffset,proto3" json:"StartOffset,omitempty"`
}

func (x *RequestFilePacketSq) Reset() {
	*x = RequestFilePacketSq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_packet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestFilePacketSq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestFilePacketSq) ProtoMessage() {}

func (x *RequestFilePacketSq) ProtoReflect() protoreflect.Message {
	mi := &file_file_packet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestFilePacketSq.ProtoReflect.Descriptor instead.
func (*RequestFilePacketSq) Descriptor() ([]byte, []int) {
	return file_file_packet_proto_rawDescGZIP(), []int{0}
}

func (x *RequestFilePacketSq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *RequestFilePacketSq) GetRequestType() FileRequestType {
	if x != nil {
		return x.RequestType
	}
	return FileRequestType_GetFileInfo
}

func (x *RequestFilePacketSq) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *RequestFilePacketSq) GetStartOffset() uint64 {
	if x != nil {
		return x.StartOffset
	}
	return 0
}

type RequestFilePacketCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master      *MasterPacket   `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	RequestType FileRequestType `protobuf:"varint,2,opt,name=RequestType,proto3,enum=ghostnet.packets.FileRequestType" json:"RequestType,omitempty"`
	Filename    string          `protobuf:"bytes,3,opt,name=Filename,proto3" json:"Filename,omitempty"`
	FileLength  uint64          `protobuf:"varint,4,opt,name=FileLength,proto3" json:"FileLength,omitempty"`
	StartOffset uint64          `protobuf:"varint,5,opt,name=StartOffset,proto3" json:"StartOffset,omitempty"`
	Result      bool            `protobuf:"varint,6,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *RequestFilePacketCq) Reset() {
	*x = RequestFilePacketCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_packet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RequestFilePacketCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RequestFilePacketCq) ProtoMessage() {}

func (x *RequestFilePacketCq) ProtoReflect() protoreflect.Message {
	mi := &file_file_packet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RequestFilePacketCq.ProtoReflect.Descriptor instead.
func (*RequestFilePacketCq) Descriptor() ([]byte, []int) {
	return file_file_packet_proto_rawDescGZIP(), []int{1}
}

func (x *RequestFilePacketCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *RequestFilePacketCq) GetRequestType() FileRequestType {
	if x != nil {
		return x.RequestType
	}
	return FileRequestType_GetFileInfo
}

func (x *RequestFilePacketCq) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *RequestFilePacketCq) GetFileLength() uint64 {
	if x != nil {
		return x.FileLength
	}
	return 0
}

func (x *RequestFilePacketCq) GetStartOffset() uint64 {
	if x != nil {
		return x.StartOffset
	}
	return 0
}

func (x *RequestFilePacketCq) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

type ResponseFilePacketSq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master      *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	SequenceNum uint32        `protobuf:"varint,2,opt,name=SequenceNum,proto3" json:"SequenceNum,omitempty"`
	StartPos    uint64        `protobuf:"varint,3,opt,name=StartPos,proto3" json:"StartPos,omitempty"`
	Filename    string        `protobuf:"bytes,4,opt,name=Filename,proto3" json:"Filename,omitempty"`
	FileData    []byte        `protobuf:"bytes,5,opt,name=FileData,proto3" json:"FileData,omitempty"`
	BufferSize  uint32        `protobuf:"varint,6,opt,name=BufferSize,proto3" json:"BufferSize,omitempty"`
	FileLength  uint64        `protobuf:"varint,7,opt,name=FileLength,proto3" json:"FileLength,omitempty"`
	TimeId      uint64        `protobuf:"varint,8,opt,name=TimeId,proto3" json:"TimeId,omitempty"`
}

func (x *ResponseFilePacketSq) Reset() {
	*x = ResponseFilePacketSq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_packet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseFilePacketSq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseFilePacketSq) ProtoMessage() {}

func (x *ResponseFilePacketSq) ProtoReflect() protoreflect.Message {
	mi := &file_file_packet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseFilePacketSq.ProtoReflect.Descriptor instead.
func (*ResponseFilePacketSq) Descriptor() ([]byte, []int) {
	return file_file_packet_proto_rawDescGZIP(), []int{2}
}

func (x *ResponseFilePacketSq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *ResponseFilePacketSq) GetSequenceNum() uint32 {
	if x != nil {
		return x.SequenceNum
	}
	return 0
}

func (x *ResponseFilePacketSq) GetStartPos() uint64 {
	if x != nil {
		return x.StartPos
	}
	return 0
}

func (x *ResponseFilePacketSq) GetFilename() string {
	if x != nil {
		return x.Filename
	}
	return ""
}

func (x *ResponseFilePacketSq) GetFileData() []byte {
	if x != nil {
		return x.FileData
	}
	return nil
}

func (x *ResponseFilePacketSq) GetBufferSize() uint32 {
	if x != nil {
		return x.BufferSize
	}
	return 0
}

func (x *ResponseFilePacketSq) GetFileLength() uint64 {
	if x != nil {
		return x.FileLength
	}
	return 0
}

func (x *ResponseFilePacketSq) GetTimeId() uint64 {
	if x != nil {
		return x.TimeId
	}
	return 0
}

type ResponseFilePacketCq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Master      *MasterPacket `protobuf:"bytes,1,opt,name=Master,proto3" json:"Master,omitempty"`
	SequenceNum uint32        `protobuf:"varint,2,opt,name=SequenceNum,proto3" json:"SequenceNum,omitempty"`
	Result      bool          `protobuf:"varint,3,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *ResponseFilePacketCq) Reset() {
	*x = ResponseFilePacketCq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_file_packet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResponseFilePacketCq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResponseFilePacketCq) ProtoMessage() {}

func (x *ResponseFilePacketCq) ProtoReflect() protoreflect.Message {
	mi := &file_file_packet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResponseFilePacketCq.ProtoReflect.Descriptor instead.
func (*ResponseFilePacketCq) Descriptor() ([]byte, []int) {
	return file_file_packet_proto_rawDescGZIP(), []int{3}
}

func (x *ResponseFilePacketCq) GetMaster() *MasterPacket {
	if x != nil {
		return x.Master
	}
	return nil
}

func (x *ResponseFilePacketCq) GetSequenceNum() uint32 {
	if x != nil {
		return x.SequenceNum
	}
	return 0
}

func (x *ResponseFilePacketCq) GetResult() bool {
	if x != nil {
		return x.Result
	}
	return false
}

var File_file_packet_proto protoreflect.FileDescriptor

var file_file_packet_proto_rawDesc = []byte{
	0x0a, 0x11, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x10, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x1a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xd0, 0x01, 0x0a, 0x13, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x53, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x43, 0x0a, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x21, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x52, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x22, 0x88, 0x02,
	0x0a, 0x13, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x43, 0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74,
	0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x43, 0x0a,
	0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x21, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61,
	0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12, 0x20,
	0x0a, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0b, 0x53, 0x74, 0x61, 0x72, 0x74, 0x4f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x12, 0x16, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x9c, 0x02, 0x0a, 0x14, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x53,
	0x71, 0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63,
	0x6b, 0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x71,
	0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b,
	0x53, 0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x08, 0x53,
	0x74, 0x61, 0x72, 0x74, 0x50, 0x6f, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x46, 0x69, 0x6c, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12,
	0x1e, 0x0a, 0x0a, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x42, 0x75, 0x66, 0x66, 0x65, 0x72, 0x53, 0x69, 0x7a, 0x65, 0x12,
	0x1e, 0x0a, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x46, 0x69, 0x6c, 0x65, 0x4c, 0x65, 0x6e, 0x67, 0x74, 0x68, 0x12,
	0x16, 0x0a, 0x06, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x64, 0x18, 0x08, 0x20, 0x01, 0x28, 0x04, 0x52,
	0x06, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x64, 0x22, 0x88, 0x01, 0x0a, 0x14, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x43, 0x71,
	0x12, 0x36, 0x0a, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x2e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x52, 0x06, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x12, 0x20, 0x0a, 0x0b, 0x53, 0x65, 0x71, 0x75,
	0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x0b, 0x53,
	0x65, 0x71, 0x75, 0x65, 0x6e, 0x63, 0x65, 0x4e, 0x75, 0x6d, 0x12, 0x16, 0x0a, 0x06, 0x52, 0x65,
	0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x2a, 0x33, 0x0a, 0x0f, 0x46, 0x69, 0x6c, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c, 0x65,
	0x49, 0x6e, 0x66, 0x6f, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x46, 0x69, 0x6c,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x10, 0x01, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d, 0x44,
	0x65, 0x76, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d, 0x43, 0x6f, 0x72, 0x65,
	0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_file_packet_proto_rawDescOnce sync.Once
	file_file_packet_proto_rawDescData = file_file_packet_proto_rawDesc
)

func file_file_packet_proto_rawDescGZIP() []byte {
	file_file_packet_proto_rawDescOnce.Do(func() {
		file_file_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_file_packet_proto_rawDescData)
	})
	return file_file_packet_proto_rawDescData
}

var file_file_packet_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_file_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_file_packet_proto_goTypes = []interface{}{
	(FileRequestType)(0),         // 0: ghostnet.packets.FileRequestType
	(*RequestFilePacketSq)(nil),  // 1: ghostnet.packets.RequestFilePacketSq
	(*RequestFilePacketCq)(nil),  // 2: ghostnet.packets.RequestFilePacketCq
	(*ResponseFilePacketSq)(nil), // 3: ghostnet.packets.ResponseFilePacketSq
	(*ResponseFilePacketCq)(nil), // 4: ghostnet.packets.ResponseFilePacketCq
	(*MasterPacket)(nil),         // 5: ghostnet.packets.MasterPacket
}
var file_file_packet_proto_depIdxs = []int32{
	5, // 0: ghostnet.packets.RequestFilePacketSq.Master:type_name -> ghostnet.packets.MasterPacket
	0, // 1: ghostnet.packets.RequestFilePacketSq.RequestType:type_name -> ghostnet.packets.FileRequestType
	5, // 2: ghostnet.packets.RequestFilePacketCq.Master:type_name -> ghostnet.packets.MasterPacket
	0, // 3: ghostnet.packets.RequestFilePacketCq.RequestType:type_name -> ghostnet.packets.FileRequestType
	5, // 4: ghostnet.packets.ResponseFilePacketSq.Master:type_name -> ghostnet.packets.MasterPacket
	5, // 5: ghostnet.packets.ResponseFilePacketCq.Master:type_name -> ghostnet.packets.MasterPacket
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_file_packet_proto_init() }
func file_file_packet_proto_init() {
	if File_file_packet_proto != nil {
		return
	}
	file_common_packet_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_file_packet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestFilePacketSq); i {
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
		file_file_packet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RequestFilePacketCq); i {
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
		file_file_packet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseFilePacketSq); i {
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
		file_file_packet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResponseFilePacketCq); i {
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
			RawDescriptor: file_file_packet_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_file_packet_proto_goTypes,
		DependencyIndexes: file_file_packet_proto_depIdxs,
		EnumInfos:         file_file_packet_proto_enumTypes,
		MessageInfos:      file_file_packet_proto_msgTypes,
	}.Build()
	File_file_packet_proto = out.File
	file_file_packet_proto_rawDesc = nil
	file_file_packet_proto_goTypes = nil
	file_file_packet_proto_depIdxs = nil
}
