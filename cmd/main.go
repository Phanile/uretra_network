package main

import (
	"bytes"
	"fmt"
	"github.com/Phanile/uretra_network/core"
	"github.com/Phanile/uretra_network/crypto"
	"github.com/Phanile/uretra_network/network"
	"io"
	"log"
	"net/http"
	"time"
)

var (
	defaultListenPort    = ":3228"
	defaultAPIListenPort = ":3229"
)

func main() {
	makeServer().Start()
}

func makeServer() *network.Server {
	network.SetConfigToDefaultPeers()
	conf, errConf := network.GetConfig()

	if errConf != nil {
		panic("failed to load config")
	}

	ip, errIp := getPublicIP()

	if errIp != nil {
		panic("node is out network")
	}

	nodeId := "node_" + ip + defaultListenPort
	network.AddPeerToConfig(ip + defaultListenPort)

	privateKey := crypto.GeneratePrivateKey()

	opts := network.ServerOptions{
		APIListenAddress: ip + defaultAPIListenPort,
		PrivateKey:       &privateKey,
		ID:               nodeId,
		SeedNodes:        conf.Peers,
		ListenAddress:    ip + defaultListenPort,
		PeersConfig:      conf,
	}

	s, err := network.NewServer(&opts)

	if err != nil {
		log.Fatal(err)
	}

	return s
}

func getPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org?format=text")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	ip, errRead := io.ReadAll(resp.Body)

	if errRead != nil {
		return "", errRead
	}

	return string(ip), nil
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

		request, err := http.NewRequest("POST", "http://localhost:3229/tx", buf)

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
