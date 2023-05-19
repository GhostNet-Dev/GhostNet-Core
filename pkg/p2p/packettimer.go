package p2p

import (
	"container/list"
	"sync"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
)

type TimerEntry struct {
	response  *ResponseHeaderInfo
	listEntry *list.Element
	time      int64
	retry     int
}

type PacketTimer struct {
	mutex        *sync.Mutex
	packetPool   map[string]*TimerEntry
	packetQ      *container.Queue
	retryHandler func(*ResponseHeaderInfo)
	loopSec      int
	timeout      int
}

func NewPacketTimer(loopSec int, timeout int, retryHandler func(*ResponseHeaderInfo)) *PacketTimer {
	timer := &PacketTimer{
		mutex:        &sync.Mutex{},
		packetPool:   make(map[string]*TimerEntry),
		packetQ:      container.NewQueue(),
		retryHandler: retryHandler,
		loopSec:      loopSec,
		timeout:      timeout,
	}

	go func() {
		for {
			if loopSec > 0 {
				time.Sleep(time.Second * time.Duration(loopSec))
			}
			if len(timer.packetPool) > 0 {
				limit := time.Now().Unix() - int64(timeout)
				timer.mutex.Lock()
				req := timer.packetQ.Peek().(*TimerEntry)
				timer.mutex.Unlock()
				if req.time < limit && req.retry < 3 {
					timer.RetryPacket(req)
					timer.mutex.Lock()
					timer.packetQ.Push(timer.packetQ.Pop())
					timer.mutex.Unlock()
				} else if req.retry > 3 {
					timer.ReleaseSqPacket(req.response.RequestId)
					if req.response.Callback != nil {
						req.response.Callback(false)
					}
				}
			}
		}
	}()
	return timer
}

func (packetTimer *PacketTimer) RetryPacket(req *TimerEntry) {
	req.retry++
	packetTimer.retryHandler(req.response)
}

func (packetTimer *PacketTimer) RegisterSqPacket(header *ResponseHeaderInfo) {
	packetTimer.mutex.Lock()
	defer packetTimer.mutex.Unlock()
	if _, exist := packetTimer.packetPool[string(header.RequestId)]; exist {
		return
	}
	entry := &TimerEntry{
		response: header,
		time:     time.Now().Unix(),
		retry:    0,
	}

	entry.listEntry = packetTimer.packetQ.Push(entry)
	packetTimer.packetPool[string(header.RequestId)] = entry
}

func (packetTimer *PacketTimer) ReleaseSqPacket(requestId []byte) {
	packetTimer.mutex.Lock()
	req := packetTimer.packetPool[string(requestId)]
	if req != nil {
		packetTimer.packetQ.Remove(req.listEntry)
		delete(packetTimer.packetPool, string(requestId))
	}
	packetTimer.mutex.Unlock()

	if req != nil && req.response.Callback != nil {
		req.response.Callback(true)
	}
}
