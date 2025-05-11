package core

import (
	"encoding/hex"
	"fmt"
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
	"github.com/go-kit/log"
	"sync"
)

type Blockchain struct {
	logger        log.Logger
	Store         Storage
	lock          sync.RWMutex
	headers       []*Header
	validator     Validator
	state         *State
	accountsState *Accounts
}

func NewBlockchain(l log.Logger, genesis *Block) *Blockchain {
	bc := &Blockchain{
		logger:  l,
		headers: []*Header{},
		state:   NewState(),
	}

	bc.Store = NewMemoryStorage(bc)
	bc.validator = NewBlockValidator(bc)
	bc.accountsState = NewAccounts()

	bc.accountsState.NewAccount(crypto.ZeroPublicKey().Address()) //coinbase account

	// TEST
	addrBytes, _ := hex.DecodeString("b2f1c7c07b3eb376ad89f3e8afba8b005616cb63")
	bc.accountsState.NewAccount(types.AddressFromBytes(addrBytes))
	_ = bc.accountsState.AddBalance(types.AddressFromBytes(addrBytes), 1000000)
	// TEST

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
	defer bc.lock.Unlock()

	validTxs := make([]*Transaction, 0, len(b.Transactions))

	for i := 0; i < len(b.Transactions); i++ {
		err := bc.handleTransaction(b.Transactions[i])

		if err != nil {
			_ = bc.logger.Log(
				"msg", "transaction failed",
				"hash", TxHasher{}.Hash(b.Transactions[i]),
				"error", err,
			)
			continue
		}

		validTxs = append(validTxs, b.Transactions[i])
	}

	b.Transactions = validTxs
	hash, hashErr := CalculateDataHash(b.Transactions)

	if hashErr != nil {
		return hashErr
	}

	b.Header.DataHash = hash

	bc.headers = append(bc.headers, b.Header)

	_ = bc.logger.Log("msg", "new block", "hash", b.Hash(HeaderHasher{}), "height", b.Header.Height, "txs", len(b.Transactions))

	return bc.Store.Put(b)
}

func (bc *Blockchain) handleTransaction(t *Transaction) error {
	if len(t.Data) > 0 {
		vm := NewVM(t.Data, bc.state)
		err := vm.Run()

		if err != nil {
			return err
		}

		fmt.Printf("state: %+v\n", bc.state)
	}

	if t.Value > 0 {
		err := bc.accountsState.Transfer(t.From.Address(), t.To, t.Value)

		if err != nil {
			return err
		}
	}

	return nil
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

func (bc *Blockchain) Height() uint32 {
	return uint32(len(bc.headers) - 1)
}

func (bc *Blockchain) GetAccounts() *Accounts {
	return bc.accountsState
}
