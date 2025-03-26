package network

import (
	"bytes"
	"github.com/go-kit/log"
	"os"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
	"uretra-network/types"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	Transports    []Transport
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	so          *ServerOptions
	memPool     *TxPool
	chain       *core.Blockchain
	isValidator bool
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts *ServerOptions) (*Server, error) {
	if opts.BlockTime == time.Duration(0) {
		opts.BlockTime = defaultBlockTime
	}

	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain := core.NewBlockchain(genesisBlock())

	s := &Server{
		so:          opts,
		memPool:     NewTxPool(),
		chain:       chain,
		isValidator: opts.PrivateKey != nil,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}),
	}

	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	return s, nil
}

func (s *Server) Start() {
	s.initTransports()

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			msg, err := s.so.RPCDecodeFunc(rpc)

			if err != nil {
				_ = s.so.Logger.Log("error", "cannot decode rpc")
			}

			errMessage := s.so.RPCProcessor.ProcessMessage(msg)

			if errMessage != nil {
				_ = s.so.Logger.Log("error", "cannot process message from rpc channel")
			}

		case <-s.quitChannel:
			break free
		}
	}

	_ = s.so.Logger.Log("msg", "Server shutdown")
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

func (s *Server) ProcessMessage(m *DecodedMessage) error {
	switch data := m.Data.(type) {
	case *core.Transaction:
		return s.ProcessTransaction(data)
	}

	return nil
}

func (s *Server) ProcessTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	_ = s.so.Logger.Log("msg", "adding new tx to mempool", "hash", hash, "mempool length", s.memPool.Len())

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

func (s *Server) initTransports() {
	for _, transport := range s.so.Transports {
		go func(transport Transport) {
			for rpc := range transport.Consume() {
				s.rpcChannel <- rpc
			}
		}(transport)
	}
}

func (s *Server) validatorLoop() {
	ticker := time.NewTicker(s.so.BlockTime)

	_ = s.so.Logger.Log("msg", "starting validator loop")

	for {
		<-ticker.C
		_ = s.createNewBlock()
	}
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

func (s *Server) createNewBlock() error {
	header, err := s.chain.GetHeader(s.chain.Height())

	if err != nil {
		return err
	}

	block, e := core.NewBlockFromPrevHeader(header, nil)

	signErr := block.Sign(*s.so.PrivateKey)

	if signErr != nil {
		return signErr
	}

	if e != nil {
		return e
	}

	s.chain.AddBlock(block)

	return nil
}

func genesisBlock() *core.Block {
	h := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: time.Now().UnixNano(),
		Height:    0,
	}
	return core.NewBlock(h, nil)
}
