package crypto

type Wallet struct {
	myAddr GhostAddress
}

func NewWallet(path string) *Wallet {
	return new(Wallet)
}

func (w *Wallet) MyPubKey() []byte {
	return w.myAddr.PubKey
}
