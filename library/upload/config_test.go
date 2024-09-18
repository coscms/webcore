package upload

import (
	"mime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMime(t *testing.T) {
	cfg := NewConfig()
	result := cfg.DetectType(`.iso`)
	assert.Equal(t, `file`, result)
	result = cfg.DetectType(`.dmg`)
	assert.Equal(t, `file`, result)
	mimeType := mime.TypeByExtension(`.iso`)
	assert.Equal(t, `application/vnd.efi.iso`, mimeType) // application/x-cd-image
	mimeType = mime.TypeByExtension(`.dmg`)
	assert.Equal(t, `application/x-apple-diskimage`, mimeType)
}
