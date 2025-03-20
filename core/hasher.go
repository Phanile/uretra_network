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

func (bh *BlockHasher) Hash(b *Block) types.Hash {
	return sha256.Sum256(b.HeaderData())
}
