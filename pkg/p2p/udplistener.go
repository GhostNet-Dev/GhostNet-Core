package p2p

import (
	"fmt"
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

type UdpServer struct {
	UdpConn *net.UDPConn
	Ip      string
	Port    string
	Pf      *PacketFactory
}

type RequestPacketInfo struct {
	Addr       *net.UDPAddr
	Err        error
	PacketByte []byte
}

type ResponsePacketInfo struct {
	ToAddr     *net.UDPAddr
	PacketType packets.PacketType
	PacketData []byte
	SqFlag     bool
}

func NewUdpServer(ip string, port string) *UdpServer {
	return &UdpServer{
		Ip:   ip,
		Port: port,
		Pf:   NewPacketFactory(),
	}
}

func (udp *UdpServer) Start(netChannel chan RequestPacketInfo) {
	service := udp.Ip + ":" + udp.Port

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	if udp.UdpConn, err = net.ListenUDP("udp", udpAddr); err != nil {
		defer udp.UdpConn.Close()
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on addr= ", service)

	if netChannel == nil {
		netChannel = make(chan RequestPacketInfo)
		go func() {
			for {
				select {
				case packetInfo := <-netChannel:
					packetByte := packetInfo.PacketByte
					recvPacket := packets.Any{}
					if err := proto.Unmarshal(packetByte, &recvPacket); err != nil {
						// packet type별로 callback handler를 만들어야한다.
						log.Fatal(err)
					}

					if recvPacket.SqFlag == true {
						if response := udp.Pf.packetSqHandler[recvPacket.Type](recvPacket.PacketData, packetInfo.Addr); response != nil {
							for _, packet := range response {
								udp.SendPacket(&packet)
							}
						}
					} else {
						udp.Pf.packetCqHandler[recvPacket.Type](recvPacket.PacketData, packetInfo.Addr)
					}
				}
			}
		}()
	}

	buffer := make([]byte, 1024)
	go func() {
		defer udp.UdpConn.Close()

		for {
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

func (udp *UdpServer) SendPacket(sendInfo *ResponsePacketInfo) {
	anyData := packets.Any{
		Type:       sendInfo.PacketType,
		SqFlag:     sendInfo.SqFlag,
		PacketData: sendInfo.PacketData,
	}
	sendData, err := proto.Marshal(&anyData)
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
