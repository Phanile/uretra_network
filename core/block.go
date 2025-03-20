package core

import (
	"bytes"
	"encoding/gob"
	"io"
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

func NewBlock(h *Header, tr []*Transaction) *Block {
	return &Block{
		Header:       h,
		Transactions: tr,
	}
}

func (b *Block) Sign(key crypto.PrivateKey) error {
	sign, err := key.Sign(b.HeaderData())

	if err != nil {
		return err
	}

	b.Validator = key.PublicKey()
	b.Signature = sign

	return nil
}

func (b *Block) Verify(signature crypto.Signature) bool {
	if b.Signature == nil {
		return false
	}

	return signature.VerifySignature(&b.Validator, b.HeaderData())
}

func (b *Block) Decode(r io.Reader, decoder Decoder[*Block]) error {
	return decoder.Decode(r, b)
}

func (b *Block) Encode(w io.Writer, encoder Encoder[*Block]) error {
	return encoder.Encode(w, b)
}

func (b *Block) Hash(hasher Hasher[*Block]) types.Hash {
	if b.hash.IsEmptyOrZero() {
		b.hash = hasher.Hash(b)
	}

	return b.hash
}

func (b *Block) HeaderData() []byte {
	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err := enc.Encode(b.Header)

	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}
