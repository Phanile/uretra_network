package core

import (
	"time"
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
