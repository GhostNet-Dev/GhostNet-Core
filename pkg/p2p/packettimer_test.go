package p2p

import (
	"testing"

	"github.com/rs/xid"
)

func TestTimer(t *testing.T) {
	ret := true
	testPacketTimer := NewPacketTimer(0, 0, func(rhi *ResponseHeaderInfo) {
		ret = false
	})
	requestId := xid.New().Bytes()
	testPacketTimer.RegisterSqPacket(&ResponseHeaderInfo{
		RequestId: requestId,
	})
	for ret {}
}
