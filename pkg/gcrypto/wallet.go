package gcrypto

type Wallet struct {
	myAddr   GhostAddress
	nickname string
}

func NewWallet() *Wallet {
	return new(Wallet)
}

func (w *Wallet) MyPubKey() []byte {
	return w.myAddr.Get160PubKey()
}
