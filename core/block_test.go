package core

import (
	"bytes"
	"fmt"
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBlock_Hash(t *testing.T) {
	b := randomBlockWithSignature(t, 0, types.RandomHash())
	fmt.Println(b.Hash(&HeaderHasher{}))
}

func TestBlock_Sign(t *testing.T) {
	privKey := crypto.GeneratePrivateKey()
	b := randomBlockWithSignature(t, 0, types.RandomHash())
	err := b.Sign(privKey)

	if err != nil {
		return
	}

	assert.Nil(t, b.Sign(privKey))
}

func TestBlock_Verify(t *testing.T) {
	assert.True(t, randomBlockWithSignature(t, 0, types.RandomHash()).Verify())
}

func TestBlock_DecodeEncode(t *testing.T) {
	b1 := randomBlockWithSignature(t, 1, types.Hash{})
	buf := &bytes.Buffer{}

	assert.Nil(t, b1.Encode(NewGobBlockEncoder(buf)))

	b2 := &Block{}

	assert.Nil(t, b2.Decode(NewGobBlockDecoder(buf)))

	b2.Hash(HeaderHasher{})
	assert.Equal(t, b1, b2)
}

func randomBlockWithSignature(t *testing.T, height uint32, prevBlockHash types.Hash) *Block {
	privateKey := crypto.GeneratePrivateKey()

	h := &Header{
		Version:       1,
		PrevBlockHash: prevBlockHash,
		Timestamp:     time.Now().UnixNano(),
		Height:        height,
	}

	tr1 := RandomTxWithSignature(t)

	b := NewBlock(h, []*Transaction{tr1})
	b.Validator = privateKey.PublicKey()
	assert.Nil(t, b.Sign(privateKey))

	dataHash, err := b.CalculateDataHash(b.Transactions)
	assert.Nil(t, err)

	b.Header.DataHash = dataHash
	b.Hash(&HeaderHasher{})

	return b
}
