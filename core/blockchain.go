package core

import (
	"fmt"
	"github.com/go-kit/log"
	"sync"
)

type Blockchain struct {
	logger    log.Logger
	store     Storage
	lock      sync.RWMutex
	headers   []*Header
	validator Validator
}

func NewBlockchain(l log.Logger, genesis *Block) *Blockchain {
	bc := &Blockchain{
		logger:  l,
		headers: []*Header{},
		store:   &MemoryStorage{},
	}

	bc.validator = NewBlockValidator(bc)

	err := bc.addBlockWithoutValidation(genesis)

	if err != nil {
		panic(err)
	}

	return bc
}

func (bc *Blockchain) AddBlock(b *Block) bool {
	if bc.validator.ValidateBlock(b) {
		err := bc.addBlockWithoutValidation(b)

		if err != nil {
			return false
		}

		return true
	}

	return false
}

func (bc *Blockchain) addBlockWithoutValidation(b *Block) error {
	bc.lock.Lock()
	bc.headers = append(bc.headers, b.Header)
	defer bc.lock.Unlock()

	_ = bc.logger.Log("msg", "new block", "hash", b.Hash(HeaderHasher{}), "height", b.Header.Height, "txs", len(b.Transactions))

	return bc.store.Put(b)
}

func (bc *Blockchain) HasBlock(height uint32) bool {
	bc.lock.RLock()
	defer bc.lock.RUnlock()

	return height <= bc.Height()
}

func (bc *Blockchain) GetHeader(height uint32) (*Header, error) {
	if height > bc.Height() {
		return nil, fmt.Errorf("trying get too high header (%d)", height)
	}

	bc.lock.Lock()
	defer bc.lock.Unlock()

	return bc.headers[height], nil
}

func (bc *Blockchain) SetValidator(val Validator) {
	bc.validator = val
}

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}
