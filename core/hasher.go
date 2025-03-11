package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"uretra-network/types"
)

type BlockHasher struct {
}

type Hasher[T any] interface {
	Hash(T) types.Hash
}

func (bh *BlockHasher) Hash(b *Block) types.Hash {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)

	if err := enc.Encode(b.Header); err != nil {
		panic(err)
	}

	return sha256.Sum256(buf.Bytes())
}
