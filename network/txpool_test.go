package network

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uretra-network/core"
)

func TestTxPool(t *testing.T) {
	pool := NewTxPool()
	assert.Equal(t, pool.Len(), uint32(0))
}

func TestTxPool_Add(t *testing.T) {
	pool := NewTxPool()
	tr := core.NewTransaction([]byte("a lot of data"))

	assert.True(t, pool.Add(tr))
	assert.Equal(t, pool.Len(), uint32(1))

	pool.Flush()
	assert.Equal(t, pool.Len(), uint32(0))
}
