package ipfilter

import (
	"net/netip"
	"strings"

	"github.com/admpub/log"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/param"
	"golang.org/x/sync/singleflight"
)

func NewFactory[UK com.Number, GK com.Number]() *Factory[UK, GK] {
	return &Factory[UK, GK]{
		userFilters:  com.NewSafeMap[UK, *IPFilter](),
		groupFilters: com.NewSafeMap[GK, *IPFilter](),
	}
}

type Factory[UK com.Number, GK com.Number] struct {
	userFilters  *com.SafeMap[UK, *IPFilter]
	groupFilters *com.SafeMap[GK, *IPFilter]
	groupGetter  func(ctx echo.Context, groupID GK) (ipBlacklist string, ipWhitelist string, err error)
	sg           singleflight.Group
}

func (f *Factory[UK, GK]) SetGroupGetter(fn func(ctx echo.Context, groupID GK) (ipBlacklist string, ipWhitelist string, err error)) *Factory[UK, GK] {
	f.groupGetter = fn
	return f
}

func (f *Factory[UK, GK]) DeleteUser(userID UK) *Factory[UK, GK] {
	f.userFilters.Delete(userID)
	return f
}

func (f *Factory[UK, GK]) DeleteGroup(groupID GK) *Factory[UK, GK] {
	f.groupFilters.Delete(groupID)
	return f
}

func (f *Factory[UK, GK]) IsAllowed(ctx echo.Context, userID UK, groupID GK, ipBlacklist string, ipWhitelist string, ip netip.Addr) bool {
	filter, ok := f.userFilters.GetOk(userID)
	if !ok {
		v, _, _ := f.sg.Do(`user:`+param.AsString(userID), func() (interface{}, error) {
			filter := NewWithIP(ipBlacklist, ipWhitelist)
			f.userFilters.Set(userID, filter)
			return filter, nil
		})
		filter = v.(*IPFilter)
	}
	if !filter.IsAllowedAddr(ip) {
		return false
	}
	if groupID > 0 {
		return f.isAllowedGroup(ctx, groupID, ip)
	}
	return true
}

func (f *Factory[UK, GK]) isAllowedGroup(ctx echo.Context, groupID GK, ip netip.Addr) bool {
	filter, ok := f.groupFilters.GetOk(groupID)
	if !ok && f.groupGetter != nil {
		v, _, _ := f.sg.Do(`group:`+param.AsString(groupID), func() (interface{}, error) {
			ipBlacklist, ipWhitelist, err := f.groupGetter(ctx, groupID)
			if err != nil {
				return nil, err
			}
			filter := NewWithIP(ipBlacklist, ipWhitelist)
			f.groupFilters.Set(groupID, filter)
			return filter, nil
		})
		filter, _ = v.(*IPFilter)
	}
	if filter == nil {
		return true
	}
	return filter.IsAllowedAddr(ip)
}

func NewWithIP(ipBlacklist, ipWhitelist string) *IPFilter {
	filter := New()
	ipBlacklist = strings.TrimSpace(ipBlacklist)
	if len(ipBlacklist) > 0 {
		ips := param.Split(ipBlacklist, "\n").Filter(isNotEmptyString).Unique().String()
		err := filter.AddBlacklist(ips...)
		if err != nil {
			log.Errorf(`failed to add ipBlacklist: %v`, err)
		}
	}
	ipWhitelist = strings.TrimSpace(ipWhitelist)
	if len(ipWhitelist) > 0 {
		ips := param.Split(ipWhitelist, "\n").Filter(isNotEmptyString).Unique().String()
		err := filter.AddWhitelist(ips...)
		if err != nil {
			log.Errorf(`failed to add ipWhitelist: %v`, err)
		}
	}
	return filter
}
