package crypto

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestKeypair_Sign_Verify(t *testing.T) {
	privateKey := GeneratePrivateKey()
	publicKey := privateKey.PublicKey()

	randomPrivateKey := GeneratePrivateKey()
	randomPublicKey := randomPrivateKey.PublicKey()

	msg := []byte("Test message. Transaction 1 approved by blockchain")

	sign, err := privateKey.Sign(msg)
	assert.Nil(t, err)
	assert.True(t, sign.VerifySignature(&publicKey, msg))
	assert.False(t, sign.VerifySignature(&publicKey, []byte("Random data")))
	assert.False(t, sign.VerifySignature(&randomPublicKey, msg))
}
