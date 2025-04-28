package core

import (
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type Storage interface {
	Put(*Block) error
	Get(height uint32) (*Block, error)
}

type MemoryStorage struct {
	mu         sync.RWMutex
	blockchain *Blockchain
	baseDir    string
}

func NewMemoryStorage(blockchain *Blockchain) *MemoryStorage {
	return &MemoryStorage{
		blockchain: blockchain,
		baseDir:    "./storageBlocks/",
	}
}

func (ms *MemoryStorage) Put(b *Block) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	data, err := json.Marshal(b)

	if err != nil {
		return err
	}

	filename := fmt.Sprintf("%s/%d.json", ms.baseDir, b.Header.Height)
	return os.WriteFile(filename, data, 0600)
}

func (ms *MemoryStorage) Get(height uint32) (*Block, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	filename := fmt.Sprintf("%s/%d.json", ms.baseDir, height)
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var block Block
	if errUnmarshall := json.Unmarshal(data, &block); errUnmarshall != nil {
		return nil, errUnmarshall
	}

	return &block, nil
}
