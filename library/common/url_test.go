package common

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/webx-top/com"
)

func TestReverseURL(t *testing.T) {
	r := ReverseURL(`https%3A%2F%2Fwww.coscms.com%2Fuser%2Fprofile%2Fbinding%3Ftype%3Doauth`)
	assert.Equal(t, `https://www.coscms.com/user/profile/binding?type=oauth`, r)
	r = ReverseURL(`http%3A%2F%2Fwww.coscms.com%2Fuser%2Fprofile%2Fbinding%3Ftype%3Doauth`)
	assert.Equal(t, `http://www.coscms.com/user/profile/binding?type=oauth`, r)
	r = ReverseURL(`%2F`)
	assert.Equal(t, `/`, r)
	v := url.QueryEscape(`./`)
	assert.Equal(t, `.%2F`, v)
	r = ReverseURL(v)
	assert.Equal(t, `./`, r)
	r = ReverseURL(`..%2F`)
	assert.Equal(t, `../`, r)
}

func TestSortedURLValues(t *testing.T) {
	r := NewSortedURLValues(`a=b&b=100&a=c`)
	assert.Equal(t, `a`, r[0].Key)
	assert.Equal(t, []string{`b`, `c`}, r[0].Values)
	assert.Equal(t, `b`, r[1].Key)
	assert.Equal(t, []string{`100`}, r[1].Values)
	r.Del(`a`)
	assert.Equal(t, 1, len(r))
	r.Del(`b`)
	assert.Equal(t, 0, len(r))
	r.ParseQuery(`aa=1&ab=2&ac=3&ad=4`)
	com.Dump(r)
	assert.Equal(t, 4, len(r))
	assert.Equal(t, `aa`, r[0].Key)
	assert.Equal(t, []string{`1`}, r[0].Values)
	assert.Equal(t, `ab`, r[1].Key)
	assert.Equal(t, []string{`2`}, r[1].Values)
	assert.Equal(t, `ac`, r[2].Key)
	assert.Equal(t, []string{`3`}, r[2].Values)
	assert.Equal(t, `ad`, r[3].Key)
	assert.Equal(t, []string{`4`}, r[3].Values)
}
