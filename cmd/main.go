package main

import (
	"bytes"
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
	trRemote := network.NewLocalTransport("REMOTE")

	trLocal.Connect(trRemote)
	trRemote.Connect(trLocal)

	go func() {
		for {
			err := sendTransaction(trRemote, trLocal.Address())

			if err != nil {
				return
			}

			time.Sleep(1 * time.Second)
		}
	}()

	privateKey := crypto.GeneratePrivateKey()

	opts := network.ServerOptions{
		PrivateKey: &privateKey,
		ID:         "LOCAL",
		Transports: []network.Transport{trLocal},
	}

	s, err := network.NewServer(&opts)

	if err != nil {
		log.Fatal(err)
	}

	s.Start()
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
