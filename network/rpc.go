package network

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"github.com/Phanile/uretra_network/core"
	"io"
	"net"
)

type MessageType byte

const (
	MessageTypeTx MessageType = iota + 0x1
	MessageTypeBlock
	MessageTypeGetBlocks
	MessageTypeStatus
	MessageTypeGetStatus
	MessageTypeBlocks
	MessageTypePing
	MessageTypePong
)

type RPC struct {
	From    net.Addr
	Payload io.Reader
}

type Message struct {
	Header MessageType
	Data   []byte
}

type DecodedMessage struct {
	From net.Addr
	Data any
}

type RPCDecodeFunc func(RPC) (*DecodedMessage, error)

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

	case MessageTypeStatus:
		statusMsg := &StatusMessage{}

		err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(statusMsg)

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: statusMsg,
		}, nil

	case MessageTypeGetStatus:
		return &DecodedMessage{
			From: rpc.From,
			Data: &GetStatusMessage{},
		}, nil

	case MessageTypeGetBlocks:
		getBlocksMsg := &GetBlocksMessage{}

		err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(getBlocksMsg)

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: getBlocksMsg,
		}, nil

	case MessageTypeBlocks:
		blocksMsg := &BlocksMessage{}
		err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(blocksMsg)

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: blocksMsg,
		}, nil

	case MessageTypePing:
		pingMsg := &PingMessage{}

		err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(pingMsg)

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: pingMsg,
		}, nil

	case MessageTypePong:
		pongMsg := &PongMessage{}

		err := gob.NewDecoder(bytes.NewReader(msg.Data)).Decode(pongMsg)

		if err != nil {
			return nil, err
		}

		return &DecodedMessage{
			From: rpc.From,
			Data: pongMsg,
		}, nil

	default:
		return nil, fmt.Errorf("invalid message type %x", msg.Header)
	}
}

type RPCProcessor interface {
	ProcessMessage(m *DecodedMessage) error
}
