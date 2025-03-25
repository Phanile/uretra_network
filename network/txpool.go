package network

import (
	"sort"
	"uretra-network/core"
	"uretra-network/types"
)

type TxMapSorter struct {
	transactions []*core.Transaction
}

func NewTxMapSorter(txMap map[types.Hash]*core.Transaction) *TxMapSorter {
	txx := make([]*core.Transaction, len(txMap))

	i := 0
	for _, val := range txMap {
		txx[i] = val
		i++
	}

	s := &TxMapSorter{
		transactions: txx,
	}

	sort.Sort(s)

	return s
}

func (s *TxMapSorter) Len() int {
	return len(s.transactions)
}

func (s *TxMapSorter) Less(i, j int) bool {
	return s.transactions[i].FirstSeen() < s.transactions[j].FirstSeen()
}

func (s *TxMapSorter) Swap(i, j int) {
	s.transactions[i], s.transactions[j] = s.transactions[j], s.transactions[i]
}

type TxPool struct {
	transactions map[types.Hash]*core.Transaction
}

func NewTxPool() *TxPool {
	return &TxPool{
		transactions: make(map[types.Hash]*core.Transaction),
	}
}

func (txp *TxPool) Transactions() []*core.Transaction {
	return NewTxMapSorter(txp.transactions).transactions
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
