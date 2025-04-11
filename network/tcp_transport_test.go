package network

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

func TestTCPTransport_Start(t *testing.T) {
	tr := NewTCPTransport(":3500", make(chan *TCPPeer))
	assert.Nil(t, tr.Start())
}

func TestTCPTransport_acceptLoop(t *testing.T) {
	tr := NewTCPTransport(":3500", make(chan *TCPPeer))

	go tr.Start()

	time.Sleep(time.Second * 1)
	for i := 0; i < 10; i++ {
		conn, err := net.Dial("tcp", ":3500")
		assert.Nil(t, err)
		fmt.Println("connection: ", conn)
		conn.Write([]byte("Hello world"))
	}
}
