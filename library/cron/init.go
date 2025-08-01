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

package cron

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync/atomic"
	"time"

	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
)

var historyJobsRunning atomic.Bool

func InitJobs(ctx context.Context) error {
	m := new(dbschema.NgingTask)
	limit := 1000
	cnt, err := m.ListByOffset(nil, nil, 0, limit, "disabled", `N`)
	if err != nil {
		return fmt.Errorf(`failed to query nging_task list: %w`, err)
	}
	total := int(cnt())
	for offset := 0; offset < total; offset += limit {
		if offset > 0 {
			_, err := m.ListByOffset(nil, nil, offset, limit, "disabled", `N`)
			if err != nil {
				return err
			}
		}
		for _, task := range m.Objects() {
			if err := SaveScriptFile(task); err != nil {
				log.Error(err.Error())
			}
			job, err := NewJobFromTask(ctx, task)
			if err != nil {
				log.Errorf("failed to cron.InitJobs(%d): %v", task.Id, err.Error())
				continue
			}
			if AddJob(task.CronSpec, job) {
				log.Infof("cron.InitJobs: 添加任务[%d]", task.Id)
				continue
			}
		}
	}
	historyJobsRunning.Store(true)
	return nil
}

func HistoryJobsRunning() bool {
	return historyJobsRunning.Load()
}

func runCmdWithTimeout(cmd *exec.Cmd, timeout time.Duration, ctx context.Context) (error, bool) {
	if ctx == nil {
		ctx = context.Background()
	}
	done := make(chan error)
	go func() {
		done <- cmd.Wait()
		close(done)
	}()
	var err error
	kill := func() {
		go func() {
			<-done // 读出上面的goroutine数据，避免阻塞导致无法退出
		}()
		if err = cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
			log.Errorf("进程[%d]无法关闭, 错误信息: %s", cmd.Process.Pid, err)
		}
	}
	t := time.NewTimer(timeout)
	defer t.Stop()
	select {
	case <-t.C:
		log.Warnf("任务执行时间超过%d秒，强制关闭进程: %d", int(timeout/time.Second), cmd.Process.Pid)
		kill()
		return err, true
	case <-ctx.Done():
		kill()
		return err, false
	case err = <-done:
		return err, false
	}
}
