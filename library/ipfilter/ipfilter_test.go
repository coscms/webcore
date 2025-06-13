package ipfilter

import (
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/defaults"
)

func TestToPrefixes(t *testing.T) {
	start := `127.0.0.0`
	end := `129.0.0.1`
	pfx, err := ParseIPRange(start, end)
	assert.NoError(t, err)
	assert.Equal(t, []netip.Prefix{
		netip.MustParsePrefix(`127.0.0.0/8`),
		netip.MustParsePrefix(`128.0.0.0/8`),
		netip.MustParsePrefix(`129.0.0.0/31`),
	}, pfx)
	for _, pf := range pfx {
		t.Log(pf.String())
	}
	start = `fe80::`
	end = `febf:ffff:ffff:ffff:ffff:ffff:ffff:ffff`
	pfx, err = ParseIPRange(start, end)
	assert.NoError(t, err)
	assert.Equal(t, []netip.Prefix{netip.MustParsePrefix(`fe80::/10`)}, pfx)
	for _, pf := range pfx {
		t.Log(pf.String())
	}

	start = `127.0.0.0`
	end = `127.0.0.0`
	pfx, err = ParseIPRange(start, end)
	assert.NoError(t, err)
	assert.Equal(t, []netip.Prefix{netip.MustParsePrefix(`127.0.0.0/32`)}, pfx)
	for _, pf := range pfx {
		t.Log(pf.String())
	}

	start = `fe80::`
	end = `fe80::`
	pfx, err = ParseIPRange(start, end)
	assert.NoError(t, err)
	assert.Equal(t, []netip.Prefix{netip.MustParsePrefix(`fe80::/128`)}, pfx)
	for _, pf := range pfx {
		t.Log(pf.String())
	}
}

func TestContains(t *testing.T) {
	i := New()
	err := i.AddWhitelist(`127.0.0.0-129.0.0.1`)
	assert.NoError(t, err)
	ip := `127.0.0.2`
	y := i.IsAllowed(ip)
	assert.True(t, y)
	ip = `128.0.0.2`
	y = i.IsAllowed(ip)
	assert.True(t, y)
	ip = `129.0.0.2`
	y = i.IsAllowed(ip)
	assert.False(t, y)
	ip = `127.0.0.0`
	y = i.IsAllowed(ip)
	assert.True(t, y)

	f := NewWithIP(``, ip).SetDisallow(true)
	y = f.IsAllowed(ip)
	assert.True(t, y)
	f = NewWithIP(ip, ``)
	y = f.IsAllowed(ip)
	assert.False(t, y)
	f = NewWithIP(``, ``)
	y = f.IsAllowed(ip)
	assert.True(t, y)
}

func TestValidate(t *testing.T) {
	err := Validate(`192.168.0.0-192.168.255.255`)
	assert.NoError(t, err)
	err = Validate(`192.168.0.0/16`)
	assert.NoError(t, err)
	err = Validate(`192.168.0.1`)
	assert.NoError(t, err)
	err = Validate(`fe80::`)
	assert.NoError(t, err)

	ctx := defaults.NewMockContext()
	err = ValidateRows(ctx, `192.168.0.0-192.168.254.255 
	
127.0.0.0/16`)
	assert.NoError(t, err)
}
