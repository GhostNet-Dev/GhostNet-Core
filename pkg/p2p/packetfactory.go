package p2p

import (
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type PacketHandler interface {
	GetGhostNetVersionSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //GhostNetVersion = 0;
	GetGhostNetVersionCq(packet []byte, from *net.UDPAddr)                      //GhostNetVersion = 0;
	//*SendPacketInfo // master node packet type
	NotificationMasterNodeSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //NotificationMasterNode = 1;
	NotificationMasterNodeCq(packet []byte, from *net.UDPAddr)                      //NotificationMasterNode = 1;
	ConnectToMasterNodeSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo    //ConnectToMasterNode = 3;
	ConnectToMasterNodeCq(packet []byte, from *net.UDPAddr)                         //ConnectToMasterNode = 3;
	SearchGhostPubKeySq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo      //SearchGhostPubKey = 4;
	SearchGhostPubKeyCq(packet []byte, from *net.UDPAddr)                           //SearchGhostPubKey = 4;
	RequestMasterNodeListSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo  //RequestMasterNodeList = 5;
	RequestMasterNodeListCq(packet []byte, from *net.UDPAddr)                       //RequestMasterNodeList = 5;
	ResponseMasterNodeListSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //ResponseMasterNodeList = 6;
	ResponseMasterNodeListCq(packet []byte, from *net.UDPAddr)                      //ResponseMasterNodeList = 6;
	SearchMasterPubKeySq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo     //SearchMasterPubKey = 9;
	SearchMasterPubKeyCq(packet []byte, from *net.UDPAddr)                          //SearchMasterPubKey = 9;
	//RegistBadBlock = 11;

	//*SendPacketInfo // blockchain packet type
	GetHeightestBlockSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //GetHeightestBlock = 12;
	GetHeightestBlockCq(packet []byte, from *net.UDPAddr)                      //GetHeightestBlock = 12;
	NewBlockSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo          //NewBlock = 13;
	NewBlockCq(packet []byte, from *net.UDPAddr)                               //NewBlock = 13;
	GetBlockSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo          //GetBlock = 14;
	GetBlockCq(packet []byte, from *net.UDPAddr)                               //GetBlock = 14;
	SendBlockSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo         //SendBlock = 15;
	SendBlockCq(packet []byte, from *net.UDPAddr)                              //SendBlock = 15;
	//ScanAddrBlock = 16; *SendPacketInfo // not used...
	SendTransactionSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo       //SendTransaction = 17;
	SendTransactionCq(packet []byte, from *net.UDPAddr)                            //SendTransaction = 17;
	SearchTransactionSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo     //SearchTransaction = 18;
	SearchTransactionCq(packet []byte, from *net.UDPAddr)                          //SearchTransaction = 18;
	SendDataTransactionSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo   //SendDataTransaction = 19;
	SendDataTransactionCq(packet []byte, from *net.UDPAddr)                        //SendDataTransaction = 19;
	SearchDataTransactionSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //SearchDataTransaction = 20;
	SearchDataTransactionCq(packet []byte, from *net.UDPAddr)                      //SearchDataTransaction = 20;
	//ScanBlockChain = 21;
	//CheckGhostNickname = 22;
	SendDataTxIdListSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //SendDataTxIdList = 23;
	SendDataTxIdListCq(packet []byte, from *net.UDPAddr)                      //SendDataTxIdList = 23;
	GetDataTxIdListSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo  //GetDataTxIdList = 24;
	GetDataTxIdListCq(packet []byte, from *net.UDPAddr)                       //GetDataTxIdList = 24;
	//ReportBlockError = 25;
	GetBlockHashSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo      //GetBlockHash = 26;
	GetBlockHashCq(packet []byte, from *net.UDPAddr)                           //GetBlockHash = 26;
	SendBlockHashSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo     //SendBlockHash = 27;
	SendBlockHashCq(packet []byte, from *net.UDPAddr)                          //SendBlockHash = 27;
	GetBlockPrevHashSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo  //GetBlockPrevHash = 28;
	GetBlockPrevHashCq(packet []byte, from *net.UDPAddr)                       //GetBlockPrevHash = 28;
	SendBlockPrevHashSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo //SendBlockPrevHash = 29;
	SendBlockPrevHashCq(packet []byte, from *net.UDPAddr)                      //SendBlockPrevHash = 29;
	GetTxStatusSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo       //GetTxStatus = 30;
	GetTxStatusCq(packet []byte, from *net.UDPAddr)                            //GetTxStatus = 30;
	SendTxStatusSq(packet []byte, from *net.UDPAddr) []ResponsePacketInfo      //SendTxStatus = 31;
	SendTxStatusCq(packet []byte, from *net.UDPAddr)                           //SendTxStatus = 31;

}

type PacketFactory struct {
	packetSqHandler map[packets.PacketType]func([]byte, *net.UDPAddr) []ResponsePacketInfo
	packetCqHandler map[packets.PacketType]func([]byte, *net.UDPAddr)
}

func NewPacketFactory() *PacketFactory {
	return &PacketFactory{
		packetSqHandler: make(map[packets.PacketType]func([]byte, *net.UDPAddr) []ResponsePacketInfo),
		packetCqHandler: make(map[packets.PacketType]func([]byte, *net.UDPAddr)),
	}
}

func MakeMasterPacket(from string, reqId uint32, clientId uint32, fromIp *ptypes.GhostIp) *packets.MasterPacket {
	return &packets.MasterPacket{
		Common: &packets.GhostPacket{
			FromPubKeyAddress: from,
			RequestId:         reqId,
			ClientId:          clientId,
			TimeId:            uint64(time.Now().Unix()),
		},
		RoutingT: packets.RoutingType_None,
		Level:    0,
	}
}

func (pf *PacketFactory) SingleRegisterPacketHandler(packetType packets.PacketType,
	sqHandler func([]byte, *net.UDPAddr) []ResponsePacketInfo, cqHandler func([]byte, *net.UDPAddr)) {
	pf.packetSqHandler[packetType] = sqHandler
	pf.packetCqHandler[packetType] = cqHandler
}

func (pf *PacketFactory) RegisterPacketHandler(handler PacketHandler) {
	pf.packetSqHandler[packets.PacketType_GetGhostNetVersion] = handler.GetGhostNetVersionSq
	pf.packetSqHandler[packets.PacketType_NotificationMasterNode] = handler.NotificationMasterNodeSq
	pf.packetSqHandler[packets.PacketType_ConnectToMasterNode] = handler.ConnectToMasterNodeSq
	pf.packetSqHandler[packets.PacketType_SearchGhostPubKey] = handler.SearchGhostPubKeySq
	pf.packetSqHandler[packets.PacketType_RequestMasterNodeList] = handler.RequestMasterNodeListSq
	pf.packetSqHandler[packets.PacketType_ResponseMasterNodeList] = handler.ResponseMasterNodeListSq
	pf.packetSqHandler[packets.PacketType_SearchMasterPubKey] = handler.SearchMasterPubKeySq

	pf.packetSqHandler[packets.PacketType_GetHeightestBlock] = handler.GetHeightestBlockSq
	pf.packetSqHandler[packets.PacketType_NewBlock] = handler.NewBlockSq
	pf.packetSqHandler[packets.PacketType_GetBlock] = handler.GetBlockSq
	pf.packetSqHandler[packets.PacketType_SendBlock] = handler.SendBlockSq
	pf.packetSqHandler[packets.PacketType_SendTransaction] = handler.SendTransactionSq
	pf.packetSqHandler[packets.PacketType_SendDataTransaction] = handler.SendDataTransactionSq
	pf.packetSqHandler[packets.PacketType_SearchDataTransaction] = handler.SearchDataTransactionSq
	pf.packetSqHandler[packets.PacketType_SendDataTxIdList] = handler.SendDataTxIdListSq
	pf.packetSqHandler[packets.PacketType_GetDataTxIdList] = handler.GetDataTxIdListSq
	pf.packetSqHandler[packets.PacketType_GetBlockHash] = handler.GetBlockHashSq
	pf.packetSqHandler[packets.PacketType_SendBlockHash] = handler.SendBlockHashSq
	pf.packetSqHandler[packets.PacketType_GetBlockPrevHash] = handler.GetBlockPrevHashSq
	pf.packetSqHandler[packets.PacketType_SendBlockPrevHash] = handler.SendBlockPrevHashSq
	pf.packetSqHandler[packets.PacketType_GetTxStatus] = handler.GetTxStatusSq
	pf.packetSqHandler[packets.PacketType_SendTxStatus] = handler.SendTxStatusSq

	pf.packetCqHandler[packets.PacketType_GetGhostNetVersion] = handler.GetGhostNetVersionCq
	pf.packetCqHandler[packets.PacketType_NotificationMasterNode] = handler.NotificationMasterNodeCq
	pf.packetCqHandler[packets.PacketType_ConnectToMasterNode] = handler.ConnectToMasterNodeCq
	pf.packetCqHandler[packets.PacketType_SearchGhostPubKey] = handler.SearchGhostPubKeyCq
	pf.packetCqHandler[packets.PacketType_RequestMasterNodeList] = handler.RequestMasterNodeListCq
	pf.packetCqHandler[packets.PacketType_ResponseMasterNodeList] = handler.ResponseMasterNodeListCq
	pf.packetCqHandler[packets.PacketType_SearchMasterPubKey] = handler.SearchMasterPubKeyCq

	pf.packetCqHandler[packets.PacketType_GetHeightestBlock] = handler.GetHeightestBlockCq
	pf.packetCqHandler[packets.PacketType_NewBlock] = handler.NewBlockCq
	pf.packetCqHandler[packets.PacketType_GetBlock] = handler.GetBlockCq
	pf.packetCqHandler[packets.PacketType_SendBlock] = handler.SendBlockCq
	pf.packetCqHandler[packets.PacketType_SendTransaction] = handler.SendTransactionCq
	pf.packetCqHandler[packets.PacketType_SendDataTransaction] = handler.SendDataTransactionCq
	pf.packetCqHandler[packets.PacketType_SearchDataTransaction] = handler.SearchDataTransactionCq
	pf.packetCqHandler[packets.PacketType_SendDataTxIdList] = handler.SendDataTxIdListCq
	pf.packetCqHandler[packets.PacketType_GetDataTxIdList] = handler.GetDataTxIdListCq
	pf.packetCqHandler[packets.PacketType_GetBlockHash] = handler.GetBlockHashCq
	pf.packetCqHandler[packets.PacketType_SendBlockHash] = handler.SendBlockHashCq
	pf.packetCqHandler[packets.PacketType_GetBlockPrevHash] = handler.GetBlockPrevHashCq
	pf.packetCqHandler[packets.PacketType_SendBlockPrevHash] = handler.SendBlockPrevHashCq
	pf.packetCqHandler[packets.PacketType_GetTxStatus] = handler.GetTxStatusCq
	pf.packetCqHandler[packets.PacketType_SendTxStatus] = handler.SendTxStatusCq

}
