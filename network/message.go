package network

import "uretra-network/core"

type BlocksMessage struct {
	blocks []*core.Block
}

type GetBlocksMessage struct {
	From uint32
	To   uint32
}

type GetStatusMessage struct {
}

type StatusMessage struct {
	ID           string
	ActualHeight uint32
	Version      uint32
}
