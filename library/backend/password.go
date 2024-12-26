package backend

import (
	"github.com/webx-top/echo"

	"github.com/coscms/webcore/library/codec"
	"github.com/coscms/webcore/library/sessionguard"
)

func DecryptPassword(c echo.Context, pass string) (string, error) {
	var err error
	pass, err = codec.DefaultSM2DecryptHex(pass)
	if err != nil {
		return pass, err
	}
	pass, err = sessionguard.Unpack(c, pass, true)
	return pass, err
}
