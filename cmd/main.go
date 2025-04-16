package main

import (
	"fmt"
	"log"
	"net"
	"uretra-network/crypto"
	"uretra-network/network"
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

	ip, errIp := getLocalIP()

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

func getLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			ipNet, ok := addr.(*net.IPNet)
			if !ok || ipNet.IP.IsLoopback() || ipNet.IP.To4() == nil {
				continue
			}

			return ipNet.IP.String(), nil
		}
	}

	return "", fmt.Errorf("local ip address not found")
}
