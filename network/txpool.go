package network

import (
	"github.com/Phanile/uretra_network/core"
	"github.com/Phanile/uretra_network/types"
	"sync"
)

type TxSortedMap struct {
	lock   sync.RWMutex
	lookup map[types.Hash]*core.Transaction
	txs    *types.List[*core.Transaction]
}

func NewTxSortedMap() *TxSortedMap {
	return &TxSortedMap{
		lookup: make(map[types.Hash]*core.Transaction),
		txs:    types.NewList[*core.Transaction](),
	}
}

func (m *TxSortedMap) First() *core.Transaction {
	m.lock.RLock()
	defer m.lock.RUnlock()

	f := m.txs.Get(0)
	return m.lookup[f.Hash(core.TxHasher{})]
}

func (m *TxSortedMap) Get(hash types.Hash) *core.Transaction {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return m.lookup[hash]
}

func (m *TxSortedMap) Add(tx *core.Transaction) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	hash := tx.Hash(core.TxHasher{})

	_, ok := m.lookup[hash]

	if !ok {
		m.lookup[hash] = tx
		m.txs.Insert(tx)

		return true
	}

	return false
}

func (m *TxSortedMap) Remove(hash types.Hash) bool {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.txs.Remove(m.lookup[hash])
	delete(m.lookup, hash)

	return true
}

func (m *TxSortedMap) Count() uint16 {
	m.lock.RLock()
	defer m.lock.RUnlock()

	return uint16(len(m.lookup))
}

func (m *TxSortedMap) Contains(hash types.Hash) bool {
	m.lock.RLock()
	defer m.lock.RUnlock()

	_, ok := m.lookup[hash]

	return ok
}

func (m *TxSortedMap) Clear() {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.lookup = make(map[types.Hash]*core.Transaction)
	m.txs.Clear()
}
