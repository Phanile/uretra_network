package network

import (
	"encoding/json"
	"flag"
	"os"
)

var configPath string = getConfigPath()

type PeersConfig struct {
	Peers []string `json:"peers"`
}

func AddPeerToConfig(peerAddr string) {
	conf, err := GetConfig()

	if err != nil {
		panic(err)
	}

	for _, peer := range conf.Peers {
		if peer == peerAddr {
			return
		}
	}

	conf.Peers = append(conf.Peers, peerAddr)

	SaveConfig(conf)
}

func RemovePeerFromConfig(peerAddr string) {
	conf, err := GetConfig()

	if err != nil {
		panic(err)
	}

	if len(conf.Peers) == 0 {
		return
	}

	index := -1
	for i, addr := range conf.Peers {
		if addr == peerAddr {
			index = i
			break
		}
	}

	if index != -1 {
		conf.Peers = append(conf.Peers[:index], conf.Peers[index+1:]...)
		SaveConfig(conf)
	}
}

func GetConfig() (*PeersConfig, error) {
	data, err := os.ReadFile(configPath)

	if err != nil {
		panic(err)
	}

	var config PeersConfig

	errUnmarshal := json.Unmarshal(data, &config)

	if errUnmarshal != nil {
		panic(errUnmarshal)
	}

	return &config, nil
}

func getConfigPath() string {
	var result string

	flag.StringVar(&result, "config", "", "config file path")
	flag.Parse()

	if result == "" {
		result = os.Getenv("CONFIG_PATH")
	}

	return result
}

func SaveConfig(conf *PeersConfig) {
	data, err := json.MarshalIndent(conf, "", "  ")

	if err != nil {
		panic(err)
	}

	errWrite := os.WriteFile(configPath, data, 0644)

	if errWrite != nil {
		panic(errWrite)
	}
}

func SetConfigToDefaultPeers() {
	conf, err := GetConfig()

	if err != nil {
		panic(err)
	}

	conf.Peers = []string{"89.151.179.176:3228"}

	SaveConfig(conf)
}
