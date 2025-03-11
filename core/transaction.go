package core

import (
	"uretra-network/crypto"
)

type Transaction struct {
	Data      []byte
	PublicKey crypto.PublicKey
	signature *crypto.Signature
}

func (tx *Transaction) Sign(key crypto.PrivateKey) error {
	sign, err := key.Sign(tx.Data)

	if err != nil {
		return err
	}

	tx.signature = sign
	tx.PublicKey = key.PublicKey()

	return nil
}

func (tx *Transaction) Verify() bool {
	if tx.signature == nil {
		return false
	}

	return tx.signature.VerifySignature(&tx.PublicKey, tx.Data)
}
