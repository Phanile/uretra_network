package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/types"
	"time"
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

func NewBlockFromPrevHeader(prevHeader *Header, tx []*Transaction) (*Block, error) {
	dataHash, err := CalculateDataHash(tx)

	if err != nil {
		return nil, err
	}

	header := &Header{
		Version:       1,
		DataHash:      dataHash,
		PrevBlockHash: HeaderHasher{}.Hash(prevHeader),
		Timestamp:     time.Now().UnixNano(),
		Height:        prevHeader.Height + 1,
	}

	return NewBlock(header, tx), nil
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

	dataHash, _ := b.CalculateDataHash(b.Transactions)

	if dataHash != b.Header.DataHash {
		return false
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

func (b *Block) CalculateDataHash(tx []*Transaction) (types.Hash, error) {
	buf := &bytes.Buffer{}

	for _, tr := range tx {
		err := tr.Encode(NewGobTxEncoder(buf))

		if err != nil {
			return types.Hash{}, err
		}
	}

	return sha256.Sum256(buf.Bytes()), nil
}

func CalculateDataHash(tx []*Transaction) (types.Hash, error) {
	buf := &bytes.Buffer{}

	for _, tr := range tx {
		err := tr.Encode(NewGobTxEncoder(buf))

		if err != nil {
			return types.Hash{}, err
		}
	}

	return sha256.Sum256(buf.Bytes()), nil
}
