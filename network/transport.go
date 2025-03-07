package network

type NetAddress string

type RPC struct {
	From NetAddress
	Data []byte
}

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddress, []byte) error
	Address() NetAddress
}
