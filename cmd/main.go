package main

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"time"
	"uretra-network/core"
	"uretra-network/crypto"
	"uretra-network/network"
)

func main() {
	trLocal := network.NewLocalTransport("LOCAL")
	trRemoteA := network.NewLocalTransport("REMOTE_A")
	trRemoteB := network.NewLocalTransport("REMOTE_B")
	trRemoteC := network.NewLocalTransport("REMOTE_C")

	_ = trLocal.Connect(trRemoteA)
	_ = trRemoteA.Connect(trRemoteB)
	_ = trRemoteB.Connect(trRemoteC)

	_ = trRemoteA.Connect(trLocal)

	initRemoteServers([]network.Transport{trRemoteA, trRemoteB, trRemoteC})

	go func() {
		for {
			err := sendTransaction(trRemoteA, trLocal.Address())

			if err != nil {
				return
			}

			time.Sleep(2 * time.Second)
		}
	}()

	go func() {
		time.Sleep(7 * time.Second)

		trLate := network.NewLocalTransport("LATE_REMOTE")
		_ = trRemoteC.Connect(trLate)
		lateServer := makeServer(nil, "LATE", trLate)

		go lateServer.Start()
	}()

	privateKey := crypto.GeneratePrivateKey()

	makeServer(&privateKey, "LOCAL", trLocal).Start()
}

func initRemoteServers(trs []network.Transport) {
	for i := 0; i < len(trs); i++ {
		id := fmt.Sprintf("REMOTE_%d", i)
		s := makeServer(nil, id, trs[i])
		go s.Start()
	}
}

func makeServer(pk *crypto.PrivateKey, id string, transport network.Transport) *network.Server {
	opts := network.ServerOptions{
		PrivateKey: pk,
		ID:         id,
		Transports: []network.Transport{transport},
	}

	s, err := network.NewServer(&opts)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func sendTransaction(tr network.Transport, to network.NetAddress) error {
	privateKey := crypto.GeneratePrivateKey()
	data := []byte(strconv.FormatInt(int64(rand.Intn(1000000)), 10))
	tx := core.NewTransaction(data)
	err := tx.Sign(privateKey)

	if err != nil {
		return err
	}

	buf := &bytes.Buffer{}
	errEncode := tx.Encode(core.NewGobTxEncoder(buf))

	if errEncode != nil {
		return errEncode
	}

	msg := network.NewMessage(network.MessageTypeTx, buf.Bytes())

	msgData, errBytes := msg.Bytes()

	if errBytes != nil {
		return errBytes
	}

	errSend := tr.SendMessage(to, msgData)

	if errSend != nil {
		return errSend
	}

	return nil
}
