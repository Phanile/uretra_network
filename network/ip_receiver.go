package network

import (
	"fmt"
	"io"
	"net"
	"net/http"
)

func GetPublicIP() (string, error) {
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

func GetLocalIP() (string, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, errAddr := iface.Addrs()
		if errAddr != nil {
			return "", errAddr
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
