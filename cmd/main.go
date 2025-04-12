package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
	"uretra-network/network"
)

func main() {
	privateKey := crypto.GeneratePrivateKey()
	go sendTestTransactions()
	makeServer(&privateKey, "LOCAL", ":3000", []string{":3001"}, ":3228").Start()
}

func makeServer(pk *crypto.PrivateKey, id string, addr string, seedNodes []string, APIListenAddr string) *network.Server {
	opts := network.ServerOptions{
		APIListenAddress: APIListenAddr,
		PrivateKey:       pk,
		ID:               id,
		SeedNodes:        seedNodes,
		ListenAddress:    addr,
	}

	s, err := network.NewServer(&opts)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func sendTestTransactions() {
	for {
		time.Sleep(time.Second * 3)
		privateKey := crypto.GeneratePrivateKey()
		data := []byte{8, 0x01, 105, 0x03, 116, 0x03, 32, 0x03, 119, 0x03, 111, 0x03, 114, 0x03, 107, 0x03, 115, 0x03, 0x04, 21, 0x01, 0x06}
		tx := core.NewTransaction(data)
		_ = tx.Sign(privateKey)

		buf := &bytes.Buffer{}
		_ = tx.Encode(core.NewGobTxEncoder(buf))

		request, err := http.NewRequest("POST", "http://localhost:3228/tx", buf)

		if err != nil {
			fmt.Println("error while request tx to json server")
		}

		client := http.Client{}
		_, errReq := client.Do(request)

		if errReq != nil {
			panic(errReq)
		}
	}
}
