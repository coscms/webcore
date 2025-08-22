package filemanager

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestChmod(t *testing.T) {
	v := fmt.Sprintf(`0%d%d%d`, 7, 7, 7)
	t.Logf(`===%v`, v)
	n, _ := strconv.ParseUint(v, 8, 32)
	t.Logf(`===%v`, n)

	//num, _ := strconv.ParseUint(fmt.Sprintf(`%d`, n), 8, 32)
	t.Logf(`===%v`, strconv.FormatUint(n, 8))
	t.Logf(`===%+v`, FileModeToPerms(os.ModePerm))
}
