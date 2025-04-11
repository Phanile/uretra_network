package main

import (
	"log"
	"uretra-network/crypto"
	"uretra-network/network"
)

func main() {
	privateKey := crypto.GeneratePrivateKey()
	makeServer(&privateKey, "LOCAL", ":3000", []string{":3001"}).Start()
}

func makeServer(pk *crypto.PrivateKey, id string, addr string, seedNodes []string) *network.Server {
	opts := network.ServerOptions{
		PrivateKey:    pk,
		ID:            id,
		SeedNodes:     seedNodes,
		ListenAddress: addr,
	}

	s, err := network.NewServer(&opts)

	if err != nil {
		log.Fatal(err)
	}

	return s
}
