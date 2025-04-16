package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/gob"
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

func (pk *PublicKey) GobEncode() ([]byte, error) {
	if pk.Key == nil {
		return nil, nil
	}

	return elliptic.MarshalCompressed(pk.Key.Curve, pk.Key.X, pk.Key.Y), nil
}

func (pk *PublicKey) GobDecode(data []byte) error {
	if len(data) == 0 {
		pk.Key = nil
		return nil
	}
	x, y := elliptic.UnmarshalCompressed(elliptic.P256(), data)
	pk.Key = &ecdsa.PublicKey{
		Curve: elliptic.P256(),
		X:     x,
		Y:     y,
	}
	return nil
}

func (pk PrivateKey) PublicKey() PublicKey {
	return PublicKey{
		Key: &pk.key.PublicKey,
	}
}

func (pk PublicKey) Address() types.Address {
	if pk.Key == nil {
		return types.Address{}
	}

	sum := sha256.Sum256(elliptic.MarshalCompressed(pk.Key, pk.Key.X, pk.Key.Y))

	return types.AddressFromBytes(sum[len(sum)-20:])
}

type Signature struct {
	r *big.Int
	s *big.Int
}

func (s *Signature) R() *big.Int {
	if s == nil {
		return nil
	}
	return s.r
}
func (s *Signature) S() *big.Int {
	if s == nil {
		return nil
	}
	return s.s
}

func (s *Signature) GobEncode() ([]byte, error) {
	if s == nil {
		return nil, nil
	}

	return append(s.r.Bytes(), s.s.Bytes()...), nil
}

func (s *Signature) GobDecode(data []byte) error {
	if len(data) == 0 {
		return nil
	}

	half := len(data) / 2
	s.r = new(big.Int).SetBytes(data[:half])
	s.s = new(big.Int).SetBytes(data[half:])
	return nil
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
	return ecdsa.Verify(pk.Key, data, signature.r, signature.s)
}

func ZeroPublicKey() PublicKey {
	return PublicKey{
		Key: &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     big.NewInt(0),
			Y:     big.NewInt(0),
		},
	}
}

func init() {
	gob.Register(&PublicKey{})
	gob.Register(&Signature{})
}
