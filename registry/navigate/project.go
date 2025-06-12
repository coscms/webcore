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

package navigate

import (
	"strings"

	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
)

func ProjectAddNavList(name string, ident string, url string, navList *navigate.List) {
	httpserver.Backend.Navigate.AddNavList(name, ident, url, navList)
}

func ProjectAdd(index int, list ...*navigate.ProjectItem) {
	httpserver.Backend.Navigate.AddProject(index, list...)
}

func ProjectURLsIdent() map[string]string {
	return httpserver.Backend.Navigate.Projects.URLsIdent()
}

func ProjectInit() {
	httpserver.Backend.Navigate.Init()
}

func ProjectIdent(urlPath string) string {
	if len(httpserver.Backend.Prefix()) > 0 {
		urlPath = strings.TrimPrefix(urlPath, httpserver.Backend.Prefix())
	}
	return httpserver.Backend.Navigate.Projects.GetIdent(urlPath)
}

func ProjectRemoveByIdent(ident string) {
	httpserver.Backend.Navigate.Projects.RemoveByIdent(ident)
}

func ProjectRemove(index int) {
	httpserver.Backend.Navigate.Projects.Remove(index)
}

func ProjectSearchIdent(ident string) int {
	return httpserver.Backend.Navigate.Projects.List.SearchIdent(ident)
}

func ProjectSet(index int, list ...*navigate.ProjectItem) {
	httpserver.Backend.Navigate.Projects.Set(index, list...)
}

func ProjectListAll() navigate.ProjectList {
	return *httpserver.Backend.Navigate.Projects.List
}

func ProjectGet(ident string) *navigate.ProjectItem {
	return httpserver.Backend.Navigate.Projects.Get(ident)
}

func ProjectFirst(notEmptyOpts ...bool) *navigate.ProjectItem {
	return httpserver.Backend.Navigate.Projects.First(notEmptyOpts...)
}
