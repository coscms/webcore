package ipfilter

import (
	"net/netip"
	"strings"

	"github.com/admpub/bart"
	"github.com/admpub/log"
)

func New() *IPFilter {
	return &IPFilter{
		whitelist: &bart.Lite{},
		blacklist: &bart.Lite{},
	}
}

type IPFilter struct {
	whitelist *bart.Lite
	blacklist *bart.Lite
	whitesize uint
	blacksize uint
	disallow  bool
}

func (a *IPFilter) SetDisallow(dfl bool) *IPFilter {
	a.disallow = dfl
	return a
}

func (a *IPFilter) Allowed(ip string) error {
	return a.AddWhitelist(ip)
}

func (a *IPFilter) Blocked(ip string) error {
	return a.AddBlacklist(ip)
}

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

func (a *IPFilter) IsAllowed(realIP string) bool {
	addr, err := netip.ParseAddr(realIP)
	if err != nil {
		log.Warnf("failed to netip.ParseAddr(%q): %v", realIP, err)
		return false
	}
	return a.IsAllowedAddr(addr)
}

func (a *IPFilter) IsAllowedAddr(ip netip.Addr) bool {
	if a.whitesize > 0 {
		return a.whitelist.Contains(ip)
	}
	if a.blacksize > 0 {
		return !a.blacklist.Contains(ip)
	}
	return !a.disallow
}
