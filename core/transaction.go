package core

import (
	"io"
	"uretra-network/network"
)

type Transaction struct {
	data []byte
	from network.NetAddress
}

func (tr *Transaction) EncodeBinary(w io.Writer) error {
	return nil
}

func (tr *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
