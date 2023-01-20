package sqlite

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"

type GhostAccount struct {
}

func (account *GhostAccount) GetMasterNodeList() []*ptypes.GhostUser {
	return nil
}

func (account *GhostAccount) GetMasterNodeSearch(pubKey string) []*ptypes.GhostUser {
	return nil
}

func (account *GhostAccount) GetMasterNodeSearchPick(pubKey string) *ptypes.GhostUser {
	return nil
}
