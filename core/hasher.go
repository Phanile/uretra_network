package core

import (
	"crypto/sha256"
	"uretra-network/types"
)

type BlockHasher struct {
}

type Hasher[T any] interface {
	Hash(T) types.Hash
}

func (BlockHasher) Hash(h *Header) types.Hash {
	return sha256.Sum256(h.Bytes())
}
