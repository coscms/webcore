package settings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStructFieldConvert(t *testing.T) {
	assert.Equal(t, `[key][objkey]`, StructFieldConvert(`group[key][value][objkey]`))
	assert.Equal(t, `[key]`, StructFieldConvert(`group[key][value]`))
	assert.Equal(t, `group`, StructFieldConvert(`group`))
	assert.Equal(t, `[key]`, StructFieldConvert(`group[key]`))
}
