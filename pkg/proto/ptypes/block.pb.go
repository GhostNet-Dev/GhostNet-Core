// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.6.1
// source: block.proto

package ptypes

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

type PairedBlocks struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Block     *GhostNetBlock     `protobuf:"bytes,1,opt,name=Block,proto3" json:"Block,omitempty"`
	DataBlock *GhostNetDataBlock `protobuf:"bytes,2,opt,name=DataBlock,proto3" json:"DataBlock,omitempty"`
}

func (x *PairedBlocks) Reset() {
	*x = PairedBlocks{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PairedBlocks) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PairedBlocks) ProtoMessage() {}

func (x *PairedBlocks) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PairedBlocks.ProtoReflect.Descriptor instead.
func (*PairedBlocks) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{0}
}

func (x *PairedBlocks) GetBlock() *GhostNetBlock {
	if x != nil {
		return x.Block
	}
	return nil
}

func (x *PairedBlocks) GetDataBlock() *GhostNetDataBlock {
	if x != nil {
		return x.DataBlock
	}
	return nil
}

type GhostNetBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header      *GhostNetBlockHeader `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	Alice       []*GhostTransaction  `protobuf:"bytes,2,rep,name=Alice,proto3" json:"Alice,omitempty"`
	Transaction []*GhostTransaction  `protobuf:"bytes,3,rep,name=Transaction,proto3" json:"Transaction,omitempty"`
}

func (x *GhostNetBlock) Reset() {
	*x = GhostNetBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GhostNetBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GhostNetBlock) ProtoMessage() {}

func (x *GhostNetBlock) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GhostNetBlock.ProtoReflect.Descriptor instead.
func (*GhostNetBlock) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{1}
}

func (x *GhostNetBlock) GetHeader() *GhostNetBlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *GhostNetBlock) GetAlice() []*GhostTransaction {
	if x != nil {
		return x.Alice
	}
	return nil
}

func (x *GhostNetBlock) GetTransaction() []*GhostTransaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

type GhostNetBlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                      uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Version                 uint32 `protobuf:"varint,2,opt,name=Version,proto3" json:"Version,omitempty"`
	PreviousBlockHeaderHash []byte `protobuf:"bytes,3,opt,name=PreviousBlockHeaderHash,proto3" json:"PreviousBlockHeaderHash,omitempty"`
	MerkleRoot              []byte `protobuf:"bytes,4,opt,name=MerkleRoot,proto3" json:"MerkleRoot,omitempty"`
	DataBlockHeaderHash     []byte `protobuf:"bytes,5,opt,name=DataBlockHeaderHash,proto3" json:"DataBlockHeaderHash,omitempty"`
	TimeStamp               uint64 `protobuf:"varint,6,opt,name=TimeStamp,proto3" json:"TimeStamp,omitempty"`
	Bits                    uint32 `protobuf:"varint,7,opt,name=Bits,proto3" json:"Bits,omitempty"`
	Nonce                   uint32 `protobuf:"varint,8,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	AliceCount              uint32 `protobuf:"varint,9,opt,name=AliceCount,proto3" json:"AliceCount,omitempty"`
	TransactionCount        uint32 `protobuf:"varint,10,opt,name=TransactionCount,proto3" json:"TransactionCount,omitempty"`
	SignatureSize           uint32 `protobuf:"varint,11,opt,name=SignatureSize,proto3" json:"SignatureSize,omitempty"`
	BlockSignature          []byte `protobuf:"bytes,12,opt,name=BlockSignature,proto3" json:"BlockSignature,omitempty"`
}

func (x *GhostNetBlockHeader) Reset() {
	*x = GhostNetBlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GhostNetBlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GhostNetBlockHeader) ProtoMessage() {}

func (x *GhostNetBlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GhostNetBlockHeader.ProtoReflect.Descriptor instead.
func (*GhostNetBlockHeader) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{2}
}

