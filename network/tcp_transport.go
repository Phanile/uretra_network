package network

import (
	"bytes"
	"io"
	"net"
	"time"
)

type PeerInfo struct {
	Peer             *TCPPeer
	BlockchainHeight uint32
	PingTime         time.Duration
}

type TCPPeer struct {
	conn     net.Conn
	Outgoing bool
}

func (peer *TCPPeer) Send(data []byte) error {
	_, err := peer.conn.Write(data)
	return err
}

func (peer *TCPPeer) readLoop(rpcCh chan RPC) {
	buf := make([]byte, 4096)

	for {
		n, err := peer.conn.Read(buf)

		if err == io.EOF {
			continue
		}

		if err != nil {
			continue
		}

		msg := buf[:n]

		rpcCh <- RPC{
			From:    peer.conn.RemoteAddr(),
			Payload: bytes.NewReader(msg),
		}
	}
}

type TCPTransport struct {
	peerCh     chan *TCPPeer
	ListenAddr string
	listener   net.Listener
}

func NewTCPTransport(addr string, peerCh chan *TCPPeer) *TCPTransport {
	return &TCPTransport{
		ListenAddr: addr,
		peerCh:     peerCh,
	}
}

func (t *TCPTransport) Start() error {
	ln, err := net.Listen("tcp", t.ListenAddr)

	if err != nil {
		return err
	}

	t.listener = ln

	go t.acceptLoop()

	return nil
}

func (t *TCPTransport) acceptLoop() {
	for {
		connection, err := t.listener.Accept()

		if err != nil {
			continue
		}

		peer := &TCPPeer{
			conn: connection,
		}

		t.peerCh <- peer
	}
}
