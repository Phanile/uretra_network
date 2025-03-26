package network

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestLocalTransport_Broadcast(t *testing.T) {
	tra := NewLocalTransport("228.228.288.228")
	trb := NewLocalTransport("228.224.248.224")
	trc := NewLocalTransport("228.214.241.211")

	assert.Nil(t, tra.Connect(trb))
	assert.Nil(t, tra.Connect(trc))

	msg := []byte("airdrop 5000 btc")
	assert.Nil(t, tra.Broadcast(msg))

	rpcb := <-trb.Consume()
	bb, errb := ioutil.ReadAll(rpcb.Payload)
	assert.Nil(t, errb)
	assert.Equal(t, bb, msg)

	rpcc := <-trc.Consume()
	bc, errc := ioutil.ReadAll(rpcc.Payload)
	assert.Nil(t, errc)
	assert.Equal(t, bc, msg)
}
