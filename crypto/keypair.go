package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/gob"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/Phanile/uretra_network/types"
	"math/big"
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

func (pk PrivateKey) Bytes() ([]byte, error) {
	if pk.key == nil {
		return nil, errors.New("private key is nil")
	}
	x509Encoded, err := x509.MarshalECPrivateKey(pk.key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: x509Encoded,
	}), nil
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

func (pk PublicKey) Bytes() ([]byte, error) {
	if pk.Key == nil {
		return nil, errors.New("public key is nil")
	}
	x509Encoded, err := x509.MarshalPKIXPublicKey(pk.Key)
	if err != nil {
		return nil, err
	}
	return pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509Encoded,
	}), nil
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

func PrivateKeyFromBytes(b []byte) (PrivateKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return PrivateKey{}, errors.New("failed to parse PEM block")
	}
	key, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return PrivateKey{}, err
	}
	return PrivateKey{key: key}, nil
}

func PublicKeyFromBytes(b []byte) (PublicKey, error) {
	block, _ := pem.Decode(b)
	if block == nil {
		return PublicKey{}, errors.New("failed to parse PEM block")
	}
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return PublicKey{}, err
	}
	return PublicKey{Key: pub.(*ecdsa.PublicKey)}, nil
}

func (pk PublicKey) MarshalJSON() ([]byte, error) {
	if pk.Key == nil {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		X     string `json:"x"`
		Y     string `json:"y"`
		Curve string `json:"curve"`
	}{
		X:     pk.Key.X.String(),
		Y:     pk.Key.Y.String(),
		Curve: "P256",
	})
}

func (pk *PublicKey) UnmarshalJSON(data []byte) error {
	var tmp struct {
		X     string `json:"x"`
		Y     string `json:"y"`
		Curve string `json:"curve"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	x := new(big.Int)
	x.SetString(tmp.X, 10)

	y := new(big.Int)
	y.SetString(tmp.Y, 10)

	var curve elliptic.Curve
	switch tmp.Curve {
	case "P256":
		curve = elliptic.P256()
	default:
		return fmt.Errorf("unknown curve")
	}

	pk.Key = &ecdsa.PublicKey{
		Curve: curve,
		X:     x,
		Y:     y,
	}

	return nil
}

func (s *Signature) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	return json.Marshal(struct {
		R string `json:"r"`
		S string `json:"s"`
	}{
		R: s.R().String(),
		S: s.S().String(),
	})
}

func (s *Signature) UnmarshalJSON(data []byte) error {
	var tmp struct {
		R string `json:"r"`
		S string `json:"s"`
	}

	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}

	s.R().SetString(tmp.R, 10)
	s.S().SetString(tmp.S, 10)

	return nil
}

func init() {
	gob.Register(&PublicKey{})
	gob.Register(&Signature{})
}
