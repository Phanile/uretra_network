package network

import (
	"fmt"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	RPCHandler RPCHandler
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
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	s := &Server{
		so:          opts,
		blockTime:   opts.BlockTime,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}),
	}

	if opts.RPCHandler == nil {
		opts.RPCHandler = NewDefaultRPCHandler(s)
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.blockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			err := s.so.RPCHandler.HandleRPC(rpc)

			if err != nil {
				fmt.Println(err)
			}

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

func (s *Server) ProcessTransaction(from NetAddress, transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	fmt.Printf("process transaction from %s - %s, memPool - %d\n", from, hash, s.memPool.Len())

	if s.memPool.Has(transaction.Hash(core.TxHasher{})) {
		return nil
	}

	if transaction.Verify() {
		s.memPool.Add(transaction)
	}

	transaction.SetFirstSeen(time.Now().UnixNano())

	return nil
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
