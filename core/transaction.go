package core

import (
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	To        crypto.PublicKey
	Value     uint64
	Signature *crypto.Signature
	hash      types.Hash
}

func NewTransaction(data []byte) *Transaction {
	return &Transaction{
		Data: data,
	}
}

func (tx *Transaction) Sign(key crypto.PrivateKey) error {
	sign, err := key.Sign(tx.Data)

	if err != nil {
		return err
	}

	tx.Signature = sign
	tx.From = key.PublicKey()

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.Signature == nil {
		return false
	}

	return tx.Signature.VerifySignature(&tx.From, tx.Data)
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsEmptyOrZero() {
		tx.hash = hasher.Hash(tx)
	}

	return tx.hash
}

func (tx *Transaction) Decode(decoder Decoder[*Transaction]) error {
	return decoder.Decode(tx)
}

func (tx *Transaction) Encode(encoder Encoder[*Transaction]) error {
	return encoder.Encode(tx)
}
