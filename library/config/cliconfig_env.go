package config

import (
	"os"
	"path/filepath"

	"github.com/admpub/godotenv"
	"github.com/admpub/log"
	"github.com/coscms/webcore/library/common"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func findEnvFile() []string {
	var envFiles []string
	envFile := filepath.Join(echo.Wd(), `.env`)
	if fi, err := os.Stat(envFile); err == nil && !fi.IsDir() {
		envFiles = append(envFiles, envFile)
	}
	return envFiles
}

func (c *CLIConfig) InitEnviron(needFindEnvFile ...bool) (err error) {
	c.envLock.Lock()
	defer c.envLock.Unlock()
	if len(needFindEnvFile) > 0 && needFindEnvFile[0] {
		c.envFiles = findEnvFile()
	}
	var newEnvVars map[string]string
	if len(c.envFiles) > 0 {
		log.Infof(`Loading env file: %#v`, c.envFiles)
		newEnvVars, err = godotenv.Read(c.envFiles...)
		if err != nil {
			return
		}
	}
	if len(newEnvVars) == 0 {
		if c.envVars == nil {
			return
		}
		for k := range c.envVars {
			log.Infof(`Unset env var: %s`, k)
			os.Unsetenv(k)
		}
		c.envVars = nil
		return
	}
	if c.envVars != nil {
		for k, v := range c.envVars {
			newV, ok := newEnvVars[k]
			if !ok {
				log.Infof(`Unset env var: %s`, k)
				os.Unsetenv(k)
				delete(c.envVars, k)
				continue
			}
			if v != newV {
				log.Infof(`Set env var: %s`, k)
				os.Setenv(k, v)
				c.envVars[k] = newV
			}
			delete(newEnvVars, k)
		}
	} else {
		c.envVars = map[string]string{}
	}
	for k, v := range newEnvVars {
		log.Infof(`Set env var: %s`, k)
		os.Setenv(k, v)
		c.envVars[k] = v
	}
	return
}

func (c *CLIConfig) WatchEnvConfig() {
	if c.envMonitor != nil {
		c.envMonitor.Close()
		c.envMonitor = nil
	}
	if len(c.envFiles) == 0 {
		return
	}
	c.envMonitor = &com.MonitorEvent{
		Modify: func(file string) {
			log.Info(`Start reloading env file: ` + file)
			err := c.InitEnviron()
			if err == nil {
				log.Info(`Succcessfully reload the env file: ` + file)
				return
			}
			if err == common.ErrIgnoreConfigChange {
				log.Info(`No need to reload the env file: ` + file)
				return
			}
			log.Errorf(`failed to cliconfig.InitEnviron: %v`, err)
		},
	}
	for _, envFile := range c.envFiles {
		err := c.envMonitor.AddFile(envFile)
		if err != nil {
			log.Errorf(`failed to envMonitor.AddFile(%q): %v`, envFile, err)
		}
	}
	c.envMonitor.Watch()
}
