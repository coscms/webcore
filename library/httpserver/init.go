package httpserver

import (
	"github.com/coscms/webcore/library/common"
)

func init() {
	common.SetProcessError(ProcessError)
}
