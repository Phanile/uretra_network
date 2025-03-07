package main

import (
	"time"
	"uretra-network/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			trRemote.SendMessage(trLocal.Address(), []byte("HELLO FROM REMOTE"))
			time.Sleep(1 * time.Second)
		}
	}()

	opts := network.ServerOptions{
		Transports: []network.Transport{trLocal},
	}

	s := network.NewServer(&opts)
	s.Start()
}
