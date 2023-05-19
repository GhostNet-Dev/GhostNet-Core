package gnetwork

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

func (master *MasterNetwork) GetGhostNetVersionSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.VersionInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.VersionInfoCq{
		Master:  p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
		Version: master.ghostNetVersion,
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_GetGhostNetVersion,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) GetGhostNetVersionCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	cq := &packets.VersionInfoCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	if cq.Version > uint32(master.ghostNetVersion) {
		// TODO: 새로운 Version을 받아야한다.
	}
	return nil
}

func (master *MasterNetwork) NotificationMasterNodeSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	master.RegisterMasterNode(sq.User)

	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
		User:   master.getGhostUser(),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_NotificationMasterNode,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) NotificationMasterNodeCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}

func (master *MasterNetwork) ConnectToMasterNodeSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.MasterNodeUserInfoSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	sq.User.Ip = header.Source
	master.account.AddUserNode(sq.User)

	cq := packets.MasterNodeUserInfoCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
		User:   master.getGhostUser(),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ConnectToMasterNode,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) ConnectToMasterNodeCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	master.RegisterMyMasterNode(cq.User)
	return nil
}

func (master *MasterNetwork) RequestMasterNodeListSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.RequestMasterNodeListSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	// make response sq
	userList, totalItem := master.account.GetMasterNodeUserList(sq.StartIndex)

	resSq := packets.ResponseMasterNodeListSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		User:   userList,
		Total:  totalItem,
	}
	responseData, err := proto.Marshal(&resSq)
	if err != nil {
		log.Fatal(err)
	}
	// make request cq
	cq := packets.RequestMasterNodeListCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_RequestMasterNodeList,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ResponseMasterNodeList,
			RequestId:  resSq.Master.GetRequestId(),
			PacketData: responseData,
			SqFlag:     true,
		},
	}
}

func (master *MasterNetwork) RequestMasterNodeListCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}

func (master *MasterNetwork) ResponseMasterNodeListSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	master.account.AddMasterUserList(sq.User)

	cq := packets.ResponseMasterNodeListCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ResponseMasterNodeList,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) ResponseMasterNodeListCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}

func (master *MasterNetwork) SearchGhostPubKeySq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	// TODO: Db에서 찾아야하므로 별도의 nick table이 필요
	node := master.account.GetNodeByNickname(sq.Nickname)

	cq := packets.SearchGhostPubKeyCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
	}

	if node != nil {
		cq.User = []*ptypes.GhostUser{node}
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_SearchGhostPubKey,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) SearchGhostPubKeyCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (master *MasterNetwork) SearchMasterPubKeySq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.SearchGhostPubKeySq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}
	node := master.account.GetNodeInfo(sq.PubKey)

	if node == nil {
		return nil
	}

	cq := packets.SearchGhostPubKeyCq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, master.localGhostIp),
		User:   []*ptypes.GhostUser{node},
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_SearchUserInfoByPubKey,
			PacketData: sendData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) SearchMasterPubKeyCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	// TODO:
	return nil
}

func (master *MasterNetwork) RegisterBlockHandler(handlerSq func(*packets.Header, *net.UDPAddr) []p2p.ResponseHeaderInfo,
	handlerCq func(*packets.Header, *net.UDPAddr)) {
	master.blockHandlerSq = handlerSq
	master.blockHandlerCq = handlerCq
}

func (master *MasterNetwork) BlockChainSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	if master.blockHandlerSq != nil {
		infos := master.blockHandlerSq(header, header.Source.GetUdpAddr())
		for idx := range infos {
			infos[idx].PacketType = packets.PacketType_MasterNetwork
			infos[idx].SecondType = packets.PacketSecondType_BlockChain
		}
		return infos
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_BlockChain,
			RequestId:  header.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (master *MasterNetwork) BlockChainCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	if master.blockHandlerCq != nil {
		master.blockHandlerCq(header, header.Source.GetUdpAddr())
	}
	return nil
}

func (master *MasterNetwork) ForwardingSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.ForwardingSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	//forwardingFrom, _ := net.ResolveUDPAddr("udp", sq.ForwardingHeader.Source.Ip+":"+sq.ForwardingHeader.Source.Port)
	headerInfo := master.packetSqHandler[sq.ForwardingHeader.SecondType](&p2p.RequestHeaderInfo{
		FromAddr: sq.ForwardingHeader.Source.GetUdpAddr(),
		Header:   sq.ForwardingHeader,
	})
	routingHeaderInfo := master.makeForwadingPacket(sq.Master.RoutingT, sq.Master.Level+1, header)

	return append(append([]p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_Forwarding,
			RequestId:  sq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}, headerInfo...), routingHeaderInfo...)
}

func (master *MasterNetwork) ForwardingCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}
