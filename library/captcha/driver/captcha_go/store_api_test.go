package captcha_go

import (
	"testing"

	"github.com/webx-top/com"
)

func TestStoreAPI(t *testing.T) {
	body := map[string]interface{}{
		`key`: `234`,
	}
	b, _ := com.JSONEncode(body)
	token := com.Token(`1.0.1`, b)
	t.Log(token)
}
