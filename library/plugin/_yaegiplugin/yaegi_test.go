package yaegiplugin

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	err := Load(`./_testdata/main.go`)
	assert.NoError(t, err)
}
