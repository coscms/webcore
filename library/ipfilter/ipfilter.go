package ipfilter

import (
	"net/netip"
	"strings"

	"github.com/admpub/bart"
	"github.com/admpub/log"
)

// New creates and returns a new IPFilter instance with initialized whitelist and blacklist.
func New() *IPFilter {
	return &IPFilter{
		whitelist: &bart.Lite{},
		blacklist: &bart.Lite{},
	}
}

// IPFilter represents a structure for managing IP address filtering.
// It contains whitelist and blacklist using bart.Lite for efficient IP range storage,
// along with counters for the size of each list and a disallow flag to control default behavior.
type IPFilter struct {
	whitelist *bart.Lite
	blacklist *bart.Lite
	whitesize uint
	blacksize uint
	disallow  bool
}

// SetDisallow sets whether the IP filter should disallow IPs by default.
// Returns the IPFilter instance for method chaining.
func (a *IPFilter) SetDisallow(dfl bool) *IPFilter {
	a.disallow = dfl
	return a
}

// Allowed adds the given IP address to the whitelist, effectively allowing it through the filter.
func (a *IPFilter) Allowed(ip string) error {
	return a.AddWhitelist(ip)
}

// Blocked adds the specified IP address to the blacklist.
func (a *IPFilter) Blocked(ip string) error {
	return a.AddBlacklist(ip)
}

// AddWhitelist adds one or more IP addresses or ranges to the whitelist.
// Each IP can be a single address (e.g. "192.168.1.1") or a range (e.g. "192.168.1.1-192.168.1.100").
// Returns an error if any IP address is invalid or cannot be parsed.
func (a *IPFilter) AddWhitelist(ips ...string) error {
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if len(ip) == 0 {
			continue
		}
		parts := strings.SplitN(ip, `-`, 2)
		if len(parts) == 2 {
			pfxList, err := ParseIPRange(parts[0], parts[1])
			if err != nil {
				return err
			}
			for _, pfx := range pfxList {
				a.insertWhitelist(pfx)
			}
			continue
		}
		pfx, err := ParsePrefix(ip)
		if err != nil {
			return err
		}
		a.insertWhitelist(pfx)
	}
	return nil
}

func (a *IPFilter) insertWhitelist(pfx netip.Prefix) {
	a.whitesize++
	a.whitelist.Insert(pfx)
}

// AddBlacklist adds one or more IP addresses or ranges to the blacklist.
// Each IP can be a single address (e.g., "192.168.1.1") or a range (e.g., "192.168.1.1-192.168.1.100").
// Returns an error if any IP address or range is invalid.
func (a *IPFilter) AddBlacklist(ips ...string) error {
	for _, ip := range ips {
		ip = strings.TrimSpace(ip)
		if len(ip) == 0 {
			continue
		}
		parts := strings.SplitN(ip, `-`, 2)
		if len(parts) == 2 {
			pfxList, err := ParseIPRange(parts[0], parts[1])
			if err != nil {
				return err
			}
			for _, pfx := range pfxList {
				a.insertBlacklist(pfx)
			}
			continue
		}
		pfx, err := ParsePrefix(ip)
		if err != nil {
			return err
		}
		a.insertBlacklist(pfx)
	}
	return nil
}

func (a *IPFilter) insertBlacklist(pfx netip.Prefix) {
	a.blacksize++
	a.blacklist.Insert(pfx)
}

// IsAllowed checks if the given IP address string is allowed by the filter.
// Returns false if the IP address is invalid or not allowed.
func (a *IPFilter) IsAllowed(realIP string) bool {
	addr, err := netip.ParseAddr(realIP)
	if err != nil {
		log.Warnf("failed to netip.ParseAddr(%q): %v", realIP, err)
		return false
	}
	return a.IsAllowedAddr(addr)
}

// IsAllowedAddr checks if the given IP address is allowed based on the filter's whitelist, blacklist, and disallow settings.
// Returns true if the IP is allowed, false otherwise.
func (a *IPFilter) IsAllowedAddr(ip netip.Addr) bool {
	if a.whitesize > 0 {
		return a.whitelist.Contains(ip)
	}
	if a.blacksize > 0 {
		return !a.blacklist.Contains(ip)
	}
	return !a.disallow
}
