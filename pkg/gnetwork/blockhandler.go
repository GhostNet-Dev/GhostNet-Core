package gnetwork

import (
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
)

// // blockchain packet type
func (master *MasterNetwork) GetHeightestBlockSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	return nil
} //GetHeightestBlock = 12;
func (master *MasterNetwork) GetHeightestBlockCq(packet []byte, from *net.UDPAddr) {

} //GetHeightestBlock = 12;
func (master *MasterNetwork) NewBlockSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //NewBlock = 13;
func (master *MasterNetwork) NewBlockCq(packet []byte, from *net.UDPAddr) {

} //NewBlock = 13;
func (master *MasterNetwork) GetBlockSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //GetBlock = 14;
func (master *MasterNetwork) GetBlockCq(packet []byte, from *net.UDPAddr) {

} //GetBlock = 14;
func (master *MasterNetwork) SendBlockSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendBlock = 15;
func (master *MasterNetwork) SendBlockCq(packet []byte, from *net.UDPAddr) {

} //SendBlock = 15;
// ScanAddrBlock = 16; // not used...
func (master *MasterNetwork) SendTransactionSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendTransaction = 17;
func (master *MasterNetwork) SendTransactionCq(packet []byte, from *net.UDPAddr) {

} //SendTransaction = 17;
func (master *MasterNetwork) SearchTransactionSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SearchTransaction = 18;
func (master *MasterNetwork) SearchTransactionCq(packet []byte, from *net.UDPAddr) {

} //SearchTransaction = 18;
func (master *MasterNetwork) SendDataTransactionSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendDataTransaction = 19;
func (master *MasterNetwork) SendDataTransactionCq(packet []byte, from *net.UDPAddr) {

} //SendDataTransaction = 19;
func (master *MasterNetwork) SearchDataTransactionSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SearchDataTransaction = 20;
func (master *MasterNetwork) SearchDataTransactionCq(packet []byte, from *net.UDPAddr) {

} //SearchDataTransaction = 20;
// ScanBlockChain = 21;
// CheckGhostNickname = 22;
func (master *MasterNetwork) SendDataTxIdListSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendDataTxIdList = 23;
func (master *MasterNetwork) SendDataTxIdListCq(packet []byte, from *net.UDPAddr) {

} //SendDataTxIdList = 23;
func (master *MasterNetwork) GetDataTxIdListSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //GetDataTxIdList = 24;
func (master *MasterNetwork) GetDataTxIdListCq(packet []byte, from *net.UDPAddr) {

} //GetDataTxIdList = 24;
// ReportBlockError = 25;
func (master *MasterNetwork) GetBlockHashSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //GetBlockHash = 26;
func (master *MasterNetwork) GetBlockHashCq(packet []byte, from *net.UDPAddr) {

} //GetBlockHash = 26;
func (master *MasterNetwork) SendBlockHashSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendBlockHash = 27;
func (master *MasterNetwork) SendBlockHashCq(packet []byte, from *net.UDPAddr) {

} //SendBlockHash = 27;
func (master *MasterNetwork) GetBlockPrevHashSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //GetBlockPrevHash = 28;
func (master *MasterNetwork) GetBlockPrevHashCq(packet []byte, from *net.UDPAddr) {

} //GetBlockPrevHash = 28;
func (master *MasterNetwork) SendBlockPrevHashSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendBlockPrevHash = 29;
func (master *MasterNetwork) SendBlockPrevHashCq(packet []byte, from *net.UDPAddr) {

} //SendBlockPrevHash = 29;
func (master *MasterNetwork) GetTxStatusSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //GetTxStatus = 30;
func (master *MasterNetwork) GetTxStatusCq(packet []byte, from *net.UDPAddr) {

} //GetTxStatus = 30;
func (master *MasterNetwork) SendTxStatusSq(packet []byte, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return nil
} //SendTxStatus = 31;
func (master *MasterNetwork) SendTxStatusCq(packet []byte, from *net.UDPAddr) {

} //SendTxStatus = 31;
