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

package checkinstall

import (
	"os/exec"

	"github.com/admpub/once"
)

func New(name string, checker ...func(name string) bool) *CheckInstall {
	c := &CheckInstall{name: name}
	if len(checker) > 0 && checker[0] != nil {
		c.checker = checker[0]
	} else {
		c.checker = DefaultChecker
	}
	return c
}

var DefaultChecker = func(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

type CheckInstall struct {
	name      string
	installed bool
	checkonce once.Once
	checker   func(name string) bool
}

func (c *CheckInstall) check() {
	c.installed = c.checker(c.name)
}

func (c *CheckInstall) IsInstalled() bool {
	c.checkonce.Do(c.check)
	return c.installed
}

func (c *CheckInstall) ResetCheck() {
	c.checkonce.Reset()
}
