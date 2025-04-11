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
	trLocal := network.NewTCPTransport(":3000", make(chan *network.TCPPeer))
	trRemoteA := network.NewTCPTransport(":3001", make(chan *network.TCPPeer))
	trRemoteB := network.NewTCPTransport(":3002", make(chan *network.TCPPeer))
	trRemoteC := network.NewTCPTransport(":3003", make(chan *network.TCPPeer))

	initRemoteServers(trRemoteA, trRemoteB, trRemoteC)

	//sendGetStatusMessage(trRemoteA, trLocal.Address())
	startLateServer()

	privateKey := crypto.GeneratePrivateKey()

	makeServer(&privateKey, "LOCAL", trLocal.ListenAddr, []string{":3001", ":3002"}).Start()
}

func startLateServer() {
	go func() {
		time.Sleep(7 * time.Second)

		trLate := network.NewTCPTransport(":4000", make(chan *network.TCPPeer))
		lateServer := makeServer(nil, "LATE", trLate.ListenAddr, []string{":3000"})

		go lateServer.Start()
	}()
}

func initRemoteServers(remoteA *network.TCPTransport, remoteB *network.TCPTransport, remoteC *network.TCPTransport) {
	s1 := makeServer(nil, "REMOTE_A", remoteA.ListenAddr, []string{":3002"})
	s2 := makeServer(nil, "REMOTE_B", remoteB.ListenAddr, []string{":3003"})
	s3 := makeServer(nil, "REMOTE_C", remoteC.ListenAddr, []string{":3004"})

	go s1.Start()
	go s2.Start()
	go s3.Start()
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

func sendTransaction(tr *network.TCPTransport) error {
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

	fmt.Println(msgData)

	return nil
}

func sendGetStatusMessage(tr *network.TCPTransport) {
	go func() {
		for {
			getStatusMsg := &network.GetStatusMessage{}

			buf := &bytes.Buffer{}
			_ = gob.NewEncoder(buf).Encode(getStatusMsg)

			msg := network.NewMessage(network.MessageTypeGetStatus, buf.Bytes())
			msgData, _ := msg.Bytes()

			fmt.Println(msgData)

			time.Sleep(time.Second * 2)
		}
	}()
}
