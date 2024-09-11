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

package webcore

import (
	"github.com/coscms/webcore/cmd"
	_ "github.com/coscms/webcore/initialize"
	_ "github.com/coscms/webcore/library/config/sessionstore"
	_ "github.com/coscms/webcore/library/cron/cmder"
	"github.com/coscms/webcore/library/module"
	_ "github.com/coscms/webcore/library/sqlite"
	_ "github.com/coscms/webcore/library/upload"
	_ "github.com/coscms/webcore/listener"
	_ "github.com/coscms/webcore/upgrade"
)

func Start(modules ...module.IModule) {
	module.Register(modules...)
	cmd.Execute()
}

func ReadyStart(ready func(), modules ...module.IModule) {
	module.Register(modules...)
	if ready != nil {
		ready()
	}
	cmd.Execute()
}