func (x *GhostNetBlockHeader) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GhostNetBlockHeader) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *GhostNetBlockHeader) GetPreviousBlockHeaderHash() []byte {
	if x != nil {
		return x.PreviousBlockHeaderHash
	}
	return nil
}

func (x *GhostNetBlockHeader) GetMerkleRoot() []byte {
	if x != nil {
		return x.MerkleRoot
	}
	return nil
}

func (x *GhostNetBlockHeader) GetDataBlockHeaderHash() []byte {
	if x != nil {
		return x.DataBlockHeaderHash
	}
	return nil
}

func (x *GhostNetBlockHeader) GetTimeStamp() uint64 {
	if x != nil {
		return x.TimeStamp
	}
	return 0
}

func (x *GhostNetBlockHeader) GetBits() uint32 {
	if x != nil {
		return x.Bits
	}
	return 0
}

func (x *GhostNetBlockHeader) GetNonce() uint32 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *GhostNetBlockHeader) GetAliceCount() uint32 {
	if x != nil {
		return x.AliceCount
	}
	return 0
}

func (x *GhostNetBlockHeader) GetTransactionCount() uint32 {
	if x != nil {
		return x.TransactionCount
	}
	return 0
}

func (x *GhostNetBlockHeader) GetSignatureSize() uint32 {
	if x != nil {
		return x.SignatureSize
	}
	return 0
}

func (x *GhostNetBlockHeader) GetBlockSignature() []byte {
	if x != nil {
		return x.BlockSignature
	}
	return nil
}

type GhostNetDataBlock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header      *GhostNetDataBlockHeader `protobuf:"bytes,1,opt,name=Header,proto3" json:"Header,omitempty"`
	Transaction []*GhostDataTransaction  `protobuf:"bytes,2,rep,name=Transaction,proto3" json:"Transaction,omitempty"`
}

func (x *GhostNetDataBlock) Reset() {
	*x = GhostNetDataBlock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GhostNetDataBlock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GhostNetDataBlock) ProtoMessage() {}

func (x *GhostNetDataBlock) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GhostNetDataBlock.ProtoReflect.Descriptor instead.
func (*GhostNetDataBlock) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{3}
}

func (x *GhostNetDataBlock) GetHeader() *GhostNetDataBlockHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *GhostNetDataBlock) GetTransaction() []*GhostDataTransaction {
	if x != nil {
		return x.Transaction
	}
	return nil
}

type GhostNetDataBlockHeader struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id                      uint32 `protobuf:"varint,1,opt,name=Id,proto3" json:"Id,omitempty"`
	Version                 uint32 `protobuf:"varint,2,opt,name=Version,proto3" json:"Version,omitempty"`
	PreviousBlockHeaderHash []byte `protobuf:"bytes,3,opt,name=PreviousBlockHeaderHash,proto3" json:"PreviousBlockHeaderHash,omitempty"`
	MerkleRoot              []byte `protobuf:"bytes,4,opt,name=MerkleRoot,proto3" json:"MerkleRoot,omitempty"`
	Nonce                   uint32 `protobuf:"varint,5,opt,name=Nonce,proto3" json:"Nonce,omitempty"`
	TransactionCount        uint32 `protobuf:"varint,6,opt,name=TransactionCount,proto3" json:"TransactionCount,omitempty"`
}

func (x *GhostNetDataBlockHeader) Reset() {
	*x = GhostNetDataBlockHeader{}
	if protoimpl.UnsafeEnabled {
		mi := &file_block_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GhostNetDataBlockHeader) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GhostNetDataBlockHeader) ProtoMessage() {}

func (x *GhostNetDataBlockHeader) ProtoReflect() protoreflect.Message {
	mi := &file_block_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GhostNetDataBlockHeader.ProtoReflect.Descriptor instead.
func (*GhostNetDataBlockHeader) Descriptor() ([]byte, []int) {
	return file_block_proto_rawDescGZIP(), []int{4}
}

func (x *GhostNetDataBlockHeader) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *GhostNetDataBlockHeader) GetVersion() uint32 {
	if x != nil {
		return x.Version
	}
	return 0
}

