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

// NewFactory creates and returns a new Factory instance with generic type parameters UK and GK.
// UK and GK must be numeric types implementing the com.Number interface.
func NewFactory[UK com.Number, GK com.Number]() *Factory[UK, GK] {
	return &Factory[UK, GK]{
		userFilters:  com.NewSafeMap[UK, *IPFilter](),
		groupFilters: com.NewSafeMap[GK, *IPFilter](),
	}
}

// Factory is a generic struct that manages IP filters for users and groups.
// It maintains separate filter maps for users (keyed by user ID) and groups (keyed by group ID),
// and provides a mechanism to retrieve group-specific IP filter lists.
// The struct uses a singleflight.Group to prevent duplicate work for concurrent requests.
type Factory[UK com.Number, GK com.Number] struct {
	userFilters  *com.SafeMap[UK, *IPFilter]
	groupFilters *com.SafeMap[GK, *IPFilter]
	groupGetter  func(ctx echo.Context, groupID GK) (ipBlacklist string, ipWhitelist string, err error)
	sg           singleflight.Group
}

// SetGroupGetter sets the function to retrieve IP blacklist and whitelist strings for a given group ID.
// The provided function should return the IP lists or an error if the lookup fails.
func (f *Factory[UK, GK]) SetGroupGetter(fn func(ctx echo.Context, groupID GK) (ipBlacklist string, ipWhitelist string, err error)) *Factory[UK, GK] {
	f.groupGetter = fn
	return f
}

// DeleteUser deletes the user with the specified userID and returns the Factory for chaining.
func (f *Factory[UK, GK]) DeleteUser(userID UK) *Factory[UK, GK] {
	f.userFilters.Delete(userID)
	return f
}

// DeleteGroup deletes the group with the specified groupID and returns the Factory for chaining.
func (f *Factory[UK, GK]) DeleteGroup(groupID GK) *Factory[UK, GK] {
	f.groupFilters.Delete(groupID)
	return f
}

// IsAllowed checks if the given IP address is allowed based on user-specific and group-specific IP filters.
// It first checks the user's IP filter (creating one if none exists), then falls back to group-level filtering if applicable.
// Returns true if the IP is allowed by either filter, false otherwise.
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

// isAllowedGroup checks if the given IP address is allowed for the specified group.
// It first checks the cached filters for the group, and if not found, retrieves them using the groupGetter.
// Returns true if the IP is allowed by the group's filters or if no filters exist for the group.
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

// NewWithIP creates a new IPFilter instance and initializes it with the provided
// blacklist and whitelist IP addresses. Each list should be a newline-separated
// string of IP addresses or ranges. Invalid entries are silently logged and skipped.
// Returns the initialized IPFilter.
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
