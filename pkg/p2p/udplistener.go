package p2p

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

type UdpServer struct {
	UdpConn   *net.UDPConn
	Pf        *PacketFactory
	glog      *glogger.GLogger
	Ip        string
	Port      string
	StartFlag bool
}

type RequestPacketInfo struct {
	Addr       *net.UDPAddr
	Err        error
	PacketByte []byte
}

type RequestHeaderInfo struct {
	FromAddr *net.UDPAddr
	Header   *packets.Header
}

type ResponseHeaderInfo struct {
	ToAddr     *net.UDPAddr
	PacketType packets.PacketType
	SecondType packets.PacketSecondType
	ThirdType  packets.PacketThirdType
	PacketData []byte
	SqFlag     bool
}

func NewUdpServer(ip, port string, packetFactory *PacketFactory, glog *glogger.GLogger) *UdpServer {
	return &UdpServer{
		Ip:        ip,
		Port:      port,
		Pf:        packetFactory,
		glog:      glog,
		StartFlag: false,
	}
}

func (udp *UdpServer) GetLocalIp() *ptypes.GhostIp {
	return &ptypes.GhostIp{Ip: udp.Ip, Port: udp.Port}
}

func (udp *UdpServer) loadIp() *ptypes.GhostIp {
	url := "https://api.ipify.org?format=text"
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	ip, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	localIp := string(ip)
	return &ptypes.GhostIp{Ip: localIp, Port: udp.Port}
}

func (udp *UdpServer) Start(netChannel chan RequestPacketInfo, ip, port string) {
	if udp.StartFlag {
		return
	}
	udp.StartFlag = true

	if ip != "" {
		udp.Ip = ip
	}
	if port != "" {
		udp.Port = port
	}

	service := udp.Ip + ":" + udp.Port

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	if udp.UdpConn, err = net.ListenUDP("udp4", udpAddr); err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on addr= ", service)

	if netChannel == nil {
		netChannel = make(chan RequestPacketInfo)
		go func() {
			for packetInfo := range netChannel {
				packetByte := packetInfo.PacketByte
				recvPacket := packets.Header{}
				if err := proto.Unmarshal(packetByte, &recvPacket); err != nil {
					// packet type별로 callback handler를 만들어야한다.
					log.Fatal(err)
				}
				firstLevel, exist := udp.Pf.firstLevel[recvPacket.Type]
				if !exist {
					return
				}
				udp.glog.DebugOutput(udp,
					fmt.Sprint("Recv from ", packetInfo.Addr, ", ",
						packets.PacketSecondType_name[int32(recvPacket.SecondType)],
						" => ", packets.PacketThirdType_name[int32(recvPacket.ThirdType)], " SQ: ",
						recvPacket.SqFlag), glogger.PacketLog)

				// TODO: it need to refac
				// sq
				if recvPacket.SqFlag {
					go func(packetInfo *RequestPacketInfo) {
						secondLevel, exist := firstLevel.packetSqHandler[recvPacket.SecondType]
						if !exist {
							return
						}
						if response := secondLevel(&RequestHeaderInfo{FromAddr: packetInfo.Addr, Header: &recvPacket}); response != nil {
							for _, packet := range response {
								packet.PacketType = recvPacket.Type
								udp.SendResponse(&packet)
							}
						}
					}(&packetInfo)
				} else {
					go func(packetInfo *RequestPacketInfo) {
						// cq
						secondLevel, exist := firstLevel.packetCqHandler[recvPacket.SecondType]
						if !exist {
							return
						}
						if response := secondLevel(&RequestHeaderInfo{FromAddr: packetInfo.Addr, Header: &recvPacket}); response != nil {
							for _, packet := range response {
								packet.PacketType = recvPacket.Type
								udp.SendResponse(&packet)
							}
						}
					}(&packetInfo)
				}
			}
		}()
	}

	go func() {
		for {
			buffer := make([]byte, 1024)
			n, addr /*n, addr*/, err := udp.UdpConn.ReadFromUDP(buffer)
			if err != nil {
				//doneChan <-err
				continue
			}

			netChannel <- RequestPacketInfo{
				Addr:       addr,
				Err:        err,
				PacketByte: buffer[:n],
			}
			/*
				err = pc.SetWriteDeadline(deadline)
				if err != nil {
					doneChan <-err
					return
				}

				n, err = pc.WriteTo(buffer[:n], addr)
				if err != nil {
					doneChan <-err
					return
				}
			*/
		}
	}()
}

func (udp *UdpServer) TranslationToHeader(sendInfo *ResponseHeaderInfo) *packets.Header {
	return &packets.Header{
		Type:       sendInfo.PacketType,
		SecondType: sendInfo.SecondType,
		ThirdType:  sendInfo.ThirdType,
		SqFlag:     sendInfo.SqFlag,
		PacketData: sendInfo.PacketData,
		Source:     udp.GetLocalIp(),
	}
}

func (udp *UdpServer) SendUdpPacket(sendInfo *ResponseHeaderInfo, to *net.UDPAddr) {
	anyData := udp.TranslationToHeader(sendInfo)
	sendData, err := proto.Marshal(anyData)
	if err != nil {
		log.Fatal(err)
	}
	udp.RawSendPacket(to, sendData)
	udp.glog.DebugOutput(udp,
		fmt.Sprint("Send to ", to, ", ", packets.PacketSecondType_name[int32(sendInfo.SecondType)],
			" => ", packets.PacketThirdType_name[int32(sendInfo.ThirdType)], " SQ: ",
			sendInfo.SqFlag), glogger.PacketLog)
}

func (udp *UdpServer) SendPacket(sendInfo *ResponseHeaderInfo, ipAddr *ptypes.GhostIp) {
	to, _ := net.ResolveUDPAddr("udp", ipAddr.Ip+":"+ipAddr.Port)
	udp.SendUdpPacket(sendInfo, to)
}

func (udp *UdpServer) SendResponse(sendInfo *ResponseHeaderInfo) {
	udp.glog.DebugOutput(udp,
		fmt.Sprint("Response to ", sendInfo.ToAddr, ", ", packets.PacketSecondType_name[int32(sendInfo.SecondType)],
			" => ", packets.PacketThirdType_name[int32(sendInfo.ThirdType)], " SQ: ",
			sendInfo.SqFlag), glogger.PacketLog)

	anyData := udp.TranslationToHeader(sendInfo)
	sendData, err := proto.Marshal(anyData)
	if err != nil {
		log.Fatal(err)
	}
	udp.RawSendPacket(sendInfo.ToAddr, sendData)
}

func (udp *UdpServer) RawSendPacket(addr *net.UDPAddr, buf []byte) {
	_, err := udp.UdpConn.WriteToUDP(buf, addr)
	if err != nil {
		log.Fatal(err)
	}
}
