package core

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"uretra-network/crypto"
)

func TestTransaction_Sign(t *testing.T) {
	data := []byte("AMOUNT 5000 BTC")
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.NotNil(t, tx.signature)
	assert.NotNil(t, tx.PublicKey)
}

func TestTransaction_Verify(t *testing.T) {
	data := []byte("AMOUNT 5000 BTC")
	privKey := crypto.GeneratePrivateKey()

	tx := &Transaction{
		Data: data,
	}

	assert.Nil(t, tx.Sign(privKey))
	assert.True(t, tx.Verify())

	otherKey := crypto.GeneratePrivateKey()
	tx.PublicKey = otherKey.PublicKey()

	assert.False(t, tx.Verify())
}
