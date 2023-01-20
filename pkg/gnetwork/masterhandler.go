package gnetwork

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

func (master *MasterNetwork) GetGhostNetVersionSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.VersionInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.VersionInfoCq{
		Master:  p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		Version: master.config.GhostVersion,
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_GetGhostNetVersion,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) GetGhostNetVersionCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	cq := &packets.VersionInfoCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	if cq.Version > uint32(master.config.GhostVersion) {
		// TODO: 새로운 Version을 받아야한다.
	}
	return nil
}

func (master *MasterNetwork) NotificationMasterNodeSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   master.getGhostUser(),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_NotificationMasterNode,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) NotificationMasterNodeCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	return nil
}

func (master *MasterNetwork) ConnectToMasterNodeSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	master.account.AddUserNode(&GhostNode{User: sq.User, NetAddr: header.Source.GetUdpAddr()})

	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   master.getGhostUser(),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ConnectToMasterNode,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) ConnectToMasterNodeCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	master.masterInfo = &GhostNode{User: cq.User, NetAddr: header.Source.GetUdpAddr()}
	return nil
}

func (master *MasterNetwork) RequestMasterNodeListSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.RequestMasterNodeListSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	// make response sq
	userList, totalItem := master.account.GetMasterNodeUserList(sq.StartIndex)

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

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_RequestMasterNodeList,
			PacketData: sendData,
			SqFlag:     false,
		},
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ResponseMasterNodeList,
			PacketData: responseData,
			SqFlag:     true,
		},
	}
}

func (master *MasterNetwork) RequestMasterNodeListCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	return nil
}

func (master *MasterNetwork) ResponseMasterNodeListSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	master.account.AddMasterUserList(sq.User)

	cq := packets.ResponseMasterNodeListCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ResponseMasterNodeList,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) ResponseMasterNodeListCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	return nil
}

func (master *MasterNetwork) SearchGhostPubKeySq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	// TODO: Db에서 찾아야하므로 별도의 nick table이 필요
	node := master.account.GetNodeByNickname(sq.Nickname)

	cq := packets.SearchGhostPubKeyCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   []*ptypes.GhostUser{node.User},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_SearchGhostPubKey,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) SearchGhostPubKeyCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (master *MasterNetwork) SearchMasterPubKeySq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	node := master.account.GetNodeInfo(sq.PubKey)

	cq := packets.SearchGhostPubKeyCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User:   []*ptypes.GhostUser{node.User},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_SearchUserInfoByPubKey,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) SearchMasterPubKeyCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	// TODO:
	return nil
}

func (master *MasterNetwork) RegisterBlockHandler(handlerSq func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo,
	handlerCq func(*packets.Header, *net.UDPAddr)) {
	master.blockHandlerSq = handlerSq
	master.blockHandlerCq = handlerCq
}

func (master *MasterNetwork) BlockChainSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	if master.blockHandlerSq != nil {
		infos := master.blockHandlerSq(header, routingInfo.SourceIp)
		for _, info := range infos {
			info.PacketType = packets.PacketType_MasterNetwork
			info.SecondType = packets.PacketSecondType_BlockChain
		}
		return infos
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_BlockChain,
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) BlockChainCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	if master.blockHandlerCq != nil {
		master.blockHandlerCq(header, routingInfo.SourceIp)
	}
	return nil
}

func (master *MasterNetwork) ForwardingSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.ForwardingSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	routingInfo.Level = int(sq.Master.Level)
	routingInfo.RoutingType = sq.Master.RoutingT
	//forwardingFrom, _ := net.ResolveUDPAddr("udp", sq.ForwardingHeader.Source.Ip+":"+sq.ForwardingHeader.Source.Port)
	headerInfo := master.packetSqHandler[sq.ForwardingHeader.SecondType](sq.ForwardingHeader, routingInfo)
	routingHeaderInfo := master.makeForwadingPacket(sq.Master.RoutingT, sq.Master.Level+1, header)

	return append(append([]p2p.PacketHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_Forwarding,
			SqFlag:     false,
		},
	}, headerInfo...), routingHeaderInfo...)
}

func (master *MasterNetwork) ForwardingCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	return nil
}
