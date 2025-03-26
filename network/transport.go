package network

type NetAddress string

type Transport interface {
	Consume() <-chan RPC
	Connect(Transport) error
	SendMessage(NetAddress, []byte) error
	Broadcast([]byte) error
	Address() NetAddress
}
