package core

import (
	"bytes"
	"encoding/gob"
	"uretra-network/crypto"
	"uretra-network/types"
)

type Header struct {
	Version       uint32
	PrevBlockHash types.Hash
	DataHash      types.Hash
	Timestamp     int64
	Height        uint32
}

type Block struct {
	Header       *Header
	Transactions []*Transaction
	Validator    crypto.PublicKey
	Signature    *crypto.Signature
	hash         types.Hash
}

func (h *Header) Bytes() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(h)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

func NewBlock(h *Header, tr []*Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: tr,
	}
}

func NewBlockWithPrivateKey(h *Header, tr []*Transaction, pk crypto.PrivateKey) *Block {
	b := &Block{
		Header:       h,
		Transactions: tr,
	}

	err := b.Sign(pk)

	if err != nil {
		panic(err)
	}

	return b
}

func (b *Block) AddTransaction(tr *Transaction) {
	b.Transactions = append(b.Transactions, tr)
}

func (b *Block) Sign(key crypto.PrivateKey) error {
	sign, err := key.Sign(b.Header.Bytes())

	if err != nil {
		return err
	}

	b.Validator = key.PublicKey()
	b.Signature = sign

	return nil
}

func (b *Block) Verify() bool {
	if b.Signature == nil {
		return false
	}

	for i := 0; i < len(b.Transactions); i++ {
		if !b.Transactions[i].Verify() {
			return false
		}
	}

	return b.Signature.VerifySignature(&b.Validator, b.Header.Bytes())
}

func (b *Block) Decode(decoder Decoder[*Block]) error {
	return decoder.Decode(b)
}

func (b *Block) Encode(encoder Encoder[*Block]) error {
	return encoder.Encode(b)
}

func (b *Block) Hash(hasher Hasher[*Header]) types.Hash {
	if b.hash.IsEmptyOrZero() {
		b.hash = hasher.Hash(b.Header)
	}

	return b.hash
}
