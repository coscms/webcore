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

package service

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPidFile(t *testing.T) {
	pidFile := createPidFile()
	pidDir := filepath.Dir(pidFile)
	err := os.MkdirAll(pidDir+`/daemon/1`, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(pidDir+`/daemon/2`, os.ModePerm)
	assert.NoError(t, err)
	err = os.MkdirAll(pidDir+`/daemon/3`, os.ModePerm)
	assert.NoError(t, err)
	files := getPidFiles()
	assert.Equal(t, []string{}, files)
}
