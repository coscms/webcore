package request

import (
	"testing"

	"github.com/coscms/webcore/library/codec"
	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	pass, _ := codec.DefaultSM2EncryptHex(`12345678`)
	rawPwd, _ := codec.DefaultSM2DecryptHex(pass)
	assert.Equal(t, rawPwd, `12345678`)
	data := &Login{
		User: `test`,
		Pass: pass,
		Code: `12345678`,
	}
	err := echoContext.Validate(data)
	assert.NoError(t, err)
}
