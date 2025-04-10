package main

import (
	"bytes"
	"encoding/gob"
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
	//sendGetStatusMessage(trRemoteA, trLocal.Address())
	startLateServer(trLocal)

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
		_ = trLate.Connect(transport)
		lateServer := makeServer(nil, "LATE", trLate)

		go lateServer.Start()

		sendGetStatusMessage(trLate, transport.Address())
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
		Transport:  transport,
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
	//8, PushInt, i, PushBytes, t, PushBytes, space, PushBytes, w, PushBytes, o, PushBytes, r, PushBytes, k, PushBytes, s, PushBytes, Pack, 21, PushInt, Store
	data := []byte{8, 0x01, 105, 0x03, 116, 0x03, 32, 0x03, 119, 0x03, 111, 0x03, 114, 0x03, 107, 0x03, 115, 0x03, 0x04, 21, 0x01, 0x06}
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

func sendGetStatusMessage(transport network.Transport, to network.NetAddress) {
	go func() {
		for {
			getStatusMsg := &network.GetStatusMessage{}

			buf := &bytes.Buffer{}
			_ = gob.NewEncoder(buf).Encode(getStatusMsg)

			msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())
			msgData, _ := msg.Bytes()

			_ = transport.SendMessage(to, msgData)

			time.Sleep(time.Second * 2)
		}
	}()
}
