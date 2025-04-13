package network

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPeersConfig_Get(t *testing.T) {
	config, err := GetConfig()
	fmt.Println(config)
	assert.Nil(t, err)
	assert.NotNil(t, config)
}

func TestPeersConfig_Save(t *testing.T) {
	config, err := GetConfig()
	fmt.Println(config)
	assert.Nil(t, err)
	assert.NotNil(t, config)

	updatedConfig := PeersConfig{
		Peers: []string{"89.151.159.224:3228", "11.111.111.223:3228", "162.163.164.165:189"},
	}

	SaveConfig(&updatedConfig)
}
