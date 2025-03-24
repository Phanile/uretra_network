package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"uretra-network/types"
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

	errData := binary.Write(buf, binary.LittleEndian, tx.Data)
	//errFrom := binary.Write(buf, binary.LittleEndian, tx.From)

	if errData != nil {
		panic("Tx is not hashable")
	}

	return sha256.Sum256(buf.Bytes())
}
