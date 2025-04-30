package network

import (
	"github.com/Phanile/uretra_network/core"
	"time"
)

type BlocksMessage struct {
	Blocks []*core.Block
}

type GetBlocksMessage struct {
	From uint32
	To   uint32
}

type GetStatusMessage struct{}

type PingMessage struct {
	RequestTime time.Time
}

type PongMessage struct {
	BlockchainHeight uint32
	ResponseTime     time.Time
	RequestTime      time.Time
}

type StatusMessage struct {
	ID           string
	ActualHeight uint32
	Version      uint32
}
