package crypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicKey(t *testing.T) {
	keyReader := NewFileKeyReader("../crypto/pub-key-rsa.pem", "")
	keyInfo, err := keyReader.GetPublicKey("test-key", map[string]string{"key": "value"})

	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	assert.NotEmpty(t, keyInfo.GetMetadata())
	assert.NotEmpty(t, keyInfo.GetKey())
	assert.NotEmpty(t, keyInfo.GetValue())
	assert.Equal(t, "value", keyInfo.metadata["key"])
}

func TestGetPrivateKey(t *testing.T) {
	keyReader := NewFileKeyReader("", "../crypto/pri-key-rsa.pem")
	keyInfo, err := keyReader.GetPrivateKey("test-key", map[string]string{"key": "value"})

	assert.Nil(t, err)
	assert.NotNil(t, keyInfo)
	assert.NotEmpty(t, keyInfo.GetMetadata())
	assert.NotEmpty(t, keyInfo.GetKey())
	assert.NotEmpty(t, keyInfo.GetValue())
	assert.Equal(t, "value", keyInfo.metadata["key"])
}

func TestInvalidKeyPath(t *testing.T) {
	keyReader := NewFileKeyReader("../crypto/invalid-pub-key-rsa.pem", "../crypto/invalid-pri-key-rsa.pem")

	// try to read public key
	keyInfo, err := keyReader.GetPublicKey("test-pub-key", nil)
	assert.Nil(t, keyInfo)
	assert.NotNil(t, err)

	// try to read private key
	keyInfo, err = keyReader.GetPrivateKey("test-pri-key", nil)
	assert.Nil(t, keyInfo)
	assert.NotNil(t, err)
}
