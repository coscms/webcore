package model

import (
	"strings"
	"testing"

	"github.com/coscms/webcore/dbschema"
	"github.com/stretchr/testify/assert"
)

func TestTrimOverflowText(t *testing.T) {
	v := strings.Repeat(`123`, 300)
	r := dbschema.DBI.Fields.TrimOverflowText(`nging_sending_log`, `result`, v)
	assert.Equal(t, v[0:252]+`...`, r)
}
