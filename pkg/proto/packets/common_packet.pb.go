// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        v3.12.4
// source: common_packet.proto

package packets

import (
	ptypes "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
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

type PacketType int32

const (
	PacketType_Reserved0 PacketType = 0
	// master node packet type
	PacketType_MasterNetwork PacketType = 1
	PacketType_FileTransfer  PacketType = 2
)

// Enum value maps for PacketType.
var (
	PacketType_name = map[int32]string{
		0: "Reserved0",
		1: "MasterNetwork",
		2: "FileTransfer",
	}
	PacketType_value = map[string]int32{
		"Reserved0":     0,
		"MasterNetwork": 1,
		"FileTransfer":  2,
	}
)

func (x PacketType) Enum() *PacketType {
	p := new(PacketType)
	*p = x
	return p
}

func (x PacketType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_packet_proto_enumTypes[0].Descriptor()
}

func (PacketType) Type() protoreflect.EnumType {
	return &file_common_packet_proto_enumTypes[0]
}

func (x PacketType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketType.Descriptor instead.
func (PacketType) EnumDescriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{0}
}

type PacketSecondType int32

const (
	// option allow_alias = true;
	PacketSecondType_GetGhostNetVersion     PacketSecondType = 0
	PacketSecondType_NotificationMasterNode PacketSecondType = 1
	PacketSecondType_ConnectToMasterNode    PacketSecondType = 2
	PacketSecondType_SearchGhostPubKey      PacketSecondType = 3
	PacketSecondType_RequestMasterNodeList  PacketSecondType = 4
	PacketSecondType_ResponseMasterNodeList PacketSecondType = 5
	PacketSecondType_SearchUserInfoByPubKey PacketSecondType = 6
	PacketSecondType_RegistBadBlock         PacketSecondType = 7
	PacketSecondType_BlockChain             PacketSecondType = 8
	PacketSecondType_Forwarding             PacketSecondType = 9
	PacketSecondType_RequestFile            PacketSecondType = 10
	PacketSecondType_ResponseFile           PacketSecondType = 11
)

// Enum value maps for PacketSecondType.
var (
	PacketSecondType_name = map[int32]string{
		0:  "GetGhostNetVersion",
		1:  "NotificationMasterNode",
		2:  "ConnectToMasterNode",
		3:  "SearchGhostPubKey",
		4:  "RequestMasterNodeList",
		5:  "ResponseMasterNodeList",
		6:  "SearchUserInfoByPubKey",
		7:  "RegistBadBlock",
		8:  "BlockChain",
		9:  "Forwarding",
		10: "RequestFile",
		11: "ResponseFile",
	}
	PacketSecondType_value = map[string]int32{
		"GetGhostNetVersion":     0,
		"NotificationMasterNode": 1,
		"ConnectToMasterNode":    2,
		"SearchGhostPubKey":      3,
		"RequestMasterNodeList":  4,
		"ResponseMasterNodeList": 5,
		"SearchUserInfoByPubKey": 6,
		"RegistBadBlock":         7,
		"BlockChain":             8,
		"Forwarding":             9,
		"RequestFile":            10,
		"ResponseFile":           11,
	}
)

func (x PacketSecondType) Enum() *PacketSecondType {
	p := new(PacketSecondType)
	*p = x
	return p
}

func (x PacketSecondType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketSecondType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_packet_proto_enumTypes[1].Descriptor()
}

func (PacketSecondType) Type() protoreflect.EnumType {
	return &file_common_packet_proto_enumTypes[1]
}

func (x PacketSecondType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketSecondType.Descriptor instead.
func (PacketSecondType) EnumDescriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{1}
}

type PacketThirdType int32

const (
	PacketThirdType_Reserved1             PacketThirdType = 0
	PacketThirdType_GetHeightestBlock     PacketThirdType = 1
	PacketThirdType_NewBlock              PacketThirdType = 2
	PacketThirdType_GetBlock              PacketThirdType = 3
	PacketThirdType_SendBlock             PacketThirdType = 4
	PacketThirdType_ScanAddrBlock         PacketThirdType = 5 // not used...
	PacketThirdType_SendTransaction       PacketThirdType = 6
	PacketThirdType_SearchTransaction     PacketThirdType = 7
	PacketThirdType_SendDataTransaction   PacketThirdType = 8
	PacketThirdType_SearchDataTransaction PacketThirdType = 9
	PacketThirdType_ScanBlockChain        PacketThirdType = 10
	PacketThirdType_CheckGhostNickname    PacketThirdType = 11
	PacketThirdType_SendDataTxIdList      PacketThirdType = 12
	PacketThirdType_GetDataTxIdList       PacketThirdType = 13
	PacketThirdType_ReportBlockError      PacketThirdType = 14
	PacketThirdType_GetBlockHash          PacketThirdType = 15
	PacketThirdType_SendBlockHash         PacketThirdType = 16
	PacketThirdType_GetBlockPrevHash      PacketThirdType = 17
	PacketThirdType_SendBlockPrevHash     PacketThirdType = 18
	PacketThirdType_GetTxStatus           PacketThirdType = 19
	PacketThirdType_SendTxStatus          PacketThirdType = 20
	PacketThirdType_CheckRootFs           PacketThirdType = 21
)

// Enum value maps for PacketThirdType.
var (
	PacketThirdType_name = map[int32]string{
		0:  "Reserved1",
		1:  "GetHeightestBlock",
		2:  "NewBlock",
		3:  "GetBlock",
		4:  "SendBlock",
		5:  "ScanAddrBlock",
		6:  "SendTransaction",
		7:  "SearchTransaction",
		8:  "SendDataTransaction",
		9:  "SearchDataTransaction",
		10: "ScanBlockChain",
		11: "CheckGhostNickname",
		12: "SendDataTxIdList",
		13: "GetDataTxIdList",
		14: "ReportBlockError",
		15: "GetBlockHash",
		16: "SendBlockHash",
		17: "GetBlockPrevHash",
		18: "SendBlockPrevHash",
		19: "GetTxStatus",
		20: "SendTxStatus",
		21: "CheckRootFs",
	}
	PacketThirdType_value = map[string]int32{
		"Reserved1":             0,
		"GetHeightestBlock":     1,
		"NewBlock":              2,
		"GetBlock":              3,
		"SendBlock":             4,
		"ScanAddrBlock":         5,
		"SendTransaction":       6,
		"SearchTransaction":     7,
		"SendDataTransaction":   8,
		"SearchDataTransaction": 9,
		"ScanBlockChain":        10,
		"CheckGhostNickname":    11,
		"SendDataTxIdList":      12,
		"GetDataTxIdList":       13,
		"ReportBlockError":      14,
		"GetBlockHash":          15,
		"SendBlockHash":         16,
		"GetBlockPrevHash":      17,
		"SendBlockPrevHash":     18,
		"GetTxStatus":           19,
		"SendTxStatus":          20,
		"CheckRootFs":           21,
	}
)

func (x PacketThirdType) Enum() *PacketThirdType {
	p := new(PacketThirdType)
	*p = x
	return p
}

func (x PacketThirdType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PacketThirdType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_packet_proto_enumTypes[2].Descriptor()
}

func (PacketThirdType) Type() protoreflect.EnumType {
	return &file_common_packet_proto_enumTypes[2]
}

func (x PacketThirdType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PacketThirdType.Descriptor instead.
func (PacketThirdType) EnumDescriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{2}
}

type RoutingType int32

const (
	RoutingType_None                  RoutingType = 0
	RoutingType_BroadCastingLevelZero RoutingType = 1
	RoutingType_BroadCasting          RoutingType = 2
	RoutingType_Flooding              RoutingType = 3
	RoutingType_SelectiveFlooding     RoutingType = 4
)

// Enum value maps for RoutingType.
var (
	RoutingType_name = map[int32]string{
		0: "None",
		1: "BroadCastingLevelZero",
		2: "BroadCasting",
		3: "Flooding",
		4: "SelectiveFlooding",
	}
	RoutingType_value = map[string]int32{
		"None":                  0,
		"BroadCastingLevelZero": 1,
		"BroadCasting":          2,
		"Flooding":              3,
		"SelectiveFlooding":     4,
	}
)

func (x RoutingType) Enum() *RoutingType {
	p := new(RoutingType)
	*p = x
	return p
}

func (x RoutingType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RoutingType) Descriptor() protoreflect.EnumDescriptor {
	return file_common_packet_proto_enumTypes[3].Descriptor()
}

func (RoutingType) Type() protoreflect.EnumType {
	return &file_common_packet_proto_enumTypes[3]
}

func (x RoutingType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RoutingType.Descriptor instead.
func (RoutingType) EnumDescriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{3}
}

type Header struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       PacketType       `protobuf:"varint,1,opt,name=Type,proto3,enum=ghostnet.packets.PacketType" json:"Type,omitempty"`
	SecondType PacketSecondType `protobuf:"varint,2,opt,name=SecondType,proto3,enum=ghostnet.packets.PacketSecondType" json:"SecondType,omitempty"`
	ThirdType  PacketThirdType  `protobuf:"varint,3,opt,name=ThirdType,proto3,enum=ghostnet.packets.PacketThirdType" json:"ThirdType,omitempty"`
	SqFlag     bool             `protobuf:"varint,4,opt,name=SqFlag,proto3" json:"SqFlag,omitempty"`
	PacketData []byte           `protobuf:"bytes,5,opt,name=PacketData,proto3" json:"PacketData,omitempty"`
	Source     *ptypes.GhostIp  `protobuf:"bytes,6,opt,name=Source,proto3" json:"Source,omitempty"`
}

