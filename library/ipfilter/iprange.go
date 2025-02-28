package ipfilter

import (
	"fmt"
	"math/bits"
	"net"
	"net/netip"
	"strings"
)

func ParsePrefix(ip string) (pfx netip.Prefix, err error) {
	if !strings.Contains(ip, `/`) {
		ipr := net.ParseIP(ip)
		if ipr == nil {
			err = fmt.Errorf(`invalid ip: %v`, ip)
			return
		}
		if ipr.To4() != nil {
			ip += `/32`
		} else {
			ip += `/128`
		}
	}
	pfx, err = netip.ParsePrefix(ip)
	if err != nil {
		err = fmt.Errorf(`failed to parse ip(%q): %w`, ip, err)
	}
	return
}

func ParseIPRange(startIP, endIP string) ([]netip.Prefix, error) {
	start, err := netip.ParseAddr(startIP)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ip(%q): %w`, startIP, err)
	}
	end, err := netip.ParseAddr(endIP)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse ip(%q): %w`, endIP, err)
	}
	return ToPrefixes(start, end)
}

func ToPrefixes(start netip.Addr, end netip.Addr) ([]netip.Prefix, error) {
	if start.BitLen() != end.BitLen() {
		return nil, fmt.Errorf(`inconsistency between start and end ip types: %v - %v`, start.String(), end.String())
	}
	if !start.Less(end) {
		start, end = end, start
	}
	var result []netip.Prefix
	for {
		startIP := start.AsSlice()
		prx, err := xor(addOne(end.AsSlice()), startIP)
		if err != nil {
			return nil, err
		}
		cidr := max(prefixLength(prx), start.BitLen()-trailingZeros(startIP))
		mask := net.CIDRMask(cidr, start.BitLen())
		prefix := netip.PrefixFrom(start, cidr)
		result = append(result, prefix)
		tmp, err := lastIP(startIP, mask)
		if err != nil {
			return nil, err
		}
		if start.BitLen() == 128 {
			start = netip.AddrFrom16([16]byte(addOne(tmp)))
		} else {
			start = netip.AddrFrom4([4]byte(addOne(tmp)))
		}
		if !start.Less(end) {
			return result, nil
		}
	}
}

func prefixLength(ip net.IP) int {
	for index, c := range ip {
		if c != 0 {
			return index*8 + bits.LeadingZeros8(c) + 1
		}
	}
	// special case for overflow
	return 0
}

func trailingZeros(ip net.IP) int {
	ipLen := len(ip)
	for i := ipLen - 1; i >= 0; i-- {
		if c := ip[i]; c != 0 {
			return (ipLen-i-1)*8 + bits.TrailingZeros8(c)
		}
	}
	return ipLen * 8
}

func lastIP(ip net.IP, mask net.IPMask) (net.IP, error) {
	ipLen := len(ip)
	if ipLen != len(mask) {
		return nil, fmt.Errorf("unexpected IPNet %v", ip.String()+`/`+mask.String())
	}
	res := make(net.IP, ipLen)
	for i, b := range mask {
		res[i] = ip[i] | ^b
	}
	return res, nil
}

func addOne(ip net.IP) net.IP {
	ipLen := len(ip)
	res := make(net.IP, ipLen)
	for i := ipLen - 1; i >= 0; i-- {
		if t := ip[i]; t != 0xFF {
			res[i] = t + 1
			copy(res, ip[0:i])
			break
		}
	}
	return res
}

func xor(a, b net.IP) (net.IP, error) {
	ipLen := len(a)
	if ipLen != len(b) {
		return nil, fmt.Errorf("a=%v, b=%v", a, b)
	}
	res := make(net.IP, ipLen)
	for i := ipLen - 1; i >= 0; i-- {
		res[i] = a[i] ^ b[i]
	}
	return res, nil
}
