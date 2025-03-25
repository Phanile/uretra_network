package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io"
	"uretra-network/core"
)

type MessageType byte

const (
	MessageTypeTx MessageType = iota + 0x1
	MessageTypeBlock
)

type RPC struct {
	From    NetAddress
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

type RPCHandler interface {
	HandleRPC(rpc RPC) error
}

type DefaultRPCHandler struct {
	p RPCProcessor
}

func NewMessage(t MessageType, data []byte) *Message {
	return &Message{
		Header: t,
		Data:   data,
	}
}

func (m *Message) Bytes() ([]byte, error) {
	buf := &bytes.Buffer{}
	err := gob.NewEncoder(buf).Encode(m)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func NewDefaultRPCHandler(p RPCProcessor) *DefaultRPCHandler {
	return &DefaultRPCHandler{
		p: p,
	}
}

func (h *DefaultRPCHandler) HandleRPC(rpc RPC) error {
	msg := &Message{}
	err := gob.NewDecoder(rpc.Payload).Decode(msg)

	if err != nil {
		return fmt.Errorf("failed to decode RPC payload: %s", rpc.Payload)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := core.Transaction{}
		err := tx.Decode(core.NewGobTxDecoder(bytes.NewBuffer(msg.Data)))

		if err != nil {
			return err
		}

		return h.p.ProcessTransaction(rpc.From, &tx)

	default:
		return fmt.Errorf("invalid message type %x", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessTransaction(netaddr NetAddress, transaction *core.Transaction) error
}
