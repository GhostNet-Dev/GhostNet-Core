package ptypes

import "net"

func (ghostIp *GhostIp) GetUdpAddr() *net.UDPAddr {
	to, _ := net.ResolveUDPAddr("udp", ghostIp.Ip+":"+ghostIp.Port)
	return to
}
