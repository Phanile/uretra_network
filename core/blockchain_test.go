package core

import (
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"uretra-network/types"
)

func TestBlockchain_Create(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})
	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestBlockchain_HasBlock(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})
	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)

	assert.True(t, bc.HasBlock(0))
}

func TestBlockchain_AddBlock(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})

	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)
	lenBlocks := 512

	for i := 0; i < lenBlocks; i++ {
		newBlock := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.True(t, bc.AddBlock(newBlock))
	}

	assert.False(t, bc.AddBlock(randomBlockWithSignature(t, 100, types.RandomHash())))
	assert.Equal(t, bc.Height(), uint32(lenBlocks))
}

func TestBlockchain_AddBlockToHigh(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.RandomHash())

	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)

	assert.False(t, bc.AddBlock(randomBlockWithSignature(t, 3, getPrevBlockHash(t, bc, uint32(1)))))
}

func TestBlockchain_GetHeader(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})

	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)
	lenBlocks := 512

	for i := 0; i < lenBlocks; i++ {
		newBlock := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.True(t, bc.AddBlock(newBlock))
		header, err := bc.GetHeader(newBlock.Header.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, newBlock.Header)
	}
}

func getPrevBlockHash(t *testing.T, bc *Blockchain, height uint32) types.Hash {
	header, err := bc.GetHeader(height - 1)
	assert.Nil(t, err)

	return HeaderHasher{}.Hash(header)
}
