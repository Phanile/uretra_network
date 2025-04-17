package core

import (
	"errors"
	"github.com/Phanile/uretra_network/types"
	"sync"
)

var (
	AccountNotFoundError         = errors.New("account not found")
	AccountNotEnoughBalanceError = errors.New("account not enough balance")
)

type Accounts struct {
	mu    sync.RWMutex
	state map[types.Address]*Account
}

type Account struct {
	Address types.Address
	Balance uint64
}

func NewAccounts() *Accounts {
	return &Accounts{
		state: make(map[types.Address]*Account),
	}
}

func (a *Accounts) NewAccount(addr types.Address) *Account {
	a.mu.Lock()
	defer a.mu.Unlock()

	acc := &Account{
		Address: addr,
	}

	a.state[addr] = acc

	return acc
}

func (a *Accounts) GetAccount(addr types.Address) (*Account, error) {
	a.mu.RLock()
	defer a.mu.RUnlock()

	return a.getNoLockAccount(addr)
}

func (a *Accounts) getNoLockAccount(addr types.Address) (*Account, error) {
	acc, ok := a.state[addr]

	if !ok {
		return nil, AccountNotFoundError
	}

	return acc, nil
}

func (a *Accounts) GetBalance(addr types.Address) (uint64, error) {
	acc, err := a.GetAccount(addr)

	if err != nil {
		return 0, err
	}

	return acc.Balance, nil
}

func (a *Accounts) Transfer(from, to types.Address, value uint64) error {
	a.mu.Lock()
	defer a.mu.Unlock()

	fromAcc, errGetAcc := a.getNoLockAccount(from)

	if errGetAcc != nil {
		return errGetAcc
	}

	if fromAcc.Balance < value {
		return AccountNotEnoughBalanceError
	}

	fromAcc.Balance -= value

	if a.state[to] == nil {
		a.state[to] = &Account{
			Address: to,
		}
	}

	a.state[to].Balance += value

	return nil
}

func (a *Accounts) AddBalance(to types.Address, value uint64) error {
	if a.state[to] == nil {
		a.state[to] = a.NewAccount(to)
	}

	a.state[to].Balance += value

	return nil
}
