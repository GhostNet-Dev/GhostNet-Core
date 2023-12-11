package ptypes

import (
	"net"
)

func (ghostIp *GhostIp) GetUdpAddr() *net.UDPAddr {
	to, _ := net.ResolveUDPAddr("udp", ghostIp.Ip+":"+ghostIp.Port)
	return to
}

func (user *GhostUser) Validate() bool {
	if user.PubKey == "" || user.Nickname == "" {
		return false
	}

	return true
}
