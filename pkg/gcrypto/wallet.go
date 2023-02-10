package gcrypto

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"

type Wallet struct {
	masterNode *ptypes.GhostUser
	ghostIp    *ptypes.GhostIp
	myAddr     *GhostAddress
	nickname   string
}

func NewWallet(nickname string, myAddr *GhostAddress, ghostIp *ptypes.GhostIp, master *ptypes.GhostUser) *Wallet {
	return &Wallet{
		masterNode: master,
		nickname:   nickname,
		myAddr:     myAddr,
		ghostIp:    ghostIp,
	}
}

func (w *Wallet) SetMasterNode(master *ptypes.GhostUser) {
	w.masterNode = master
}

func (w *Wallet) GetGhostUser() *ptypes.GhostUser {
	return &ptypes.GhostUser{
		Nickname: w.nickname,
		PubKey:   w.GetPubAddress(),
		Ip:       w.ghostIp,
	}
}

func (w *Wallet) GetGhostAddress() *GhostAddress {
	return w.myAddr
}

func (w *Wallet) MyPubKey() []byte {
	return w.myAddr.Get160PubKey()
}

func (w *Wallet) GetPubAddress() string {
	return w.myAddr.GetPubAddress()
}

func (w *Wallet) GetNickname() string {
	return w.nickname
}
