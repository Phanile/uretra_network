package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/go-kit/log"
	"net"
	"os"
	"sync"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
	"uretra-network/types"
)

var defaultBlockTime = 5 * time.Second

type ServerOptions struct {
	SeedNodes     []string
	ListenAddress string
	ID            string
	Logger        log.Logger
	RPCDecodeFunc RPCDecodeFunc
	RPCProcessor  RPCProcessor
	BlockTime     time.Duration
	PrivateKey    *crypto.PrivateKey
}

type Server struct {
	TCPTransport *TCPTransport
	peerCh       chan *TCPPeer
	mu           sync.RWMutex
	peerMap      map[net.Addr]*TCPPeer
	so           *ServerOptions
	memPool      *TxPool
	chain        *core.Blockchain
	isValidator  bool
	rpcChannel   chan RPC
	quitChannel  chan struct{}
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

	chain := core.NewBlockchain(opts.Logger, genesisBlock())

	peerCh := make(chan *TCPPeer)
	tr := NewTCPTransport(opts.ListenAddress, peerCh)

	s := &Server{
		TCPTransport: tr,
		peerCh:       peerCh,
		peerMap:      make(map[net.Addr]*TCPPeer),
		so:           opts,
		memPool:      NewTxPool(256),
		chain:        chain,
		isValidator:  opts.PrivateKey != nil,
		rpcChannel:   make(chan RPC),
		quitChannel:  make(chan struct{}, 1),
	}
	s.TCPTransport.peerCh = peerCh

	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	if s.isValidator {
		go s.validatorLoop()
	}

	s.boostrapPeers()

	return s, nil
}

func (s *Server) Start() {
	_ = s.TCPTransport.Start()
	s.boostrapPeers()

free:
	for {
		select {
		case peer := <-s.peerCh:
			_, ok := s.peerMap[peer.conn.RemoteAddr()]

			if !ok {
				go peer.readLoop(s.rpcChannel)
			}

		case rpc := <-s.rpcChannel:
			msg, err := s.so.RPCDecodeFunc(rpc)

			if err != nil {
				_ = s.so.Logger.Log("error", "cannot decode rpc")
				continue
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

func (s *Server) boostrapPeers() {
	for _, addr := range s.so.SeedNodes {
		go func(addr string) {
			dial, err := net.Dial("tcp", addr)

			if err != nil {
				return
			}

			s.peerCh <- &TCPPeer{
				conn: dial,
			}

		}(addr)
	}
}

func (s *Server) sendGetStatusMessage(peer *TCPPeer) {
	getStatusMsg := &GetStatusMessage{}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(getStatusMsg)

	if err != nil {
		_ = s.so.Logger.Log("error", "encode get status message failed")
		return
	}

	msg := NewMessage(MessageTypeGetStatus, buf.Bytes())
	msgData, _ := msg.Bytes()

	fmt.Println(msgData)
}

func (s *Server) Broadcast(payload []byte) error {
	for _, peer := range s.peerMap {
		err := peer.Send(payload)

		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) ProcessMessage(m *DecodedMessage) error {
	switch data := m.Data.(type) {
	case *core.Transaction:
		return s.processTransaction(data)
	case *core.Block:
		return s.processBlock(data)
	case *GetStatusMessage:
		return s.processGetStatusMessage(m.From)
	case *StatusMessage:
		return s.processStatusMessage(m.From, data)
	case *GetBlocksMessage:
		return s.processGetBlocksMessage(m.From, data)
	}
	return nil
}

func (s *Server) processTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	_ = s.so.Logger.Log("msg", "adding new tx to mempool", "hash", hash, "mempool length", s.memPool.PendingCount())

	if s.memPool.Contains(hash) {
		return nil
	}

	if transaction.Verify() {
		go s.broadcastTx(transaction)

		s.memPool.Add(transaction)
	}

	return nil
}

func (s *Server) processBlock(b *core.Block) error {
	if s.chain.AddBlock(b) {
		go s.broadcastBlock(b)
	}

	return nil
}

func (s *Server) processGetStatusMessage(from net.Addr) error {
	statusMessage := &StatusMessage{
		ActualHeight: s.chain.Height(),
		ID:           s.so.ID,
	}

	buf := &bytes.Buffer{}

	err := gob.NewEncoder(buf).Encode(statusMessage)

	if err != nil {
		return err
	}

	msg := NewMessage(MessageTypeStatus, buf.Bytes())
	data, _ := msg.Bytes()

	fmt.Println(data)

	return nil
}

func (s *Server) processStatusMessage(adr net.Addr, m *StatusMessage) error {
	fmt.Println("I AM ", s.so.ID, ", ", "get from ", adr, " data: ", " ID: ", m.ID, " Height of chain: ", m.ActualHeight)

	if s.chain.Height() >= m.ActualHeight {
		return nil
	}

	getBlocksMessage := &GetBlocksMessage{
		From: s.chain.Height(),
		To:   0,
	}

	buf := &bytes.Buffer{}

	_ = gob.NewEncoder(buf).Encode(getBlocksMessage)

	msg := NewMessage(MessageTypeGetBlocks, buf.Bytes())
	data, _ := msg.Bytes()

	fmt.Println("I AM ", s.so.ID, ", ", "wanna get blocks from ", adr, " : my chain height: ", getBlocksMessage.From, " To: ", getBlocksMessage.To)

	fmt.Println(data)

	return nil
}

func (s *Server) processGetBlocksMessage(from net.Addr, m *GetBlocksMessage) error {
	return nil
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
	encMsg, errBytes := msg.Bytes()

	if errBytes != nil {
		return errBytes
	}

	return s.Broadcast(encMsg)
}

func (s *Server) broadcastBlock(b *core.Block) error {
	buf := &bytes.Buffer{}
	err := b.Encode(core.NewGobBlockEncoder(buf))

	if err != nil {
		return err
	}

	msg := NewMessage(MessageTypeBlock, buf.Bytes())
	encMsg, errBytes := msg.Bytes()

	if errBytes != nil {
		return errBytes
	}

	return s.Broadcast(encMsg)
}

func (s *Server) createNewBlock() error {
	header, err := s.chain.GetHeader(s.chain.Height())

	if err != nil {
		return err
	}

	txs := s.memPool.Pending()

	block, e := core.NewBlockFromPrevHeader(header, txs)

	signErr := block.Sign(*s.so.PrivateKey)

	if signErr != nil {
		return signErr
	}

	if e != nil {
		return e
	}

	if s.chain.AddBlock(block) {
		go s.broadcastBlock(block)

		s.memPool.ClearPending()
	}

	return nil
}

func genesisBlock() *core.Block {
	h := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: 0,
		Height:    0,
	}
	return core.NewBlock(h, nil)
}
