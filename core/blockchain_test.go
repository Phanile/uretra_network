package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBlockchain_Create(t *testing.T) {
	b := randomBlock(0)
	bc := NewBlockchain(b)

	assert.NotNil(t, bc.validator)
	assert.Equal(t, bc.Height(), uint32(0))
}

func TestBlockchain_HasBlock(t *testing.T) {
	b := randomBlockWithSignature(0)
	bc := NewBlockchain(b)

	assert.True(t, bc.HasBlock(0))
}

func TestBlockchain_AddBlock(t *testing.T) {
	b := randomBlockWithSignature(0)

	bc := NewBlockchain(b)
	lenBlocks := 512

	for i := 0; i < lenBlocks; i++ {
		newBlock := randomBlockWithSignature(uint32(i + 1))
		assert.True(t, bc.addBlock(newBlock))
	}

	assert.False(t, bc.addBlock(randomBlock(100)))
	assert.Equal(t, bc.Height(), uint32(lenBlocks))
}
