package ipfilter

import (
	"fmt"
	"math/bits"
	"net"
	"net/netip"
	"strings"
)

// ParsePrefix converts an IP string into a netip.Prefix.
// If the input doesn't contain a CIDR notation, it automatically appends /32 for IPv4 or /128 for IPv6.
// Returns the parsed prefix or an error if the IP is invalid or parsing fails.
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

// ParseIPRange parses the given start and end IP addresses and returns a list of IP prefixes
// covering the range. Returns an error if either IP fails to parse.
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

// ToPrefixes converts an IP address range (start to end) into a list of CIDR prefixes.
// It handles both IPv4 and IPv6 addresses, ensuring the resulting prefixes cover the entire range.
// Returns an error if the IP types are inconsistent or if any calculation fails.
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

// prefixLength calculates the prefix length of an IP address by finding the first non-zero byte
// and counting the leading zeros in that byte. Returns 0 for the special case of overflow.
func prefixLength(ip net.IP) int {
	for index, c := range ip {
		if c != 0 {
			return index*8 + bits.LeadingZeros8(c) + 1
		}
	}
	// special case for overflow
	return 0
}

// trailingZeros counts the number of trailing zero bits in the IP address.
// It returns the total number of trailing zeros in the IP bytes,
// with each byte contributing up to 8 trailing zeros.
func trailingZeros(ip net.IP) int {
	ipLen := len(ip)
	for i := ipLen - 1; i >= 0; i-- {
		if c := ip[i]; c != 0 {
			return (ipLen-i-1)*8 + bits.TrailingZeros8(c)
		}
	}
	return ipLen * 8
}

// lastIP calculates the last IP address in a subnet given an IP and its netmask.
// Returns the broadcast address of the subnet or an error if IP and mask lengths mismatch.
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

// addOne increments the given IP address by one.
// It handles overflow by carrying over to the next byte when a byte reaches 0xFF.
// Returns a new IP address with the incremented value.
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

// xor performs a bitwise XOR operation between two IP addresses.
// Returns the resulting IP and an error if the input IPs have different lengths.
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
