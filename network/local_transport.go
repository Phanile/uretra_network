package network

import (
	"bytes"
	"fmt"
	"sync"
)

type LocalTransport struct {
	address        NetAddress
	consumeChannel chan RPC
	lock           sync.RWMutex
	peers          map[NetAddress]*LocalTransport
}

func NewLocalTransport(addr NetAddress) Transport {
	return &LocalTransport{
		address:        addr,
		consumeChannel: make(chan RPC, 1024),
		peers:          make(map[NetAddress]*LocalTransport),
	}
}

func (t *LocalTransport) Consume() <-chan RPC {
	return t.consumeChannel
}

func (t *LocalTransport) Connect(lc Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[lc.Address()] = lc.(*LocalTransport)

	return nil
}

func (t *LocalTransport) SendMessage(to NetAddress, data []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	if t.Address() == to {
		return nil
	}

	peer, ok := t.peers[to]

	if !ok {
		return fmt.Errorf("send message error: peer does not exist")
	}

	peer.consumeChannel <- RPC{
		From:    t.address,
		Payload: bytes.NewReader(data),
	}

	return nil
}

func (t *LocalTransport) Broadcast(payload []byte) error {
	t.lock.RLock()
	defer t.lock.RUnlock()

	for _, peer := range t.peers {
		err := t.SendMessage(peer.Address(), payload)

		if err != nil {
			return err
		}
	}

	return nil
}

func (t *LocalTransport) Address() NetAddress {
	return t.address
}