func (x *GhostNetDataBlockHeader) GetPreviousBlockHeaderHash() []byte {
	if x != nil {
		return x.PreviousBlockHeaderHash
	}
	return nil
}

func (x *GhostNetDataBlockHeader) GetMerkleRoot() []byte {
	if x != nil {
		return x.MerkleRoot
	}
	return nil
}

func (x *GhostNetDataBlockHeader) GetNonce() uint32 {
	if x != nil {
		return x.Nonce
	}
	return 0
}

func (x *GhostNetDataBlockHeader) GetTransactionCount() uint32 {
	if x != nil {
		return x.TransactionCount
	}
	return 0
}

var File_block_proto protoreflect.FileDescriptor

var file_block_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0f, 0x67,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x1a, 0x11,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x86, 0x01, 0x0a, 0x0c, 0x50, 0x61, 0x69, 0x72, 0x65, 0x64, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x73, 0x12, 0x34, 0x0a, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1e, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79,
	0x70, 0x65, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x52, 0x05, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x40, 0x0a, 0x09, 0x44, 0x61, 0x74, 0x61,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x22, 0x2e, 0x67, 0x68,
	0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x68,
	0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x52,
	0x09, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x22, 0xcb, 0x01, 0x0a, 0x0d, 0x47,
	0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x3c, 0x0a, 0x06,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x67,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47,
	0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x52, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x37, 0x0a, 0x05, 0x41, 0x6c,
	0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67, 0x68, 0x6f, 0x73,
	0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73,
	0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x05, 0x41, 0x6c,
	0x69, 0x63, 0x65, 0x12, 0x43, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x21, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74,
	0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xad, 0x03, 0x0a, 0x13, 0x47, 0x68, 0x6f,
	0x73, 0x74, 0x4e, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x49, 0x64,
	0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x38, 0x0a, 0x17, 0x50, 0x72,
	0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x17, 0x50, 0x72, 0x65,
	0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x48, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f,
	0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65,
	0x52, 0x6f, 0x6f, 0x74, 0x12, 0x30, 0x0a, 0x13, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x13, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1c, 0x0a, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53, 0x74,
	0x61, 0x6d, 0x70, 0x18, 0x06, 0x20, 0x01, 0x28, 0x04, 0x52, 0x09, 0x54, 0x69, 0x6d, 0x65, 0x53,
	0x74, 0x61, 0x6d, 0x70, 0x12, 0x12, 0x0a, 0x04, 0x42, 0x69, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x04, 0x42, 0x69, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e, 0x63,
	0x65, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12, 0x1e,
	0x0a, 0x0a, 0x41, 0x6c, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x0d, 0x52, 0x0a, 0x41, 0x6c, 0x69, 0x63, 0x65, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x2a,
	0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75,
	0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61,
	0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x24, 0x0a, 0x0d, 0x53, 0x69,
	0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x0d, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x53, 0x69, 0x7a, 0x65,
	0x12, 0x26, 0x0a, 0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53, 0x69, 0x67, 0x6e, 0x61, 0x74, 0x75,
	0x72, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0e, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x53,
	0x69, 0x67, 0x6e, 0x61, 0x74, 0x75, 0x72, 0x65, 0x22, 0x9e, 0x01, 0x0a, 0x11, 0x47, 0x68, 0x6f,
	0x73, 0x74, 0x4e, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x12, 0x40,
	0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28,
	0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73,
	0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x12, 0x47, 0x0a, 0x0b, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74,
	0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x44, 0x61, 0x74,
	0x61, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0b, 0x54, 0x72,
	0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xdf, 0x01, 0x0a, 0x17, 0x47, 0x68,
	0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x0e, 0x0a, 0x02, 0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x02, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x07, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12,
	0x38, 0x0a, 0x17, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x17, 0x50, 0x72, 0x65, 0x76, 0x69, 0x6f, 0x75, 0x73, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x48, 0x61, 0x73, 0x68, 0x12, 0x1e, 0x0a, 0x0a, 0x4d, 0x65, 0x72,
	0x6b, 0x6c, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x4d,
	0x65, 0x72, 0x6b, 0x6c, 0x65, 0x52, 0x6f, 0x6f, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x4e, 0x6f, 0x6e,
	0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x4e, 0x6f, 0x6e, 0x63, 0x65, 0x12,
	0x2a, 0x0a, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x10, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x38, 0x5a, 0x36, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e,
	0x65, 0x74, 0x2d, 0x44, 0x65, 0x76, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d,
	0x43, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70,
	0x74, 0x79, 0x70, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_block_proto_rawDescOnce sync.Once
	file_block_proto_rawDescData = file_block_proto_rawDesc
)

