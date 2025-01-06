package codec

import (
	"testing"

	"github.com/stretchr/testify/assert"
	codecP "github.com/webx-top/codec"
)

type aes struct {
	codecP.Codec
}

func init() {
	RegisterStringCryptor([3]byte{'T', 'E', 'S'}, &aes{
		Codec: codecP.NewAES(`AES-256-CBC`),
	})
}

// EncryptString 加密
func (r *aes) EncryptString(input string) (string, error) {
	encrypted := r.Encode(input, `123`)
	return encrypted, nil
}

// DecryptString 解密
func (r *aes) DecryptString(input string) (string, error) {
	decrypted := r.Decode(input, `123`)
	return decrypted, nil
}

func TestStringCryptor(t *testing.T) {
	encrypted, err := AutoEncrypt(`TestStringCryptor`, `TES`)
	assert.NoError(t, err)
	assert.NotEmpty(t, encrypted)
	t.Log(encrypted)

	decrypted, err := AutoDecrypt(encrypted)
	assert.NoError(t, err)
	assert.Equal(t, `TestStringCryptor`, decrypted)
}
