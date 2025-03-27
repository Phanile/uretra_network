package network

import (
	"sync"
	"uretra-network/core"
	"uretra-network/types"
)

type TxPool struct {
	all       *TxSortedMap
	pending   *TxSortedMap
	maxLength uint16
}

func NewTxPool(maxLength uint16) *TxPool {
	return &TxPool{
		maxLength: maxLength,
		all:       NewTxSortedMap(),
		pending:   NewTxSortedMap(),
	}
}

func (p *TxPool) Add(tr *core.Transaction) bool {
	if p.all.Count() == p.maxLength {
		return false
	}

	if p.all.Contains(tr.Hash(core.TxHasher{})) {
		return false
	}

	return p.all.Add(tr) && p.pending.Add(tr)
}

func (p *TxPool) Contains(hash types.Hash) bool {
	return p.all.Contains(hash)
}

func (p *TxPool) ClearPending() {
	p.pending.Clear()
}

func (p *TxPool) PendingCount() uint16 {
	return p.pending.Count()
}

func (p *TxPool) Pending() []*core.Transaction {
	return p.pending.txs.Data
}

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
