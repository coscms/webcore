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

package manager

import (
	"github.com/coscms/webcore/cmd/bootconfig"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/library/navigate"
	navRegistry "github.com/coscms/webcore/registry/navigate"
)

var Project = navigate.NewProject(`Nging`, `nging`, `/index`, navRegistry.LeftNavigate)

func init() {
	bootconfig.SupportManager = true
	httpserver.Backend.Navigate.AddProject(2, Project)
}
