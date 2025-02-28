package ipfilter

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/param"
)

func isNotEmptyString(s *string) bool {
	if s == nil {
		return false
	}
	*s = strings.TrimSpace(*s)
	return len(*s) > 0
}

func ValidateRows(ctx echo.Context, iplist string) error {
	ips := param.Split(iplist, "\n").Filter(isNotEmptyString).Unique().String()
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

func ValidateRange(startIP, endIP string) error {
	start, err := netip.ParseAddr(startIP)
	if err != nil {
		return fmt.Errorf(`%w(%q): %w`, ErrParseStartIPAddress, startIP, err)
	}
	var end netip.Addr
	end, err = netip.ParseAddr(endIP)
	if err != nil {
		return fmt.Errorf(`%w(%q): %w`, ErrParseEndIPAddress, endIP, err)
	}
	if start.BitLen() != end.BitLen() {
		return fmt.Errorf(`%w: %v - %v`, ErrStartAndEndIPMismatchType, start.String(), end.String())
	}
	return err
}
