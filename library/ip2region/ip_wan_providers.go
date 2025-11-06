package ip2region

import (
	"regexp"

	"github.com/coscms/webcore/library/config/extend"
)

// WANIPProvider defines a structure for WAN IP address providers.
// It contains configuration for retrieving public IP addresses (both IPv4 and IPv6)
// from external services, including request details and response parsing rules.
// Name: Provider's display name
// Description: Brief description of the provider
// URL: Endpoint URL for IP lookup
// Method: HTTP method to use (GET/POST/etc)
// IP4Rule: Regex pattern for extracting IPv4 from response (empty if not applicable)
// IP6Rule: Regex pattern for extracting IPv6 from response ("=" means entire response body)
// ip4regexp: Compiled regex for IPv4 parsing
// ip6regexp: Compiled regex for IPv6 parsing
// Disabled: Whether this provider is temporarily disabled
type WANIPProvider struct {
	Name        string
	Description string
	URL         string
	Method      string
	IP4Rule     string // 如果是IPv4规则，此项可以为空
	IP6Rule     string // 如果是IPv6规则，此项不能为空（“=”代表整个body数据）
	ip4regexp   *regexp.Regexp
	ip6regexp   *regexp.Regexp
	Disabled    bool
}

// HasIPv6Rule reports whether the WANIPProvider has any IPv6 rules configured.
func (w *WANIPProvider) HasIPv6Rule() bool {
	return len(w.IP6Rule) > 0
}

// HasIPv4Rule returns true if there are IPv4 rules defined or if there are no IPv6 rules.
func (w *WANIPProvider) HasIPv4Rule() bool {
	return len(w.IP4Rule) > 0 || !w.HasIPv6Rule()
}

type WANIPProviders map[string]*WANIPProvider

// Reload reloads all registered WAN IP providers, registering valid ones and unregistering invalid ones.
// Returns nil if the receiver is nil or after completing the reload operation.
func (w *WANIPProviders) Reload() error {
	if w == nil {
		return nil
	}
	for key, value := range *w {
		if value != nil && len(value.Name) > 0 && len(value.URL) > 0 {
			Register(value)
		} else {
			Unregister(key)
		}
	}
	return nil
}

var (
	wanIPProviders   = map[string]*WANIPProvider{}
	defaultProviders = []*WANIPProvider{
		// IPv4
		{
			Name:        `oray.com`,
			Description: `oray`,
			URL:         `https://ddns.oray.com/checkip`,
			IP4Rule:     IPv4Rule,
		}, {
			Name:        `ip-api.com`,
			Description: ``,
			URL:         `http://ip-api.com/json/?fields=query`,
			IP4Rule:     `"query":"` + IPv4Rule + `"`,
		}, {
			Name:        `myip.la`,
			Description: ``,
			URL:         `https://api.myip.la/`,
			IP4Rule:     ``,
		}, {
			Name:        `ipify.org`,
			Description: ``,
			URL:         `https://api.ipify.org`,
			IP4Rule:     ``,
		}, {
			Name:        `3322.org`,
			Description: ``,
			URL:         `http://members.3322.org/dyndns/getip`,
			IP4Rule:     ``,
		}, {
			Name:        `ip.sb`,
			Description: ``,
			URL:         `https://api.ip.sb/ip`,
			IP4Rule:     ``,
		}, {
			Name:        `ipconfig.io`,
			Description: ``,
			URL:         `https://ipconfig.io/ip`,
			IP4Rule:     ``,
		},
		// IPv6
		{
			Name:        `ident.me`,
			Description: ``,
			URL:         `https://v6.ident.me`,
			IP6Rule:     `=`,
		}, {
			Name:        `api-ipv6.ip.sb`,
			Description: ``,
			URL:         `https://api-ipv6.ip.sb/ip`,
			IP6Rule:     `=`,
		}, {
			Name:        `v6.myip.la`,
			Description: ``,
			URL:         `https://v6.myip.la/`,
			IP6Rule:     `=`,
		},
	}
)

func init() {
	extend.Register(`wanIPProvider`, func() interface{} {
		return &WANIPProviders{} // 更新时会自动调用 WANIPProviders.Reload()
	})
	for _, provider := range defaultProviders {
		if err := Register(provider); err != nil {
			panic(err)
		}
	}
}

// Register adds a WANIPProvider to the available providers map after validating and compiling its IP rules.
// It validates and compiles the IPv4 and IPv6 regex patterns if provided, and sets a default HTTP method if none specified.
// Returns an error if any regex compilation fails.
func Register(p *WANIPProvider) (err error) {
	if len(p.IP4Rule) > 0 && p.IP4Rule != `=` {
		p.ip4regexp, err = regexp.Compile(p.IP4Rule)
		if err != nil {
			return
		}
	}
	if len(p.IP6Rule) > 0 && p.IP6Rule != `=` {
		p.ip6regexp, err = regexp.Compile(p.IP6Rule)
		if err != nil {
			return
		}
	}
	if len(p.Method) == 0 {
		p.Method = `GET`
	}
	wanIPProviders[p.Name] = p
	return
}

// Get returns the WANIPProvider instance registered with the given name.
// Returns nil if no provider is found for the specified name.
func Get(name string) *WANIPProvider {
	p, _ := wanIPProviders[name]
	return p
}

// Unregister removes the specified WAN IP providers from the registry.
// It takes one or more provider names as arguments and deletes them from the active providers list.
func Unregister(names ...string) {
	for _, name := range names {
		delete(wanIPProviders, name)
	}
}
