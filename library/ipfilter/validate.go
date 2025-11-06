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

// isNotEmptyString checks if the given string pointer is not nil and contains non-whitespace characters.
// It trims whitespace from the string and returns true if the resulting string is not empty.
// Returns false if the pointer is nil or the string is empty after trimming.
func isNotEmptyString(s *string) bool {
	if s == nil {
		return false
	}
	*s = strings.TrimSpace(*s)
	return len(*s) > 0
}

// ValidateRows validates a list of IP addresses separated by newlines.
// It checks each IP for validity and returns appropriate error messages if any IP is invalid.
// Returns nil if all IPs are valid, otherwise returns an echo.HTTPError with details about the first invalid IP encountered.
// The error messages distinguish between mismatched IP types, invalid start IP, invalid end IP, and general invalid IP cases.
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

// Validate checks if the given IP string is valid, supporting both single IP addresses
// and IP ranges (in the format "start-end"). Returns an error if the IP is invalid.
func Validate(ip string) error {
	parts := strings.SplitN(ip, `-`, 2)
	if len(parts) == 2 {
		return ValidateRange(parts[0], parts[1])
	}
	_, err := ParsePrefix(ip)
	return err
}

// ValidateRange validates that the given start and end IP addresses form a valid IP range.
// It checks that both IPs are valid and of the same type (IPv4 or IPv6).
// Returns an error if either IP is invalid or if they are of different types.
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
