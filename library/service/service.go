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
	"fmt"
	stdLog "log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/admpub/log"
	"github.com/admpub/service"
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"
)

func ValidServiceAction(action string) error {
	for _, act := range service.ControlAction {
		if act == action {
			return nil
		}
	}
	return fmt.Errorf("available actions: %q", service.ControlAction)
}

func slog() *log.Logger {
	return log.GetLogger(`service`)
}

// New 以服务的方式启动nging
// 服务支持的操作有：
// nging service install  	-- 安装服务
// nging service uninstall  -- 卸载服务
// nging service start 		-- 启动服务
// nging service stop 		-- 停止服务
// nging service restart 	-- 重启服务
func New(cfg *Config, action string) error {
	p := NewProgram(cfg)
	s, err := service.New(p, &p.Config.Config)
	if err != nil {
		return err
	}
	p.service = s

	// Service
	if action != `run` {
		if err := ValidServiceAction(action); err != nil {
			return err
		}
		err = service.Control(s, action)
		if err != nil {
			slog().Errorf(`%s: %s`, action, err.Error())
		} else {
			slog().Okayf(`%s: success`, action)
		}
		return err
	}
	return s.Run()
}

func getPidFiles() []string {
	pidFile := []string{}
	pidFilePath := filepath.Join(com.SelfDir(), `data/pid`)
	err := filepath.Walk(pidFilePath, func(pidPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if info.Name() == `daemon` { // 忽略进程值守创建的进程ID，避免被清理
				err = filepath.SkipDir
			}
			return err
		}
		if filepath.Ext(pidPath) == `.pid` {
			pidFile = append(pidFile, pidPath)
		}
		return nil
	})
	if err != nil {
		stdLog.Println(err)
	}
	return pidFile
}

func NewProgram(cfg *Config) *program {
	pidFile := filepath.Join(com.SelfDir(), `data/pid`)
	err := com.MkdirAll(pidFile, os.ModePerm)
	if err != nil {
		stdLog.Println(err)
	}
	pidFile = filepath.Join(pidFile, `nging.pid`)
	return &program{
		Config:  cfg,
		pidFile: pidFile,
	}
}

type program struct {
	*Config
	service  service.Service
	cmd      *exec.Cmd
	fullExec string
	pidFile  string
}

func (p *program) Start(s service.Service) (err error) {
	if service.Interactive() {
		p.logger.Info("Running in terminal.")
	} else {
		p.logger.Info("Running under service manager.")
	}
	if filepath.Base(p.Exec) == p.Exec {
		p.fullExec, err = exec.LookPath(p.Exec)
		if err != nil {
			return fmt.Errorf("failed to find executable %q: %v", p.Exec, err)
		}
	} else {
		p.fullExec = p.Exec
	}

	p.createCmd()
	// Start should not block. Do the actual work async.
	go p.retryableRun()
	return nil
}

func (p *program) createCmd() {
	p.cmd = exec.Command(p.fullExec, p.Args...)
	p.cmd.Dir = p.Dir
	p.cmd.Env = param.StringSlice(append(os.Environ(), p.Env...)).Unique().String()
	if p.Stderr != nil {
		p.cmd.Stderr = p.Stderr
	}
	if p.Stdout != nil {
		p.cmd.Stdout = p.Stdout
	}
	p.logger.Infof("Running cmd: %s %#v", p.fullExec, p.Args)
	p.logger.Infof("Workdir: %s", p.cmd.Dir)
	//p.logger.Infof("Env var: %s", com.Dump(p.cmd.Env, false))
}

func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	p.killCmd()
	p.logger.Infof("Stopping %s", p.DisplayName)
	if service.Interactive() {
		os.Exit(0)
	}
	return nil
}

func (p *program) killCmd() {
	err := com.CloseProcessFromCmd(p.cmd)
	if err != nil {
		p.logger.Error(err)
	}
	err = com.CloseProcessFromPidFile(p.pidFile)
	if err != nil {
		p.logger.Error(p.pidFile+`:`, err)
	}
	for _, pidFile := range getPidFiles() {
		err = com.CloseProcessFromPidFile(pidFile)
		if err != nil {
			p.logger.Error(pidFile+`:`, err)
		}
	}
}

func (p *program) close() {
	if service.Interactive() {
		p.Stop(p.service)
	} else {
		p.service.Stop()
		p.killCmd()
	}
	if p.Config.OnExited != nil {
		err := p.Config.OnExited()
		if err != nil {
			p.logger.Error(err)
		}
	}
}

func (p *program) run(logPrefix string) error {
	p.logger.Infof(logPrefix+"Starting %s", p.DisplayName)
	err := p.cmd.Start()
	if err == nil {
		stdLog.Println(logPrefix+"APP PID:", p.cmd.Process.Pid)
		os.WriteFile(p.pidFile, []byte(strconv.Itoa(p.cmd.Process.Pid)), os.ModePerm)
		err = p.cmd.Wait()
	}
	if err != nil {
		p.logger.Error(logPrefix+"Error running: ", err)
	}
	return err
}

func (p *program) isSelfUpgradeRestart() bool {
	return p.cmd != nil &&
		p.cmd.ProcessState != nil &&
		p.cmd.ProcessState.ExitCode() == ExitCodeSelfRestart
}

func (p *program) retryableRun() {
	var err error
	defer func() {
		if p.cmd != nil && p.cmd.ProcessState != nil {
			var result string
			switch p.cmd.ProcessState.ExitCode() {
			case 0:
				result = `successful`
			case ExitCodeSelfRestart:
				result = `self-restart`
			default:
				result = `failed`
			}
			p.logger.Infof("Process execution result: %s", result)
			if err != nil {
				p.killCmd()
				log.Close()
				os.Exit(p.cmd.ProcessState.ExitCode())
				return
			}
		}
		p.close()
	}()
	// Do work here
	err = p.run(``)
	if err == nil || p.isSelfUpgradeRestart() {
		return
	}
	maxRetries := p.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}
	retryInterval := p.RetryInterval
	if retryInterval <= 0 {
		retryInterval = DefaultRetryInterval
	}
	wait := time.Second * time.Duration(retryInterval)
	for i := 1; i < maxRetries; i++ {
		progress := fmt.Sprintf(`[%d/%d]`, i, maxRetries)
		p.logger.Infof(progress+"[retry] start %s after %v", p.DisplayName, wait)
		time.Sleep(wait)
		p.killCmd()
		p.createCmd()
		err = p.run(progress + `[retry]`)
		if err == nil || p.isSelfUpgradeRestart() {
			return
		}
	}
}
