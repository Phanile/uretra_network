package core

import (
	"errors"
	"sync"
)

type Storage interface {
	Put(*Block) error
	Get(height uint32) (*Block, error)
}

type MemoryStorage struct {
	mu         sync.RWMutex
	blocks     []*Block
	blockchain *Blockchain
}

func NewMemoryStorage(blockchain *Blockchain) *MemoryStorage {
	return &MemoryStorage{
		blockchain: blockchain,
	}
}

func (ms *MemoryStorage) Put(b *Block) error {
	ms.blocks = append(ms.blocks, b)
	return nil
}

func (ms *MemoryStorage) Get(height uint32) (*Block, error) {
	if height > ms.blockchain.Height() {
		return nil, errors.New("too high height of block")
	}

	ms.mu.Lock()
	defer ms.mu.Unlock()

	ok := ms.blocks[height]

	if ok != nil {
		return ok, nil
	}

	return nil, errors.New("block not found")
}
