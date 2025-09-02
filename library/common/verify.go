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

package common

import "github.com/webx-top/com"

const (
	BoolY = `Y`
	BoolN = `N`
)

const (
	ContentTypeHTML     = `html`
	ContentTypeMarkdown = `markdown`
	ContentTypeText     = `text`
)

var (
	boolFlags = []string{BoolY, BoolN}
	contypes  = []string{ContentTypeHTML, ContentTypeMarkdown, ContentTypeText}
)

func GetBoolFlag(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, boolFlags) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return BoolN
	}
	return value
}

func IsBoolFlag(value string) bool {
	return com.InSlice(value, boolFlags)
}

func BoolToFlag(v bool) string {
	if v {
		return BoolY
	}
	return BoolN
}

func FlagToBool(v string) bool {
	return v == BoolY
}

func GetContype(value string, defaults ...string) string {
	if len(value) == 0 || !com.InSlice(value, contypes) {
		if len(defaults) > 0 {
			return defaults[0]
		}
		return ContentTypeText
	}
	return value
}

func IsContype(value string) bool {
	return com.InSlice(value, contypes)
}

func GetEnumValue(enums []string, value string, defaults string) string {
	if len(value) == 0 || !com.InSlice(value, enums) {
		return defaults
	}
	return value
}
