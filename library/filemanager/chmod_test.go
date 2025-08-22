package filemanager

import (
	"fmt"
	"strconv"
	"testing"
)

func TestChmod(t *testing.T) {
	v := fmt.Sprintf(`0%d%d%d`, 7, 7, 7)
	t.Logf(`===%v`, v)
	n, _ := strconv.ParseUint(v, 8, 32)
	t.Logf(`===%v`, n)
}
