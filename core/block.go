package core

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"io"
	"uretra-network/types"
)

type Header struct {
	version   uint32
	prevBlock types.Hash
	timestamp int64
	height    uint32
	nonce     uint64
}

func (h *Header) DecodeBinary(r io.Reader) error {
	if err := binary.Read(r, binary.LittleEndian, &h.version); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.prevBlock); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.timestamp); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.height); err != nil {
		return err
	}

	if err := binary.Read(r, binary.LittleEndian, &h.nonce); err != nil {
		return err
	}

	return nil
}

func (h *Header) EncodeBinary(w io.Writer) error {
	if err := binary.Write(w, binary.LittleEndian, &h.version); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.prevBlock); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.timestamp); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.height); err != nil {
		return err
	}

	if err := binary.Write(w, binary.LittleEndian, &h.nonce); err != nil {
		return err
	}

	return nil
}

type Block struct {
	header       Header
	transactions []Transaction
	hash         types.Hash
}

func (b *Block) Hash() types.Hash {
	buf := &bytes.Buffer{}
	err := b.header.EncodeBinary(buf)

	if err != nil {
		panic("Something went wrong while hashing the header of the block")
	}

	if b.hash.IsEmptyOrZero() {
		b.hash = sha256.Sum256(buf.Bytes())
	}

	return b.hash
}

func (b *Block) DecodeBinary(r io.Reader) error {
	err := b.header.DecodeBinary(r)

	if err != nil {
		return err
	}

	for _, tr := range b.transactions {
		trErr := tr.DecodeBinary(r)

		if trErr != nil {
			return trErr
		}
	}

	errRead := binary.Read(r, binary.LittleEndian, &b.hash)

	if errRead != nil {
		return errRead
	}

	return nil
}

func (b *Block) EncodeBinary(w io.Writer) error {
	err := b.header.EncodeBinary(w)

	if err != nil {
		return err
	}

	for _, tr := range b.transactions {
		trErr := tr.EncodeBinary(w)

		if trErr != nil {
			return trErr
		}
	}

	errWrite := binary.Write(w, binary.LittleEndian, &b.hash)

	if errWrite != nil {
		return errWrite
	}

	return nil
}
