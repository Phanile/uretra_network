package main

import (
	"bytes"
	"fmt"
	"log"
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

	sendTestTransactions(trRemoteA, trLocal)
	//startLateServer(trRemoteC)

	privateKey := crypto.GeneratePrivateKey()

	makeServer(&privateKey, "LOCAL", trLocal).Start()
}

func sendTestTransactions(tr1, tr2 network.Transport) {
	go func() {
		for {
			_ = sendTransaction(tr1, tr2.Address())
			time.Sleep(2 * time.Second)
		}
	}()
}

func startLateServer(transport network.Transport) {
	go func() {
		time.Sleep(7 * time.Second)

		trLate := network.NewLocalTransport("LATE_REMOTE")
		_ = transport.Connect(trLate)
		lateServer := makeServer(nil, "LATE", trLate)

		go lateServer.Start()
	}()
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
	data := []byte{10, 0x01, 20, 0x01, 0x02} // 10 Push 20 Push Add
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
