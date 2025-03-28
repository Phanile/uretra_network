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
	MessageTypeTx MessageType = iota
	MessageTypeBlock
	MessageTypeGetBlocks
)

type RPC struct {
	From    NetAddress
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

type DecodedMessage struct {
	From NetAddress
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

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

func DefaultRPCDecodeFunc(rpc RPC) (*DecodedMessage, error) {
	msg := &Message{}
	err := gob.NewDecoder(rpc.Payload).Decode(msg)

	if err != nil {
		return nil, fmt.Errorf("failed to decode RPC payload: %s", rpc.Payload)
	}

	switch msg.Header {
	case MessageTypeTx:
		tx := &core.Transaction{}
		err := tx.Decode(core.NewGobTxDecoder(bytes.NewReader(msg.Data)))

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: tx,
		}, nil
	case MessageTypeBlock:
		b := &core.Block{}
		err := b.Decode(core.NewGobBlockDecoder(bytes.NewReader(msg.Data)))

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: b,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message type %x", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(m *DecodedMessage) error
}
