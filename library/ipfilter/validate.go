package ipfilter

import (
	"errors"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func ValidateRows(ctx echo.Context, iplist string) error {
	ips := param.Split(iplist, "\n").Filter().Unique().String()
	var err error
	for _, ip := range ips {
		err = Validate(ip)
		if err != nil {
			if errors.Is(err, ErrStartAndEndIPMismatchType) {
				return ctx.NewError(code.InvalidParameter, `起始和结束IP不能使用不同的IP类型: %s`, ip)
			}
			if errors.Is(err, ErrParseStartIPAddress) {
				return ctx.NewError(code.InvalidParameter, `起始IP的类型不正确: %s`, ip)
			}
			if errors.Is(err, ErrParseEndIPAddress) {
				return ctx.NewError(code.InvalidParameter, `结束IP的类型不正确: %s`, ip)
			}
			return ctx.NewError(code.InvalidParameter, `IP地址无效: %s`, ip)
		}
	}
	return err
}

func Validate(ip string) error {
	parts := strings.SplitN(ip, `-`, 2)
	if len(parts) == 2 {
		return ValidateRange(parts[0], parts[1])
	}
	_, err := ParsePrefix(ip)
	if err != nil {
		return err
	}
	return err
}
