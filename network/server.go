package network

import (
	"bytes"
	"fmt"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	so          *ServerOptions
	memPool     *TxPool
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts *ServerOptions) *Server {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	s := &Server{
		so:          opts,
		memPool:     NewTxPool(),
		isValidator: opts.PrivateKey != nil,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}),
	}

	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	return s
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(s.so.BlockTime)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			msg, err := s.so.RPCDecodeFunc(rpc)

			if err != nil {
				_ = fmt.Errorf("cannot decode rpc")
			}

			errMessage := s.so.RPCProcessor.ProcessMessage(msg)

			if errMessage != nil {
				_ = fmt.Errorf("cannot process message from rpc channel")
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

func (s *Server) Broadcast(payload []byte) error {
	for _, transport := range s.so.Transports {
		err := transport.Broadcast(payload)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) broadcastTx(tx *core.Transaction) error {
	buf := &bytes.Buffer{}
	err := tx.Encode(core.NewGobTxEncoder(buf))

	if err != nil {
		return err
	}

	msg := NewMessage(MessageTypeTx, buf.Bytes())
	bytes, errBytes := msg.Bytes()

	if errBytes != nil {
		return errBytes
	}

	return s.Broadcast(bytes)
}

func (s *Server) ProcessMessage(m *DecodedMessage) error {
	switch data := m.Data.(type) {
	case *core.Transaction:
		return s.ProcessTransaction(data)
	}

	return nil
}

func (s *Server) ProcessTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	fmt.Printf("process transaction %s, memPool - %d\n", hash, s.memPool.Len())

	if s.memPool.Has(transaction.Hash(core.TxHasher{})) {
		return nil
	}

	if transaction.Verify() {
		transaction.SetFirstSeen(time.Now().UnixNano())

		go func() {
			_ = s.broadcastTx(transaction)
		}()

		s.memPool.Add(transaction)
	}

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
