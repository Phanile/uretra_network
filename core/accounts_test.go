package core

import (
	"github.com/Phanile/uretra_network/crypto"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAccounts_GetBalance(t *testing.T) {
	pk := crypto.GeneratePrivateKey()
	a := NewAccounts()
	assert.Nil(t, a.AddBalance(pk.PublicKey().Address(), uint64(100)))

	balance, ok := a.GetBalance(pk.PublicKey().Address())
	assert.Nil(t, ok)
	assert.Equal(t, balance, uint64(100))
}

func TestAccounts_Transfer(t *testing.T) {
	pkAlice := crypto.GeneratePrivateKey()
	pkBob := crypto.GeneratePrivateKey()
	a := NewAccounts()

	assert.Nil(t, a.AddBalance(pkAlice.PublicKey().Address(), uint64(1000)))
	assert.Nil(t, a.Transfer(pkAlice.PublicKey().Address(), pkBob.PublicKey().Address(), uint64(700)))

	balanceAlice, okAlice := a.GetBalance(pkAlice.PublicKey().Address())
	balanceBob, okBob := a.GetBalance(pkBob.PublicKey().Address())
	assert.Nil(t, okAlice)
	assert.Nil(t, okBob)

	assert.Equal(t, balanceAlice, uint64(300))
	assert.Equal(t, balanceBob, uint64(700))
}
