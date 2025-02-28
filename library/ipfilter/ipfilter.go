package ipfilter

import (
	"net/netip"
	"strings"

	"github.com/admpub/bart"
	"github.com/admpub/log"
)

func New() *IPFilter {
	return &IPFilter{
		Whitelist: &bart.Table[struct{}]{},
		Blacklist: &bart.Table[struct{}]{},
	}
}

type IPFilter struct {
	Whitelist *bart.Table[struct{}]
	Blacklist *bart.Table[struct{}]
}

func (f *IPFilter) Allowed(ip string) error {
	return f.AddWhitelist(ip)
}

func (f *IPFilter) Blocked(ip string) error {
	return f.AddBlacklist(ip)
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
				a.Whitelist.Insert(pfx, struct{}{})
			}
			continue
		}
		pfx, err := netip.ParsePrefix(ip)
		if err != nil {
			return err
		}
		a.Whitelist.Insert(pfx, struct{}{})
	}
	return nil
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
				a.Whitelist.Insert(pfx, struct{}{})
			}
			continue
		}
		pfx, err := netip.ParsePrefix(ip)
		if err != nil {
			return err
		}
		a.Blacklist.Insert(pfx, struct{}{})
	}
	return nil
}

func (a *IPFilter) IsAllowed(realIP string) bool {
	ip, err := netip.ParseAddr(realIP)
	if err != nil {
		log.Warnf("failed to netip.ParseAddr(%q): %v", realIP, err)
		return false
	}
	if a.Whitelist.Size() > 0 {
		return a.Whitelist.Contains(ip)
	}
	if a.Blacklist.Size() > 0 {
		return !a.Blacklist.Contains(ip)
	}
	return true
}
