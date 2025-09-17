/*
   Nging is a toolbox for webmasters
   Copyright (C) 2018-present Wenhui Shen <swh@admpub.com>

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published
   by the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package clitranslator

import (
	"github.com/admpub/once"
	"github.com/coscms/webcore/library/config"
	"github.com/webx-top/echo/middleware/language"
)

var translate *language.Translate
var translock once.Once
var LangCode = `zh-CN` // 默认语言

func initTranslate() {
	translate = config.FromFile().BuildTranslator(LangCode)
}

func GetTranslator() *language.Translate {
	translock.Do(initTranslate)
	return translate
}

func ResetTranslator() {
	translock.Reset()
}
