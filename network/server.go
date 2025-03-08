package network

import (
	"fmt"
	"time"
)

type ServerOptions struct {
	Transports []Transport
}

type Server struct {
	so          *ServerOptions
	rpcChannel  chan RPC
	quitChannel chan struct{}
}

func NewServer(opts *ServerOptions) *Server {
	return &Server{
		so:          opts,
		rpcChannel:  make(chan RPC),
		quitChannel: make(chan struct{}),
	}
}

func (s *Server) Start() {
	s.initTransports()
	ticker := time.NewTicker(5 * time.Second)

free:
	for {
		select {
		case rpc := <-s.rpcChannel:
			fmt.Printf("%+v\n", rpc)
		case <-s.quitChannel:
			break free
		case <-ticker.C:
			fmt.Println("do stuff every 5 seconds")
		}
	}

	fmt.Println("Server shutdown")
}

func (s *Server) initTransports() {
	for _, transport := range s.so.Transports {
		go func(transport Transport) {
			for rpc := range transport.Consume() {
				s.rpcChannel <- rpc
			}
		}(transport)
	}
}
