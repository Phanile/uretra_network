package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/Phanile/uretra_network/types"
)

type HeaderHasher struct {
}

type Hasher[T any] interface {
	Hash(T) types.Hash
}

func (HeaderHasher) Hash(h *Header) types.Hash {
	return sha256.Sum256(h.Bytes())
}

type TxHasher struct {
}

func (TxHasher) Hash(tx *Transaction) types.Hash {
	buf := &bytes.Buffer{}

	err := gob.NewEncoder(buf).Encode(tx)

	if err != nil {
		panic("Tx is not hashable")
	}

	return sha256.Sum256(buf.Bytes())
}
