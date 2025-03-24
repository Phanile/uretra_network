package network

import (
	"uretra-network/core"
	"uretra-network/types"
)

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (txp *TxPool) Len() uint32 {
	return uint32(len(txp.transactions))
}

func (txp *TxPool) Flush() {
	txp.transactions = make(map[types.Hash]*core.Transaction)
}

func (txp *TxPool) Add(t *core.Transaction) bool {
	hash := t.Hash(core.TxHasher{})

	txp.transactions[hash] = t

	return true
}

func (txp *TxPool) Has(hash types.Hash) bool {
	_, ok := txp.transactions[hash]
	return ok
}
