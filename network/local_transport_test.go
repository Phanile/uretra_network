package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnect(t *testing.T) {
	a := NewLocalTransport("a")
	b := NewLocalTransport("b")

	a.Connect(b)
	b.Connect(a)
	assert.Equal(t, a.peers[b.address], b)
	assert.Equal(t, b.peers[a.address], a)
}

func TestSendMessage(t *testing.T) {
	a := NewLocalTransport("a")
	b := NewLocalTransport("b")

	a.Connect(b)
	b.Connect(a)

	msg := []byte("Hi From A")
	assert.Nil(t, a.SendMessage(b.Address(), msg))

	rpc := <-b.Consume()
	assert.Equal(t, rpc.Data, msg)
	assert.Equal(t, rpc.From, a.Address())
}
