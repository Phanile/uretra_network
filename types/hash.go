package types

import (
	"crypto/rand"
	"encoding/hex"
)

type Hash [32]uint8

func (h Hash) String() string {
	return hex.EncodeToString(h[:])
}

func HashFromBytes(b []byte) Hash {
	if len(b) != 32 {
		panic("hash length must be 32 bytes")
	}

	var res [32]uint8

	for i := 0; i < 32; i++ {
		res[i] = b[i]
	}

	return res
}

func RandomHash() Hash {
	rnd := make([]byte, 32)
	rand.Read(rnd)
	return HashFromBytes(rnd)
}

func (h Hash) IsEmptyOrZero() bool {
	for i := 0; i < 32; i++ {
		if h[i] == 0 {
			return true
		}
	}

	return false
}
