package network

import (
	"fmt"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
)

type ServerOptions struct {
	Transports []Transport
	BlockTime  time.Duration
	PrivateKey *crypto.PrivateKey
}

type Server struct {
	so          *ServerOptions
	blockTime   time.Duration
	memPool     *TxPool
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts *ServerOptions) *Server {
	return &Server{
		so:          opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChannel:
			break free
		case <-ticker.C:
			if s.isValidator {
				s.createNewBlock()
			}
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) handleTransaction(t *core.Transaction) {
	if s.memPool.Has(t.Hash(core.TxHasher{})) {
		return
	}

	if t.Verify() {
		s.memPool.Add(t)
	}
}

func (s *Server) createNewBlock() error {
	fmt.Println("creating a new block")
	return nil
}

func (s *Server) initTransports() {
	for _, transport := range s.so.Transports {
		go func(transport Transport) {
			for rpc := range transport.Consume() {
				s.rpcChannel <- rpc
			}
		}(transport)
	}
}
