package cron

import (
	"fmt"
	"os"
	"strings"

	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/filecache"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

func SaveScriptFile(m *dbschema.NgingTask) error {
	if !com.IsWindows {
		return nil
	}
	name := param.AsString(m.Id) + `.bat`
	if !strings.Contains(m.Command, "\n") {
		_ = filecache.RemoveCache(`taskscripts`, name)
		return nil
	}
	err := filecache.WriteCache(`taskscripts`, name, com.Str2bytes(m.Command))
	if err != nil {
		err = fmt.Errorf(`failed to cron.SaveScriptFile(%q): %w`, name, err)
	}
	return err
}

func DeleteScriptFile(id uint) error {
	if !com.IsWindows {
		return nil
	}
	name := param.AsString(id) + `.bat`
	err := filecache.RemoveCache(`taskscripts`, name)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
	}
	return err
}

func ScriptFile(id uint) string {
	if !com.IsWindows {
		return ``
	}
	name := param.AsString(id) + `.bat`
	return filecache.CacheFile(`taskscripts`, name)
}

func ScriptCommand(id uint, command string) string {
	if !com.IsWindows {
		return command
	}
	if !strings.Contains(command, "\n") {
		return command
	}
	scriptFile := ScriptFile(id)
	if len(scriptFile) == 0 {
		return command
	}
	return scriptFile
}
