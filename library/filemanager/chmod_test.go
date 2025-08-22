package filemanager

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestChmod(t *testing.T) {
	v := fmt.Sprintf(`0%d%d%d`, 7, 7, 7)
	t.Logf(`===%v`, v)
	n, _ := strconv.ParseUint(v, 8, 32)
	t.Logf(`===%v`, n)

	//num, _ := strconv.ParseUint(fmt.Sprintf(`%d`, n), 8, 32)
	t.Logf(`===%v`, strconv.FormatUint(n, 8))
	t.Logf(`===%+v`, FileModeToPerms(os.ModePerm))

	fi, err := os.Stat(`..`)
	if err == nil {
		md := fi.Mode()
		str := strconv.FormatUint(uint64(md), 8)
		t.Logf(`===%v %s %+v`, str, md, FileModeToPerms(md))
	}

	assert.True(t, ValidatePermCodes([3]uint32{4, 0, 0}))
	assert.True(t, ValidatePermNumber(4))
}
