package core

import (
	"uretra-network/crypto"
	"uretra-network/types"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	signature *crypto.Signature
	hash      types.Hash
	firstSeen int64
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

	tx.signature = sign
	tx.From = key.PublicKey()

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.signature == nil {
		return false
	}

	return tx.signature.VerifySignature(&tx.From, tx.Data)
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

func (tx *Transaction) SetFirstSeen(firstSeen int64) {
	tx.firstSeen = firstSeen
}

func (tx *Transaction) FirstSeen() int64 {
	return tx.firstSeen
}
