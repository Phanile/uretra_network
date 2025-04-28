package core

import (
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
)

type Transaction struct {
	Data      []byte
	From      crypto.PublicKey
	To        types.Address
	Value     uint64
	Nonce     uint64
	Signature *crypto.Signature
	hash      types.Hash
}

func NewTransaction(
	data []byte,
	from crypto.PublicKey,
	to types.Address,
	value uint64,
	nonce uint64) *Transaction {
	return &Transaction{
		Data:  data,
		From:  from,
		To:    to,
		Value: value,
		Nonce: nonce,
	}
}

func (tx *Transaction) Sign(key crypto.PrivateKey) error {
	txHash := tx.Hash(TxHasher{})

	sign, err := key.Sign(txHash[:])

	if err != nil {
		return err
	}

	tx.Signature = sign
	tx.From = key.PublicKey()

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.From == crypto.ZeroPublicKey() && tx.To == crypto.ZeroPublicKey().Address() { //coinbase transactions
		return true
	}

	if tx.From.Address() == tx.To {
		return false
	}

	if tx.Signature == nil {
		return false
	}

	txHash := tx.Hash(TxHasher{})

	return tx.Signature.VerifySignature(&tx.From, txHash[:])
}

func (tx *Transaction) Hash(hasher Hasher[*Transaction]) types.Hash {
	if tx.hash.IsEmptyOrZero() {
		return hasher.Hash(tx)
	}

	return tx.hash
}

func (tx *Transaction) Decode(decoder Decoder[*Transaction]) error {
	return decoder.Decode(tx)
}

func (tx *Transaction) Encode(encoder Encoder[*Transaction]) error {
	return encoder.Encode(tx)
}
