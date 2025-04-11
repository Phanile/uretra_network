package core

import (
	"github.com/go-kit/log"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"uretra-network/types"
)

func TestStorage_Get(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.Hash{})
	bc := NewBlockchain(log.NewLogfmtLogger(os.Stderr), b)

	lenBlocks := 64

	for i := 0; i < lenBlocks; i++ {
		newBlock := randomBlockWithSignature(t, uint32(i+1), getPrevBlockHash(t, bc, uint32(i+1)))
		assert.True(t, bc.AddBlock(newBlock))
		header, err := bc.GetHeader(newBlock.Header.Height)
		assert.Nil(t, err)
		assert.Equal(t, header, newBlock.Header)
		block, errGet := bc.Store.Get(uint32(i + 1))
		assert.Nil(t, errGet)
		assert.Equal(t, block, newBlock)
	}
}
