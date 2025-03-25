package network

import (
	"fmt"
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

func TestTxPool_Sort(t *testing.T) {
	pool := NewTxPool()
	tr := core.NewTransaction([]byte("a lot of data"))
	tr1 := core.NewTransaction([]byte("a lots of data"))
	tr2 := core.NewTransaction([]byte("get: 1000 usdt"))

	assert.True(t, pool.Add(tr))
	assert.True(t, pool.Add(tr1))
	assert.True(t, pool.Add(tr2))

	sorter := NewTxMapSorter(pool.transactions)

	fmt.Println(sorter.transactions)

	assert.Equal(t, sorter.Len(), 3)

	sorter.Swap(0, 2)
	fmt.Println(sorter.transactions)
}
