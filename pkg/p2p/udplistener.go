package p2p

import (
	"fmt"
	"log"
	"net"
)

type UdpServer struct {
	UdpConn *net.UDPConn
	Ip      string
	Port    string
}

func NewUdpServer(ip string, port string) *UdpServer {
	return &UdpServer{
		Ip:   ip,
		Port: port,
	}
}

func (udp *UdpServer) Start() {
	service := udp.Ip + ":" + udp.Port

	udpAddr, err := net.ResolveUDPAddr("udp4", service)

	if err != nil {
		log.Fatal(err)
	}

	// setup listener for incoming UDP connection
	if udp.UdpConn, err = net.ListenUDP("udp", udpAddr); err != nil {
		log.Fatal(err)
	}

	fmt.Println("UDP server up and listening on addr= ", service)

	buffer := make([]byte, 1024)
	go func() {
		defer udp.UdpConn.Close()

		for {
			_, _ /*n, addr*/, err := udp.UdpConn.ReadFrom(buffer)
			if err != nil {
				//doneChan <-err
				continue
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