func (x *Header) Reset() {
	*x = Header{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_packet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Header) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Header) ProtoMessage() {}

func (x *Header) ProtoReflect() protoreflect.Message {
	mi := &file_common_packet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Header.ProtoReflect.Descriptor instead.
func (*Header) Descriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{0}
}

func (x *Header) GetType() PacketType {
	if x != nil {
		return x.Type
	}
	return PacketType_Reserved0
}

func (x *Header) GetSecondType() PacketSecondType {
	if x != nil {
		return x.SecondType
	}
	return PacketSecondType_GetGhostNetVersion
}

func (x *Header) GetThirdType() PacketThirdType {
	if x != nil {
		return x.ThirdType
	}
	return PacketThirdType_Reserved1
}

func (x *Header) GetSqFlag() bool {
	if x != nil {
		return x.SqFlag
	}
	return false
}

func (x *Header) GetPacketData() []byte {
	if x != nil {
		return x.PacketData
	}
	return nil
}

func (x *Header) GetSource() *ptypes.GhostIp {
	if x != nil {
		return x.Source
	}
	return nil
}

type GhostPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FromPubKeyAddress string `protobuf:"bytes,1,opt,name=FromPubKeyAddress,proto3" json:"FromPubKeyAddress,omitempty"`
	RequestId         uint32 `protobuf:"varint,2,opt,name=RequestId,proto3" json:"RequestId,omitempty"`
	ClientId          uint32 `protobuf:"varint,3,opt,name=ClientId,proto3" json:"ClientId,omitempty"`
	TimeId            uint64 `protobuf:"varint,4,opt,name=TimeId,proto3" json:"TimeId,omitempty"`
}

