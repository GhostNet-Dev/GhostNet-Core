package gnetwork

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

func (master *MasterNetwork) GetGhostNetVersionSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.VersionInfoSq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.VersionInfoCq{
		Master:  p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		Version: uint32(master.config.GhostVersion),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_GetGhostNetVersion,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) GetGhostNetVersionCq(packet []byte, from *net.UDPAddr) {
	cq := &packets.VersionInfoCq{}
	if err := proto.Unmarshal(packet, cq); err != nil {
		log.Fatal(err)
	}
	if cq.Version > uint32(master.config.GhostVersion) {
		// TODO: 새로운 Version을 받아야한다.
	}
}

func (master *MasterNetwork) NotificationMasterNodeSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}
	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User: &ptypes.GhostUser{
			PubKey:   master.owner.GetPubAddress(),
			Nickname: master.nickname,
		},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_GetGhostNetVersion,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) NotificationMasterNodeCq(packet []byte, from *net.UDPAddr) {}

func (master *MasterNetwork) ConnectToMasterNodeSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}

	master.AddMasterNode(&MasterNode{User: sq.User, NetAddr: from})

	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User: &ptypes.GhostUser{
			PubKey:   master.owner.GetPubAddress(),
			Nickname: master.nickname,
		},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_GetGhostNetVersion,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
} //ConnectToMasterNode = 3;

func (master *MasterNetwork) ConnectToMasterNodeCq(packet []byte, from *net.UDPAddr) {
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(packet, cq); err != nil {
		log.Fatal(err)
	}
	master.masterInfo = &MasterNode{User: cq.User, NetAddr: from}
} //ConnectToMasterNode = 3;

func (master *MasterNetwork) RequestMasterNodeListSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.RequestMasterNodeListSq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}

	// make response sq
	userList, totalItem := master.GetMasterNodeUserList(sq.StartIndex)

	resSq := packets.ResponseMasterNodeListSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   userList,
		Total:  totalItem,
	}
	responseData, err := proto.Marshal(&resSq)

	// make request cq
	cq := packets.RequestMasterNodeListCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_GetGhostNetVersion,
			PacketData: sendData,
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			PacketType: packets.PacketType_ResponseMasterNodeList,
			PacketData: responseData,
			SqFlag:     true,
		},
	}
} //RequestMasterNodeList = 5;

func (master *MasterNetwork) RequestMasterNodeListCq(packet []byte, from *net.UDPAddr) {} //RequestMasterNodeList = 5;

func (master *MasterNetwork) ResponseMasterNodeListSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}

	return nil
} //ResponseMasterNodeList = 6;

func (master *MasterNetwork) ResponseMasterNodeListCq(packet []byte, from *net.UDPAddr) {
} //ResponseMasterNodeList = 6;

func (master *MasterNetwork) SearchGhostPubKeySq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}
	// TODO: Db에서 찾아야하므로 별도의 nick table이 필요
	return nil
} //SearchGhostPubKey = 4;

func (master *MasterNetwork) SearchGhostPubKeyCq(packet []byte, from *net.UDPAddr) {
} //SearchGhostPubKey = 4;

func (master *MasterNetwork) SearchMasterPubKeySq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(packet, sq); err != nil {
		log.Fatal(err)
	}
	node := master.GetMasterNodeByNickname(sq.Nickname)

	cq := packets.SearchGhostPubKeyCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   []*ptypes.GhostUser{node.User},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_SearchMasterPubKey,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
} //SearchMasterPubKey = 9;

func (master *MasterNetwork) SearchMasterPubKeyCq(packet []byte, from *net.UDPAddr) {
} //SearchMasterPubKey = 9;
