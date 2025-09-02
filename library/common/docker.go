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

import (
	"bufio"
	"os"
	"strings"
)

func InDocker() bool {
	if len(os.Getenv(`container`)) > 0 {
		return true
	}
	for _, file := range []string{`/.dockerenv`, `/.dockerinit`} {
		_, err := os.Stat(file)
		if err == nil {
			return true
		}
	}

	// cgroup
	cgroupFile := `/proc/1/cgroup`
	if file, err := os.Open(cgroupFile); err == nil {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "/docker/") || strings.Contains(line, "kubepods") || strings.Contains(line, `/lxc/`) || strings.Contains(line, `/kubepod/`) {
				file.Close()
				return true
			}
		}
		file.Close()
	}

	// sched
	schedFile := `/proc/1/sched`
	if b, err := os.ReadFile(schedFile); err == nil {
		content := string(b)
		content = strings.TrimSpace(content)
		exists := strings.HasPrefix(content, `systemd `) || strings.HasPrefix(content, `init `)
		return !exists
	}
	return false
}
