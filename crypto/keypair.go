package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"math/big"
	"uretra-network/types"
)

type PrivateKey struct {
	key *ecdsa.PrivateKey
}

func GeneratePrivateKey() PrivateKey {
	key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)

	if err != nil {
		panic(err)
	}

	return PrivateKey{
		key: key,
	}
}

type PublicKey struct {
	key *ecdsa.PublicKey
}

func (pk PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		key: &pk.key.PublicKey,
	}
}

func (pk PublicKey) Address() types.Address {
	sum := sha256.Sum256(elliptic.MarshalCompressed(pk.key, pk.key.X, pk.key.Y))

	return types.AddressFromBytes(sum[len(sum)-20:])
}

type Signature struct {
	r *big.Int
	s *big.Int
}

func (pk PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, pk.key, data)

	if err != nil {
		return nil, err
	}

	return &Signature{
		r: r,
		s: s,
	}, nil
}

func (signature *Signature) VerifySignature(pk *PublicKey, data []byte) bool {
	return ecdsa.Verify(pk.key, data, signature.r, signature.s)
}
