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
	Key *ecdsa.PublicKey
}

func (pk PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &pk.key.PublicKey,
	}
}

func (pk PublicKey) Address() types.Address {
	sum := sha256.Sum256(elliptic.MarshalCompressed(pk.Key, pk.Key.X, pk.Key.Y))

	return types.AddressFromBytes(sum[len(sum)-20:])
}

type Signature struct {
	R *big.Int
	S *big.Int
}

func (pk PrivateKey) Sign(data []byte) (*Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, pk.key, data)

	if err != nil {
		return nil, err
	}

	return &Signature{
		R: r,
		S: s,
	}, nil
}

func (signature *Signature) VerifySignature(pk *PublicKey, data []byte) bool {
	return ecdsa.Verify(pk.Key, data, signature.R, signature.S)
}
