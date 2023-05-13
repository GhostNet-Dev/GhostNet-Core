package common

import (
	"errors"
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"google.golang.org/protobuf/proto"
)

type ConnectMaster struct {
	db           *store.LiteStore
	udp          *p2p.UdpServer
	wallet       *gcrypto.Wallet
	table        string
	eventChannel chan bool
	eventWait    bool
}

const RootUrl = "www.ghostnetroot.com"

func NewConnectMaster(table string, db *store.LiteStore, packetFactory *p2p.PacketFactory,
	udp *p2p.UdpServer, w *gcrypto.Wallet) *ConnectMaster {

	conn := &ConnectMaster{
		db:           db,
		udp:          udp,
		wallet:       w,
		table:        table,
		eventChannel: make(chan bool),
		eventWait:    false,
	}
	conn.RegisterPacketHandler(packetFactory)
	return conn
}

func (conn *ConnectMaster) CheckNickname(nickname string) (result bool, err error) {
	sq := &packets.SearchGhostPubKeySq{
		Master:   p2p.MakeMasterPacket(conn.wallet.GetPubAddress(), 0, 0, conn.udp.GetLocalIp()),
		Nickname: nickname,
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     conn.wallet.GetMasterNode().Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_SearchGhostPubKey,
		PacketData: sendData,
		SqFlag:     true,
	}
	conn.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)

	select {
	case result = <-conn.eventChannel:
	case <-time.After(time.Second * 8):
		return false, errors.New("timeout")
	}

	return result, nil
}

func (conn *ConnectMaster) SendTx(tx *types.GhostTransaction) (result bool, err error) {
	sq := &packets.SendTransactionSq{
		Master: p2p.MakeMasterPacket(conn.wallet.GetPubAddress(), 0, 0, conn.udp.GetLocalIp()),
		TxId:   tx.TxId,
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     conn.wallet.GetMasterNode().Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  packets.PacketThirdType_SendTransaction,
		PacketData: sendData,
		SqFlag:     true,
	}
	conn.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)

	select {
	case result = <-conn.eventChannel:
	case <-time.After(time.Second * 8):
		return false, errors.New("timeout")
	}

	return result, nil
}

func (conn *ConnectMaster) RegisterPacketHandler(packetFactory *p2p.PacketFactory) {
	packetSqHandler := make(map[packets.PacketSecondType]p2p.FuncPacketHandler)
	packetCqHandler := make(map[packets.PacketSecondType]p2p.FuncPacketHandler)

	packetCqHandler[packets.PacketSecondType_ConnectToMasterNode] = func(rhi *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
		return nil
	}
	packetCqHandler[packets.PacketSecondType_SearchGhostPubKey] = func(rhi *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
		header := rhi.Header
		cq := &packets.SearchGhostPubKeyCq{}
		if err := proto.Unmarshal(header.PacketData, cq); err != nil {
			log.Fatal(err)
		}

		conn.eventChannel <- (len(cq.User) != 0)
		return nil
	}

	packetFactory.UpdatePacketHandler(packets.PacketType_MasterNetwork, packetSqHandler, packetCqHandler)
}
