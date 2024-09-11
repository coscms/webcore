package backend

import (
	"testing"

	"github.com/webx-top/echo/testing/test"
)

func TestMakeSubdomains(t *testing.T) {
	r := MakeSubdomains(`www.coscms.com,www.coscms.com:8181`, []string{})
	test.Eq(t, []string{"www.coscms.com", "www.coscms.com:8181", "www.coscms.com:9999"}, r)
	r = MakeSubdomains(`www.coscms.com,www.coscms.com:8181`, DefaultLocalHostNames)
	test.Eq(t, []string{"www.coscms.com", "www.coscms.com:8181", "www.coscms.com:9999", "127.0.0.1:9999", "127.0.0.1:8181", "localhost:9999", "localhost:8181"}, r)
}
