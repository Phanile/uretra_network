package core

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"uretra-network/crypto"
	"uretra-network/types"
)

func randomBlock(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	tr1 := RandomTxWithSignature(t)

	return NewBlock(h, []*Transaction{tr1})
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	h := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	tr1 := RandomTxWithSignature(t)

	privateKey := crypto.GeneratePrivateKey()

	return NewBlockWithPrivateKey(h, []*Transaction{tr1}, privateKey)
}

func TestBlock_Hash(t *testing.T) {
	b := randomBlock(t, 0, types.RandomHash())
	fmt.Println(b.Hash(&HeaderHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.RandomHash())
	err := b.Sign(privKey)

	if err != nil {
		return
	}

	assert.Nil(t, b.Sign(privKey))
}

func TestBlock_Verify(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlock(t, 0, types.RandomHash())
	err := b.Sign(privKey)
	b.Validator = privKey.PublicKey()

	if err != nil {
		return
	}

	assert.True(t, b.Verify())
}
