package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
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

	binary.Write(buf, binary.LittleEndian, tx.Data)
	binary.Write(buf, binary.LittleEndian, tx.From)
	binary.Write(buf, binary.LittleEndian, tx.To)
	binary.Write(buf, binary.LittleEndian, tx.Value)
	binary.Write(buf, binary.LittleEndian, tx.Nonce)

	return sha256.Sum256(buf.Bytes())
}
