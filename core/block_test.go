package core

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
	"uretra-network/types"
)

func TestHeader_Encode_Decode_Binary(t *testing.T) {
	h := &Header{
		version:   1,
		prevBlock: types.RandomHash(),
		timestamp: time.Now().UnixNano(),
		height:    10,
		nonce:     893274,
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, h.EncodeBinary(buf))

	hDecode := &Header{}
	assert.Nil(t, hDecode.DecodeBinary(buf))

	assert.Equal(t, h, hDecode)
}

func TestBlock_Encode_Decode_Binary(t *testing.T) {
	h := Header{
		version:   1,
		prevBlock: types.RandomHash(),
		timestamp: time.Now().UnixNano(),
		height:    10,
		nonce:     893274,
	}

	b := &Block{
		header:       h,
		transactions: nil,
		hash:         types.RandomHash(),
	}

	buf := &bytes.Buffer{}
	assert.Nil(t, b.EncodeBinary(buf))

	bDecode := &Block{}
	assert.Nil(t, bDecode.DecodeBinary(buf))

	assert.Equal(t, b, bDecode)
}

func TestBlock_Hash(t *testing.T) {
	h := Header{
		version:   1,
		prevBlock: types.RandomHash(),
		timestamp: time.Now().UnixNano(),
		height:    10,
		nonce:     893274,
	}

	b := &Block{
		header:       h,
		transactions: nil,
	}

	hash := b.Hash()
	fmt.Println(hash)

	assert.False(t, hash.IsEmptyOrZero())
	assert.Equal(t, b.hash, hash)
}
