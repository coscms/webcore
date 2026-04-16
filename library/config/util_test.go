package config

import (
	"testing"

	"github.com/admpub/confl"
	"github.com/stretchr/testify/assert"
	"github.com/webx-top/echo/middleware/language"
	"github.com/webx-top/echo/param"
)

func TestParseLanguage(t *testing.T) {
	cfg := NewConfig()
	lng := language.Config{}
	lng.AllList = []string{`zh-CN`, `en`}
	lng.Default = `zh-CN`
	lng.Extra = map[string]param.Store{
		"zh-CN": {
			"label": "简体中文",
			"flag":  "🇨🇳",
		},
		"en": {
			"label": "English",
			"flag":  "🇺🇸",
		},
	}
	lng.Reload = true
	_, err := confl.Decode(`
language {
  Default      : "zh-Cn"
  Fallback     : ""
  AllList      : ["zh-CN","en"]
  Reload       : true
  Extra : {
    "zh-CN" : {
      label : "简体中文"
      flag : "🇨🇳"
    }
    "en" : {
      label : "English"
      flag : "🇺🇸"
    }
  }
}
`, cfg)
	assert.NoError(t, err)
	cfg.SetDefaults()
	assert.Equal(t, lng, cfg.Language)
}
