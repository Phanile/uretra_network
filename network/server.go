package network

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"github.com/Phanile/uretra_network/api"
	"github.com/Phanile/uretra_network/core"
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
	"github.com/go-kit/log"
	"net"
	"os"
	"sync"
	"time"
)

const (
	defaultListenPort    = ":3228"
	defaultAPIListenPort = ":3229"
)

type ServerOptions struct {
	SeedNodes        []string
	ListenAddress    string
	APIListenAddress string
	ID               string
	Logger           log.Logger
	RPCDecodeFunc    RPCDecodeFunc
	RPCProcessor     RPCProcessor
	PrivateKey       *crypto.PrivateKey
	PeersConfig      *PeersConfig
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
	txChannel    chan *core.Transaction
}

func MakeServer() *Server {
	SetConfigToDefaultPeers()
	conf, errConf := GetConfig()

	if errConf != nil {
		panic("failed to load config")
	}

	ip, errIp := GetPublicIP()

	if errIp != nil {
		panic("node is out network")
	}

	nodeId := "node_" + ip + defaultListenPort
	AddPeerToConfig(ip + defaultListenPort)

	privateKey := crypto.GeneratePrivateKey()

	opts := ServerOptions{
		APIListenAddress: ip + defaultAPIListenPort,
		PrivateKey:       &privateKey,
		ID:               nodeId,
		SeedNodes:        conf.Peers,
		ListenAddress:    ip + defaultListenPort,
		PeersConfig:      conf,
	}

	s, err := NewServer(&opts)

	if err != nil {
		panic(err)
	}

	return s
}

func NewServer(opts *ServerOptions) (*Server, error) {
	if opts.RPCDecodeFunc == nil {
		opts.RPCDecodeFunc = DefaultRPCDecodeFunc
	}

	if opts.Logger == nil {
		opts.Logger = log.NewLogfmtLogger(os.Stderr)
		opts.Logger = log.With(opts.Logger, "ID", opts.ID)
	}

	chain := core.NewBlockchain(opts.Logger, genesisBlock(*opts.PrivateKey))

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
		txChannel:    make(chan *core.Transaction),
	}

	s.TCPTransport.peerCh = peerCh

	if opts.RPCProcessor == nil {
		opts.RPCProcessor = s
	}

	if s.isValidator {
		//go s.validatorLoop()
	}

	if len(s.so.APIListenAddress) > 0 {
		apiServerConfig := api.ServerConfig{
			ListenAddr: s.so.APIListenAddress,
			Logger:     opts.Logger,
		}

		apiServer := api.NewServer(apiServerConfig, s.chain, s.txChannel)

		go apiServer.Start()

		_ = opts.Logger.Log("msg", "api server run on", opts.APIListenAddress)
	}

	return s, nil
}

func (s *Server) Start() {
	_ = s.TCPTransport.Start()
	s.boostrapPeers()
	go s.testMessages()

free:
	for {
		select {
		case peer := <-s.peerCh:
			s.peerMap[peer.conn.RemoteAddr()] = peer
			AddPeerToConfig(peer.conn.RemoteAddr().String())
			fmt.Println("ADDED NEW PEER: ", peer.conn.RemoteAddr().String())

			go peer.readLoop(s.rpcChannel)

			s.sendGetStatusMessage(peer)

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

		case tx := <-s.txChannel:
			err := s.processTransaction(tx)

			if err != nil {
				fmt.Println("process transaction error")
			}

		case <-s.quitChannel:
			break free
		}
	}

	_ = s.so.Logger.Log("msg", "Server shutdown")
}

func (s *Server) boostrapPeers() {
	for _, addr := range s.so.SeedNodes {
		if addr == s.so.ListenAddress {
			continue
		}

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

func (s *Server) testMessages() {
	for {
		time.Sleep(time.Second * 2)
		for _, peer := range s.peerMap {
			statusMessage := &StatusMessage{
				ActualHeight: s.chain.Height(),
				ID:           s.so.ID,
			}

			buf := &bytes.Buffer{}

			_ = gob.NewEncoder(buf).Encode(statusMessage)

			msg := NewMessage(MessageTypeStatus, buf.Bytes())
			data, _ := msg.Bytes()
			fmt.Println("sending to ", peer.conn.RemoteAddr().String(), " status message")
			_ = peer.Send(data)
		}
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

	_ = peer.Send(msgData)
}

func (s *Server) Broadcast(payload []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

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
	case *BlocksMessage:
		return s.processBlocksMessage(m.From, data)
	}
	return nil
}

func (s *Server) processTransaction(transaction *core.Transaction) error {
	hash := transaction.Hash(core.TxHasher{})
	_ = s.so.Logger.Log("msg", "adding new tx to mempool", "hash", hash, "mempool length", s.memPool.all.Count())

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

	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peerMap[from]

	if !ok {
		return errors.New("peer not found")
	}

	return peer.Send(data)
}

func (s *Server) processStatusMessage(addr net.Addr, m *StatusMessage) error {
	fmt.Println("I AM ", s.so.ID, ", ", "get from ", addr, " data: ", " ID: ", m.ID, " Height of chain: ", m.ActualHeight)

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

	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peerMap[addr]

	if !ok {
		return errors.New("peer not found")
	}

	fmt.Println("I AM ", s.so.ID, ", ", "wanna get blocks from ", addr, " : my chain height: ", getBlocksMessage.From, " To: ", getBlocksMessage.To)

	return peer.Send(data)
}

func (s *Server) processGetBlocksMessage(from net.Addr, m *GetBlocksMessage) error {
	var blocks []*core.Block
	bcHeight := s.chain.Height()

	if m.To == 0 {
		for i := m.From; i <= bcHeight; i++ {
			block, err := s.chain.Store.Get(i)

			if err != nil {
				return err
			}

			blocks = append(blocks, block)
		}
	}

	blocksMsg := BlocksMessage{
		blocks: blocks,
	}

	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(blocksMsg)

	if err != nil {
		return err
	}

	s.mu.RLock()
	defer s.mu.RUnlock()

	peer, ok := s.peerMap[from]
	if !ok {
		return errors.New("trying to send messageTypeBlocks - peer not found")
	}

	msg := NewMessage(MessageTypeBlocks, buf.Bytes())
	data, err := msg.Bytes()

	if err != nil {
		return err
	}

	errSend := peer.Send(data)

	if errSend != nil {
		return errSend
	}

	return nil
}

func (s *Server) processBlocksMessage(from net.Addr, m *BlocksMessage) error {
	_ = s.so.Logger.Log("msg", "received ", len(m.blocks), " blocks from ", from)

	for i := 0; i < len(m.blocks); i++ {
		s.chain.AddBlock(m.blocks[i])
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

func genesisBlock(key crypto.PrivateKey) *core.Block {
	h := &core.Header{
		Version:   1,
		DataHash:  types.Hash{},
		Timestamp: 0,
		Height:    0,
	}

	b := core.NewBlock(h, nil)
	_ = b.Sign(key)

	coinbase := crypto.ZeroPublicKey()

	tx := core.NewTransaction(nil, coinbase, coinbase.Address(), 1000000, 1)
	b.Transactions = append(b.Transactions, tx)

	return b
}