func file_block_proto_rawDescGZIP() []byte {
	file_block_proto_rawDescOnce.Do(func() {
		file_block_proto_rawDescData = protoimpl.X.CompressGZIP(file_block_proto_rawDescData)
	})
	return file_block_proto_rawDescData
}

var file_block_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_block_proto_goTypes = []interface{}{
	(*PairedBlocks)(nil),            // 0: ghostnet.ptypes.PairedBlocks
	(*GhostNetBlock)(nil),           // 1: ghostnet.ptypes.GhostNetBlock
	(*GhostNetBlockHeader)(nil),     // 2: ghostnet.ptypes.GhostNetBlockHeader
	(*GhostNetDataBlock)(nil),       // 3: ghostnet.ptypes.GhostNetDataBlock
	(*GhostNetDataBlockHeader)(nil), // 4: ghostnet.ptypes.GhostNetDataBlockHeader
	(*GhostTransaction)(nil),        // 5: ghostnet.ptypes.GhostTransaction
	(*GhostDataTransaction)(nil),    // 6: ghostnet.ptypes.GhostDataTransaction
}
var file_block_proto_depIdxs = []int32{
	1, // 0: ghostnet.ptypes.PairedBlocks.Block:type_name -> ghostnet.ptypes.GhostNetBlock
	3, // 1: ghostnet.ptypes.PairedBlocks.DataBlock:type_name -> ghostnet.ptypes.GhostNetDataBlock
	2, // 2: ghostnet.ptypes.GhostNetBlock.Header:type_name -> ghostnet.ptypes.GhostNetBlockHeader
	5, // 3: ghostnet.ptypes.GhostNetBlock.Alice:type_name -> ghostnet.ptypes.GhostTransaction
	5, // 4: ghostnet.ptypes.GhostNetBlock.Transaction:type_name -> ghostnet.ptypes.GhostTransaction
	4, // 5: ghostnet.ptypes.GhostNetDataBlock.Header:type_name -> ghostnet.ptypes.GhostNetDataBlockHeader
	6, // 6: ghostnet.ptypes.GhostNetDataBlock.Transaction:type_name -> ghostnet.ptypes.GhostDataTransaction
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_block_proto_init() }
func file_block_proto_init() {
	if File_block_proto != nil {
		return
	}
	file_transaction_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_block_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PairedBlocks); i {
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
		file_block_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GhostNetBlock); i {
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
		file_block_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GhostNetBlockHeader); i {
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
		file_block_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GhostNetDataBlock); i {
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
		file_block_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GhostNetDataBlockHeader); i {
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
			RawDescriptor: file_block_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_block_proto_goTypes,
		DependencyIndexes: file_block_proto_depIdxs,
		MessageInfos:      file_block_proto_msgTypes,
	}.Build()
	File_block_proto = out.File
	file_block_proto_rawDesc = nil
	file_block_proto_goTypes = nil
	file_block_proto_depIdxs = nil
}
