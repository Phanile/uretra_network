package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"uretra-network/crypto"
	"uretra-network/types"
)

func randomBlock(height uint32) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	tr1 := &Transaction{
		Data: []byte("test data for test 21 04"),
	}

	return NewBlock(h, []*Transaction{tr1})
}

func randomBlockWithSignature(height uint32) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: types.RandomHash(),
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	tr1 := &Transaction{
		Data: []byte("test data for test 21 04"),
	}

	privateKey := crypto.GeneratePrivateKey()

	return NewBlockWithPrivateKey(h, []*Transaction{tr1}, privateKey)
}

func TestBlock_Hash(t *testing.T) {
	b := randomBlock(0)
	fmt.Println(b.Hash(&BlockHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	err := b.Sign(privKey)

	if err != nil {
		return
	}

	assert.Nil(t, b.Sign(privKey))
}

func TestBlock_Verify(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(0)
	err := b.Sign(privKey)
	b.Validator = privKey.PublicKey()

	if err != nil {
		return
	}

	assert.True(t, b.Verify(*b.Signature))
}
