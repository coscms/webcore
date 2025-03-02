package license

import "strings"

func MakeUserAgent() string {
	pkgName := Package()
	if len(pkgName) > 0 {
		pkgName = `-` + pkgName
	}
	name := License().Info.Name
	if len(name) > 0 {
		name = ` ` + name
	}
	return ProductName() + pkgName + `/` + Version() + name
}

type UserAgentRaw struct {
	Product string
	Package string
	Version string
	User    string
}

func ParseUserAgent(userAgent string) UserAgentRaw {
	parts := strings.SplitN(userAgent, `/`, 2)
	prods := strings.SplitN(parts[0], `-`, 2)
	r := UserAgentRaw{}
	r.Product = prods[0]
	if len(prods) > 1 {
		r.Package = prods[1]
	}
	if len(parts) > 1 {
		_parts := strings.SplitN(parts[1], ` `, 2)
		r.Version = _parts[0]
		if len(_parts) > 1 {
			r.User = _parts[1]
		}
	}
	return r
}
