package core

import (
	"bytes"
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
	assert.NotNil(t, tx.Signature)
	assert.NotNil(t, tx.From)
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
	tx.From = otherKey.PublicKey()

	assert.False(t, tx.Verify())
}

func TestTransaction_Decode(t *testing.T) {
	tx := RandomTxWithSignature(t)
	buf := &bytes.Buffer{}
	assert.Nil(t, tx.Encode(NewGobTxEncoder(buf)))

	txDecoded := &Transaction{}

	assert.Nil(t, txDecoded.Decode(NewGobTxDecoder(buf)))

	assert.Equal(t, tx, txDecoded)
}

func RandomTxWithSignature(t *testing.T) *Transaction {
	data := []byte("AMOUNT 5000 BTC")
	privKey := crypto.GeneratePrivateKey()
	sign, err := privKey.Sign(data)
	assert.Nil(t, err)

	return &Transaction{
		Data:      data,
		From:      privKey.PublicKey(),
		Signature: sign,
	}
}
