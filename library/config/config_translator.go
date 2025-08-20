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
	return c.BuildTranslator(ctx.Lang().String())
}

func (c *Config) CloneLanguageConfig() language.Config {
	cfg := language.Config{
		Project:      c.Language.Project,
		Default:      c.Language.Default,
		Fallback:     c.Language.Fallback,
		AllList:      make([]string, len(c.Language.AllList)),
		RulesPath:    make([]string, len(c.Language.RulesPath)),
		MessagesPath: make([]string, len(c.Language.MessagesPath)),
		Reload:       false,
	}
	copy(cfg.AllList, c.Language.AllList)
	copy(cfg.RulesPath, c.Language.RulesPath)
	copy(cfg.MessagesPath, c.Language.MessagesPath)
	cfg.SetFSFunc(bootconfig.LangFSFunc)
	return cfg
}

func (c *Config) NewLanguage() *language.Language {
	cfg := c.CloneLanguageConfig()
	return language.New(&cfg)
}

func (c *Config) BuildTranslator(langCode string) *language.Translate {
	tr := &language.Translate{}
	lng := c.NewLanguage()
	tr.Reset(langCode, lng)
	return tr
}
