package store

import (	
	gsql "github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
)

type BlockContainer struct {
	GSql gsql.GSql
}


func NewBlockContainer() *BlockContainer{
	return &BlockContainer{
		GSql : *gsql.NewGSql("sqlite"),
	}
}