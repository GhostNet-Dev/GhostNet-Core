package dummyfactory

import (
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/common"
	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/workload"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type DummyFactory struct {
	master         *ptypes.GhostUser
	bootFactory    *factory.BootFactory
	defaultFactory *factory.DefaultFactory
	conn           *common.ConnectMaster

	Worker []*workload.Workload
}

func NewDummyFactory(maxWorker int, masterIp *ptypes.GhostIp, bootFactory *factory.BootFactory,
	networkFactory *factory.NetworkFactory, defaultFactory *factory.DefaultFactory,
	glog *glogger.GLogger) *DummyFactory {
	master := &ptypes.GhostUser{
		Ip: masterIp,
	}
	factory := &DummyFactory{
		master:         master,
		bootFactory:    bootFactory,
		defaultFactory: defaultFactory,
		conn: common.NewConnectMaster(store.DefaultMastersTable, bootFactory.Db,
			networkFactory.PacketFactory, networkFactory.Udp, defaultFactory.UserWallet),
	}

	for i := 0; i < maxWorker; i++ {
		factory.Worker = append(factory.Worker, workload.NewWorkload(fmt.Sprintf("worker%d", i), bootFactory.LoadWallet,
			defaultFactory.BlockServer, defaultFactory.BlockContainer, defaultFactory.Txs,
			factory.conn, defaultFactory.UserWallet, glog))
	}

	return factory
}
