package packets

func (m *MasterPacket) GetRequestId() []byte {
	return m.Common.GetRequestId()
}

func (m *MasterPacket) GetFromPubkey() string {
	return m.Common.FromPubKeyAddress
}

func (m *MasterPacket) GetTimeId() uint64 {
	return m.Common.TimeId
}
