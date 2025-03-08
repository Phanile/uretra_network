package core

import "io"

type Transaction struct {
	Data []byte
}

func (tr *Transaction) EncodeBinary(w io.Writer) error {
	return nil
}

func (tr *Transaction) DecodeBinary(r io.Reader) error {
	return nil
}
