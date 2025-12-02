package config

import (
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/middleware/language"
)

func (c *Config) initLanguage() *Config {
	if len(c.Language.Default) > 0 {
		c.Language.Default = echo.NewLangCode(c.Language.Default).Normalize()
	}
	if len(c.Language.Fallback) > 0 {
		c.Language.Fallback = echo.NewLangCode(c.Language.Fallback).Normalize()
	}
	if len(c.Language.AllList) == 0 {
		if len(c.Language.Default) == 0 {
			c.Language.Default = `en`
		}
		c.Language.AllList = []string{c.Language.Default}
		return c
	}
	for index, lang := range c.Language.AllList {
		lang = echo.NewLangCode(lang).Normalize()
		if index == 0 && len(c.Language.Default) == 0 {
			c.Language.Default = lang
		}
		c.Language.AllList[index] = lang
	}
	return c
}

func (c *Config) GetTranslator(ctx echo.Context) echo.Translator {
	return c.BuildTranslator(ctx, ctx.Lang().String())
}

func (c *Config) CloneLanguageConfig() language.Config {
	cfg := c.Language.Clone()
	cfg.SetFSFunc(bootconfig.LangFSFunc)
	cfg.Reload = false
	return cfg
}

func (c *Config) NewLanguage() *language.Language {
	cfg := c.CloneLanguageConfig()
	return language.New(&cfg)
}

func (c *Config) BuildTranslator(ctx echo.Context, langCode string) *language.Translate {
	tr := &language.Translate{}
	lng := c.NewLanguage()
	langs := make(map[string]bool, len(c.Language.AllList))
	for _, lang := range c.Language.AllList {
		langs[lang] = true
	}
	tr.Reset(ctx, langCode, lng, langs, c.Language.Default)
	return tr
}
