package node

import idApi "github.com/zodimo/go-compose/pkg/compose-identifier/api"

func NewNodeID() NodeID {
	return idApi.NewIdentifier()
}
