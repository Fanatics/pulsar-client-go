package crypto

import (
	"testing"

	pb "github.com/apache/pulsar-client-go/pulsar/internal/pulsar_proto"
	"github.com/apache/pulsar-client-go/pulsar/log"
	"github.com/stretchr/testify/assert"
)

func TestAddPublicKeyCipher(t *testing.T) {
	msgCrypto, err := NewDefaultMessageCrypto("test-default-crypto", true, log.DefaultNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, msgCrypto)

	// valid keyreader
	err = msgCrypto.AddPublicKeyCipher(
		[]string{"my-app.key"},
		NewFileKeyReader("../crypto/pub_key_rsa.pem", ""),
	)
	assert.Nil(t, err)

	// invalid keyreader
	err = msgCrypto.AddPublicKeyCipher(
		[]string{"my-app0.key"},
		NewFileKeyReader("../crypto/no_pub_key_rsa.pem", ""),
	)
	assert.NotNil(t, err)

	// empty keyreader
	err = msgCrypto.AddPublicKeyCipher(
		[]string{"my-app1.key"},
		nil,
	)
	assert.NotNil(t, err)

	// keyreader with wrong econding of public key
	err = msgCrypto.AddPublicKeyCipher(
		[]string{"my-app2.key"},
		NewFileKeyReader("../crypto/wrong_encode_pub_key_rsa.pem", ""),
	)
	assert.NotNil(t, err)

	// keyreader with truncated pub key
	err = msgCrypto.AddPublicKeyCipher(
		[]string{"my-app2.key"},
		NewFileKeyReader("../crypto/truncated_pub_key_rsa.pem", ""),
	)
	assert.NotNil(t, err)
}

func TestEncrypt(t *testing.T) {
	msgMetadata := &pb.MessageMetadata{}
	msgMetadataSupplier := NewMessageMetadataSupplier(msgMetadata)

	msg := "my-message-01"

	msgCrypto, err := NewDefaultMessageCrypto("my-app", true, log.DefaultNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, msgCrypto)

	// valid keyreader
	encryptedData, err := msgCrypto.Encrypt(
		[]string{"my-app.key"},
		NewFileKeyReader("../crypto/pub_key_rsa.pem", ""),
		msgMetadataSupplier,
		[]byte(msg),
	)

	assert.Nil(t, err)
	assert.NotNil(t, encryptedData)

	// encrypted data key and encryption param must set in
	// in the message metadata after encryption
	assert.NotNil(t, msgMetadataSupplier.GetEncryptionParam())
	assert.NotEmpty(t, msgMetadataSupplier.GetEncryptionKeys())

	// invalid keyreader
	encryptedData, err = msgCrypto.Encrypt(
		[]string{"my-app2.key"},
		NewFileKeyReader("../crypto/no_pub_key_rsa.pem", ""),
		msgMetadataSupplier,
		[]byte(msg),
	)

	assert.NotNil(t, err)
	assert.Nil(t, encryptedData)
}

func TestEncryptDecrypt(t *testing.T) {
	msgMetadata := &pb.MessageMetadata{}
	msgMetadataSupplier := NewMessageMetadataSupplier(msgMetadata)

	msg := "my-message-01"

	msgCrypto, err := NewDefaultMessageCrypto("my-app", true, log.DefaultNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, msgCrypto)

	// valid keyreader
	encryptedData, err := msgCrypto.Encrypt(
		[]string{"my-app.key"},
		NewFileKeyReader("../crypto/pub_key_rsa.pem", ""),
		msgMetadataSupplier,
		[]byte(msg),
	)

	assert.Nil(t, err)
	assert.NotNil(t, encryptedData)

	// encrypted data key and encryption param must set in
	// in the message metadata after encryption
	assert.NotNil(t, msgMetadataSupplier.GetEncryptionParam())
	assert.NotEmpty(t, msgMetadataSupplier.GetEncryptionKeys())

	// try to decrypt
	msgCryptoDecrypt, err := NewDefaultMessageCrypto("my-app", true, log.DefaultNopLogger())
	assert.Nil(t, err)
	assert.NotNil(t, msgCrypto)

	// keyreader with invalid private key
	decryptedData, err := msgCryptoDecrypt.Decrypt(
		msgMetadataSupplier,
		encryptedData,
		NewFileKeyReader("", "../crypto/no_pri_key_rsa.pem"),
	)
	assert.NotNil(t, err)
	assert.Nil(t, decryptedData)

	// keyreader with wrong encoded private key
	decryptedData, err = msgCryptoDecrypt.Decrypt(
		msgMetadataSupplier,
		encryptedData,
		NewFileKeyReader("", "../crypto/wrong_encoded_pri_key_rsa.pem"),
	)
	assert.NotNil(t, err)
	assert.Nil(t, decryptedData)

	// keyreader with valid private key
	decryptedData, err = msgCryptoDecrypt.Decrypt(
		msgMetadataSupplier,
		encryptedData,
		NewFileKeyReader("", "../crypto/pri_key_rsa.pem"),
	)

	assert.Nil(t, err)
	assert.Equal(t, msg, string(decryptedData))
}
