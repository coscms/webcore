package wasmplugin

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWASM(t *testing.T) {
	files, _ := filepath.Glob(`./testdata/*.wasm`)
	err := Load(files...)
	assert.NoError(t, err)
}