func (x *GhostPacket) Reset() {
	*x = GhostPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_packet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GhostPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GhostPacket) ProtoMessage() {}

func (x *GhostPacket) ProtoReflect() protoreflect.Message {
	mi := &file_common_packet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GhostPacket.ProtoReflect.Descriptor instead.
func (*GhostPacket) Descriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{1}
}

func (x *GhostPacket) GetFromPubKeyAddress() string {
	if x != nil {
		return x.FromPubKeyAddress
	}
	return ""
}

func (x *GhostPacket) GetRequestId() uint32 {
	if x != nil {
		return x.RequestId
	}
	return 0
}

func (x *GhostPacket) GetClientId() uint32 {
	if x != nil {
		return x.ClientId
	}
	return 0
}

func (x *GhostPacket) GetTimeId() uint64 {
	if x != nil {
		return x.TimeId
	}
	return 0
}

type MasterPacket struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Common   *GhostPacket `protobuf:"bytes,1,opt,name=Common,proto3" json:"Common,omitempty"`
	RoutingT RoutingType  `protobuf:"varint,2,opt,name=RoutingT,proto3,enum=ghostnet.packets.RoutingType" json:"RoutingT,omitempty"`
	Level    uint32       `protobuf:"varint,3,opt,name=Level,proto3" json:"Level,omitempty"`
}

func (x *MasterPacket) Reset() {
	*x = MasterPacket{}
	if protoimpl.UnsafeEnabled {
		mi := &file_common_packet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MasterPacket) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MasterPacket) ProtoMessage() {}

func (x *MasterPacket) ProtoReflect() protoreflect.Message {
	mi := &file_common_packet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MasterPacket.ProtoReflect.Descriptor instead.
func (*MasterPacket) Descriptor() ([]byte, []int) {
	return file_common_packet_proto_rawDescGZIP(), []int{2}
}

func (x *MasterPacket) GetCommon() *GhostPacket {
	if x != nil {
		return x.Common
	}
	return nil
}

func (x *MasterPacket) GetRoutingT() RoutingType {
	if x != nil {
		return x.RoutingT
	}
	return RoutingType_None
}

func (x *MasterPacket) GetLevel() uint32 {
	if x != nil {
		return x.Level
	}
	return 0
}

var File_common_packet_proto protoreflect.FileDescriptor

var file_common_packet_proto_rawDesc = []byte{
	0x0a, 0x13, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x1a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0xa9, 0x02, 0x0a, 0x06, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x30,
	0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x67,
	0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e,
	0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x42, 0x0a, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e,
	0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x65,
	0x63, 0x6f, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0a, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x3f, 0x0a, 0x09, 0x54, 0x68, 0x69, 0x72, 0x64, 0x54, 0x79, 0x70,
	0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x21, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e,
	0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x50, 0x61, 0x63, 0x6b, 0x65,
	0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x09, 0x54, 0x68, 0x69, 0x72,
	0x64, 0x54, 0x79, 0x70, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x53, 0x71, 0x46, 0x6c, 0x61, 0x67, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x53, 0x71, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x1e, 0x0a,
	0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x44, 0x61, 0x74, 0x61, 0x12, 0x30, 0x0a,
	0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e,
	0x47, 0x68, 0x6f, 0x73, 0x74, 0x49, 0x70, 0x52, 0x06, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x22,
	0x8d, 0x01, 0x0a, 0x0b, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x12,
	0x2c, 0x0a, 0x11, 0x46, 0x72, 0x6f, 0x6d, 0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x46, 0x72, 0x6f, 0x6d,
	0x50, 0x75, 0x62, 0x4b, 0x65, 0x79, 0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x09, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x08, 0x43,
	0x6c, 0x69, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x54, 0x69, 0x6d, 0x65, 0x49,
	0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x54, 0x69, 0x6d, 0x65, 0x49, 0x64, 0x22,
	0x96, 0x01, 0x0a, 0x0c, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74,
	0x12, 0x35, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x1d, 0x2e, 0x67, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x73, 0x2e, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x52,
	0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x12, 0x39, 0x0a, 0x08, 0x52, 0x6f, 0x75, 0x74, 0x69,
	0x6e, 0x67, 0x54, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x67, 0x68, 0x6f, 0x73,
	0x74, 0x6e, 0x65, 0x74, 0x2e, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x2e, 0x52, 0x6f, 0x75,
	0x74, 0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x52, 0x6f, 0x75, 0x74, 0x69, 0x6e,
	0x67, 0x54, 0x12, 0x14, 0x0a, 0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x05, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x2a, 0x40, 0x0a, 0x0a, 0x50, 0x61, 0x63, 0x6b,
	0x65, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76,
	0x65, 0x64, 0x30, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e,
	0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c, 0x46, 0x69, 0x6c, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x10, 0x02, 0x2a, 0xa0, 0x02, 0x0a, 0x10, 0x50,
	0x61, 0x63, 0x6b, 0x65, 0x74, 0x53, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x54, 0x79, 0x70, 0x65, 0x12,
	0x16, 0x0a, 0x12, 0x47, 0x65, 0x74, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x56, 0x65,
	0x72, 0x73, 0x69, 0x6f, 0x6e, 0x10, 0x00, 0x12, 0x1a, 0x0a, 0x16, 0x4e, 0x6f, 0x74, 0x69, 0x66,
	0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64,
	0x65, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x54, 0x6f,
	0x4d, 0x61, 0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x10, 0x02, 0x12, 0x15, 0x0a, 0x11,
	0x53, 0x65, 0x61, 0x72, 0x63, 0x68, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x50, 0x75, 0x62, 0x4b, 0x65,
	0x79, 0x10, 0x03, 0x12, 0x19, 0x0a, 0x15, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x4d, 0x61,
	0x73, 0x74, 0x65, 0x72, 0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x04, 0x12, 0x1a,
	0x0a, 0x16, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x4d, 0x61, 0x73, 0x74, 0x65, 0x72,
	0x4e, 0x6f, 0x64, 0x65, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x05, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x55, 0x73, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x42, 0x79, 0x50, 0x75,
	0x62, 0x4b, 0x65, 0x79, 0x10, 0x06, 0x12, 0x12, 0x0a, 0x0e, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x42, 0x61, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x10, 0x07, 0x12, 0x0e, 0x0a, 0x0a, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x10, 0x08, 0x12, 0x0e, 0x0a, 0x0a, 0x46, 0x6f,
	0x72, 0x77, 0x61, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x10, 0x09, 0x12, 0x0f, 0x0a, 0x0b, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x10, 0x0a, 0x12, 0x10, 0x0a, 0x0c, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x10, 0x0b, 0x2a, 0xc8, 0x03,
	0x0a, 0x0f, 0x50, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x54, 0x68, 0x69, 0x72, 0x64, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x0d, 0x0a, 0x09, 0x52, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x31, 0x10, 0x00,
	0x12, 0x15, 0x0a, 0x11, 0x47, 0x65, 0x74, 0x48, 0x65, 0x69, 0x67, 0x68, 0x74, 0x65, 0x73, 0x74,
	0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x10, 0x02, 0x12, 0x0c, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x10, 0x03, 0x12, 0x0d, 0x0a, 0x09, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b,
	0x10, 0x04, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x63, 0x61, 0x6e, 0x41, 0x64, 0x64, 0x72, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x06, 0x12, 0x15, 0x0a, 0x11, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x10,
	0x07, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x08, 0x12, 0x19, 0x0a, 0x15, 0x53, 0x65,
	0x61, 0x72, 0x63, 0x68, 0x44, 0x61, 0x74, 0x61, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x10, 0x09, 0x12, 0x12, 0x0a, 0x0e, 0x53, 0x63, 0x61, 0x6e, 0x42, 0x6c, 0x6f,
	0x63, 0x6b, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x10, 0x0a, 0x12, 0x16, 0x0a, 0x12, 0x43, 0x68, 0x65,
	0x63, 0x6b, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x69, 0x63, 0x6b, 0x6e, 0x61, 0x6d, 0x65, 0x10,
	0x0b, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x65, 0x6e, 0x64, 0x44, 0x61, 0x74, 0x61, 0x54, 0x78, 0x49,
	0x64, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x0c, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x44, 0x61,
	0x74, 0x61, 0x54, 0x78, 0x49, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x10, 0x0d, 0x12, 0x14, 0x0a, 0x10,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x45, 0x72, 0x72, 0x6f, 0x72,
	0x10, 0x0e, 0x12, 0x10, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x48, 0x61,
	0x73, 0x68, 0x10, 0x0f, 0x12, 0x11, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63,
	0x6b, 0x48, 0x61, 0x73, 0x68, 0x10, 0x10, 0x12, 0x14, 0x0a, 0x10, 0x47, 0x65, 0x74, 0x42, 0x6c,
	0x6f, 0x63, 0x6b, 0x50, 0x72, 0x65, 0x76, 0x48, 0x61, 0x73, 0x68, 0x10, 0x11, 0x12, 0x15, 0x0a,
	0x11, 0x53, 0x65, 0x6e, 0x64, 0x42, 0x6c, 0x6f, 0x63, 0x6b, 0x50, 0x72, 0x65, 0x76, 0x48, 0x61,
	0x73, 0x68, 0x10, 0x12, 0x12, 0x0f, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x54, 0x78, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x10, 0x13, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x65, 0x6e, 0x64, 0x54, 0x78, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x10, 0x14, 0x12, 0x0f, 0x0a, 0x0b, 0x43, 0x68, 0x65, 0x63, 0x6b,
	0x52, 0x6f, 0x6f, 0x74, 0x46, 0x73, 0x10, 0x15, 0x2a, 0x69, 0x0a, 0x0b, 0x52, 0x6f, 0x75, 0x74,
	0x69, 0x6e, 0x67, 0x54, 0x79, 0x70, 0x65, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x6f, 0x6e, 0x65, 0x10,
	0x00, 0x12, 0x19, 0x0a, 0x15, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x43, 0x61, 0x73, 0x74, 0x69, 0x6e,
	0x67, 0x4c, 0x65, 0x76, 0x65, 0x6c, 0x5a, 0x65, 0x72, 0x6f, 0x10, 0x01, 0x12, 0x10, 0x0a, 0x0c,
	0x42, 0x72, 0x6f, 0x61, 0x64, 0x43, 0x61, 0x73, 0x74, 0x69, 0x6e, 0x67, 0x10, 0x02, 0x12, 0x0c,
	0x0a, 0x08, 0x46, 0x6c, 0x6f, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11,
	0x53, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x76, 0x65, 0x46, 0x6c, 0x6f, 0x6f, 0x64, 0x69, 0x6e,
	0x67, 0x10, 0x04, 0x42, 0x39, 0x5a, 0x37, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x47, 0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d, 0x44, 0x65, 0x76, 0x2f, 0x47,
	0x68, 0x6f, 0x73, 0x74, 0x4e, 0x65, 0x74, 0x2d, 0x43, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x6b, 0x67,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x61, 0x63, 0x6b, 0x65, 0x74, 0x73, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_common_packet_proto_rawDescOnce sync.Once
	file_common_packet_proto_rawDescData = file_common_packet_proto_rawDesc
)

func file_common_packet_proto_rawDescGZIP() []byte {
	file_common_packet_proto_rawDescOnce.Do(func() {
		file_common_packet_proto_rawDescData = protoimpl.X.CompressGZIP(file_common_packet_proto_rawDescData)
	})
	return file_common_packet_proto_rawDescData
}

var file_common_packet_proto_enumTypes = make([]protoimpl.EnumInfo, 4)
var file_common_packet_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_common_packet_proto_goTypes = []interface{}{
	(PacketType)(0),        // 0: ghostnet.packets.PacketType
	(PacketSecondType)(0),  // 1: ghostnet.packets.PacketSecondType
	(PacketThirdType)(0),   // 2: ghostnet.packets.PacketThirdType
	(RoutingType)(0),       // 3: ghostnet.packets.RoutingType
	(*Header)(nil),         // 4: ghostnet.packets.Header
	(*GhostPacket)(nil),    // 5: ghostnet.packets.GhostPacket
	(*MasterPacket)(nil),   // 6: ghostnet.packets.MasterPacket
	(*ptypes.GhostIp)(nil), // 7: ghostnet.ptypes.GhostIp
}
var file_common_packet_proto_depIdxs = []int32{
	0, // 0: ghostnet.packets.Header.Type:type_name -> ghostnet.packets.PacketType
	1, // 1: ghostnet.packets.Header.SecondType:type_name -> ghostnet.packets.PacketSecondType
	2, // 2: ghostnet.packets.Header.ThirdType:type_name -> ghostnet.packets.PacketThirdType
	7, // 3: ghostnet.packets.Header.Source:type_name -> ghostnet.ptypes.GhostIp
	5, // 4: ghostnet.packets.MasterPacket.Common:type_name -> ghostnet.packets.GhostPacket
	3, // 5: ghostnet.packets.MasterPacket.RoutingT:type_name -> ghostnet.packets.RoutingType
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_common_packet_proto_init() }
func file_common_packet_proto_init() {
	if File_common_packet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_common_packet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Header); i {
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
		file_common_packet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GhostPacket); i {
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
		file_common_packet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MasterPacket); i {
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
			RawDescriptor: file_common_packet_proto_rawDesc,
			NumEnums:      4,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_common_packet_proto_goTypes,
		DependencyIndexes: file_common_packet_proto_depIdxs,
		EnumInfos:         file_common_packet_proto_enumTypes,
		MessageInfos:      file_common_packet_proto_msgTypes,
	}.Build()
	File_common_packet_proto = out.File
	file_common_packet_proto_rawDesc = nil
	file_common_packet_proto_goTypes = nil
	file_common_packet_proto_depIdxs = nil
}
