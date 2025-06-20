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
	"runtime/debug"
	"strings"
	"sync/atomic"
	"time"

	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/cron/send"
	cronWriter "github.com/coscms/webcore/library/cron/writer"
	"github.com/coscms/webcore/library/httpserver"
	"github.com/coscms/webcore/registry/alert"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
	"github.com/webx-top/echo/engine"
	"github.com/webx-top/echo/middleware/tplfunc"
	"github.com/webx-top/echo/param"
	"golang.org/x/sync/singleflight"

	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/charset"
)

var (
	defaultOuputSize uint64 = 2000
	cmdPreParams     []string

	// ErrFailure 报错:执行失败
	ErrFailure = errors.New(`Error`)
)

// Runner 命令运行
type Runner func(timeout time.Duration) (out string, runingErr string, onRunErr error, isTimeout bool)

type RunnerGetter func(string) Runner

func init() {
	if com.IsWindows {
		cmdPreParams = []string{"cmd.exe", "/c"}
		//cmdPreParams = []string{"bash.exe", "-c"}
	} else {
		shell := os.Getenv("SHELL")
		if len(shell) == 0 {
			shell = "/bin/bash"
			if _, err := os.Stat(shell); err != nil {
				shell = "/bin/sh"
			}
		}
		cmdPreParams = []string{shell, "-c"}
	}
}

var cmdReplacer = strings.NewReplacer("\r", ``)

func CmdParams(command string) []string {
	params := append([]string{}, cmdPreParams...)
	if !com.IsWindows {
		command = cmdReplacer.Replace(command)
	}
	params = append(params, command)
	return params
}

// Job 定义需要处理的job
type Job struct {
	id         uint                // 任务ID
	logID      atomic.Uint64       // 日志记录ID
	name       string              // 任务名称
	task       *dbschema.NgingTask // 任务对象
	runner     Runner              // 执行函数
	system     bool                // 是否是系统内部功能
	status     atomic.Int32        // 任务状态，大于0表示正在执行中
	concurrent bool                // 同一个任务是否允许并行执行
}

func NewJobFromTask(ctx context.Context, task *dbschema.NgingTask) (*Job, error) {
	if task.Id < 1 {
		return nil, echo.NewError("Job: missing task.Id", code.DataNotFound)
	}
	var env []string
	task.Env = strings.TrimSpace(task.Env)
	if len(task.Env) > 0 {
		for _, row := range strings.Split(task.Env, "\n") {
			row = strings.TrimSpace(row)
			if len(row) > 0 {
				env = append(env, row)
			}
		}
	}
	cmd := task.Command
	if len(cmd) > 1 && cmd[0] == '>' {
		cmd = cmd[1:]
		cmdInfo := strings.SplitN(cmd, `:`, 2)
		var param string
		switch len(cmdInfo) {
		case 2:
			param = cmdInfo[1]
			fallthrough
		case 1:
			fnName := cmdInfo[0]
			jobx, ok := systemJobs[fnName]
			if !ok {
				return nil, echo.NewError(fmt.Sprintf("Job: invalid job name: %s", fnName), code.InvalidParameter)
			}
			job := &Job{
				id:         task.Id,
				name:       task.Name,
				task:       task,
				concurrent: task.Concurrent == 1,
				runner:     jobx.RunnerGetter(param),
				system:     true,
			}
			return job, nil
		}
	}
	if task.GroupId > 0 {
		group := dbschema.NewNgingTaskGroup(task.Context())
		err := group.Get(nil, `id`, task.GroupId)
		if err != nil {
			return nil, err
		}
		if len(group.CmdPrefix) > 0 {
			cmd = group.CmdPrefix + ` ` + cmd
		}
		if len(group.CmdSuffix) > 0 {
			cmd += ` ` + group.CmdSuffix
		}
	}
	job := NewCommandJob(ctx, task.Id, task.Name, cmd, task.WorkDirectory, env...)
	job.task = task
	job.concurrent = task.Concurrent == 1
	return job, nil
}

func NewOutputWriter(sizes ...uint64) OutputWriter {
	var size uint64
	if len(sizes) > 0 {
		size = sizes[0]
	}
	if size == 0 {
		size = defaultOuputSize
	}
	return NewCmdRec(size)
}

func NewCommandJob(ctx context.Context, id uint, name string, command string, dir string, env ...string) *Job {
	job := &Job{
		id:   id,
		name: name,
	}
	command = ScriptCommand(id, command)
	job.runner = func(timeout time.Duration) (string, string, error, bool) {
		if ctx == nil {
			ctx = context.Background()
		}
		bufOut := NewCmdRec(defaultOuputSize)
		bufErr := NewCmdRec(defaultOuputSize)
		params := CmdParams(command)
		cmd := exec.Command(params[0], params[1:]...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(), env...)
		cmd.Stdout = bufOut
		cmd.Stderr = bufErr
		if err := cmd.Start(); err != nil {
			return ``, ``, err, false
		}
		err, isTimeout := runCmdWithTimeout(cmd, timeout, ctx)
		if com.IsWindows {
			bOut, e := charset.Convert(`gbk`, `utf-8`, bufOut.Bytes())
			if e != nil {
				log.Errorf(`failed to charset.Convert(bufOut, gbk, utf-8): %v`, e)
			}
			bErr, e := charset.Convert(`gbk`, `utf-8`, bufErr.Bytes())
			if e != nil {
				log.Errorf(`failed to charset.Convert(bufErr, gbk, utf-8): %v`, e)
			}
			return engine.Bytes2str(bOut), engine.Bytes2str(bErr), err, isTimeout
		}

		return bufOut.String(), bufErr.String(), err, isTimeout
	}
	return job
}

func (j *Job) Status() int32 {
	return j.status.Load()
}

func (j *Job) Name() string {
	return j.name
}

func (j *Job) Id() uint {
	return j.id
}

func (j *Job) LogID() uint64 {
	return j.logID.Load()
}

func (j *Job) addAndReturningLog(taskLog *dbschema.NgingTaskLog) *Job {
	// 插入日志
	_, err := taskLog.Insert()
	if err != nil {
		log.Error("Job: 日志写入失败: ", err)
	}
	j.logID.Store(taskLog.Id)
	return j
}

func (j *Job) send(elapsed int64, t time.Time, err error, cmdOut string, isTimeout bool, timeout time.Duration) error {
	data := param.Store{
		"task":      *j.task,
		"startTime": t.Format("2006-01-02 15:04:05"),
		"elapsed":   tplfunc.NumberTrim(float64(elapsed)/1000, 6),
		"output":    cmdOut,
	}
	var title, status, statusText string
	if isTimeout {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "超时")
		status = `timeout`
		statusText = fmt.Sprintf("超时（%d秒）", int(timeout/time.Second))
	} else if err != nil {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "失败")
		status = `failure`
		statusText = "失败（" + err.Error() + "）"
	} else {
		title = fmt.Sprintf("任务执行结果通知 #%d: %s", j.task.Id, "成功")
		status = `success`
		statusText = "成功"
	}
	data["title"] = title
	data["status"] = status
	data["statusText"] = statusText
	data["content"] = send.NewContent()
	backendURL := common.Setting(`base`).String(`backendURL`)
	backendURL = strings.TrimSuffix(backendURL, `/`)
	data["detailURL"] = backendURL + httpserver.Backend.Router.Prefix() + `/task/log_view/` + fmt.Sprint(j.LogID())
	return Send(&alert.AlertData{
		Title:   title,
		Content: send.NewContent(),
		Data:    data,
	})
}

var sg singleflight.Group

// Run 运行Job
func (j *Job) Run() {
	if j.concurrent { // 允许同一个任务同时执行多次
		j.run()
		return
	}

	// 不允许同一个任务同时执行多次
	t := time.Now()
	key := com.String(j.id)
	_, _, shared := sg.Do(key, func() (interface{}, error) {
		j.run()
		return nil, nil
	})
	if shared {
		taskLog := new(dbschema.NgingTaskLog)
		taskLog.TaskId = j.id
		taskLog.Created = uint(t.Unix())
		taskLog.Output = fmt.Sprintf("任务[ %d. %s ]上一次执行尚未结束，本次被忽略。", j.id, j.name)
		j.addAndReturningLog(taskLog)
	}
}

func (j *Job) run() {
	var (
		cmdOut    string
		cmdErr    string
		err       error
		isTimeout bool
	)
	t := time.Now()
	taskLog := new(dbschema.NgingTaskLog)
	taskLog.TaskId = j.id
	taskLog.Created = uint(t.Unix())
	defer func() {
		taskLog.Output = cmdOut
		taskLog.Error = cmdErr
		if e := recover(); e != nil {
			errMsg := fmt.Sprintf(`[NGING.PANIC] %v`, e)
			if len(taskLog.Error) > 0 {
				taskLog.Error += "\n" + errMsg
			} else {
				taskLog.Error = errMsg
			}
			log.Error(e, "\n", string(debug.Stack()))
			taskLog.Status = `failure`
		}
		if j == nil { // 异常情况
			_, err = taskLog.Insert()
			if err != nil {
				log.Error("Job: 日志写入失败: ", err)
			}
			return
		}
		if j.task.ClosedLog == `N` && !strings.HasPrefix(cmdOut, cronWriter.NotRecordPrefixFlag) && !strings.HasPrefix(cmdErr, cronWriter.NotRecordPrefixFlag) {
			j.addAndReturningLog(taskLog)
		}
	}()

	if !j.concurrent {
		if j.status.Load() > 0 {
			taskLog.Output = fmt.Sprintf("任务[ %d. %s ]上一次执行尚未结束，本次被忽略。", j.id, j.name)
			return
		}
	}

	if workPool != nil {
		workPool <- true
		defer func() {
			if workPool == nil {
				return
			}
			<-workPool
		}()
	}

	log.Debugf("开始执行任务: %d", j.id)

	j.status.Add(1)
	defer j.status.Add(-1)

	timeout := time.Duration(time.Hour * 24)
	if j.task.Timeout > 0 {
		timeout = time.Second * time.Duration(j.task.Timeout)
	}

	cmdOut, cmdErr, err, isTimeout = j.runner(timeout)
	elapsed := time.Since(t).Milliseconds()
	taskLog.Elapsed = uint(elapsed)
	if isTimeout {
		taskLog.Status = `timeout`
		taskLog.Error = fmt.Sprintf("任务执行超过 %d 秒\n----------------------\n", int64(timeout/time.Second))
	} else if err != nil {
		taskLog.Status = `failure`
		taskLog.Error = err.Error()
	} else {
		taskLog.Status = `success`
	}

	// 更新上次执行时间
	j.task.PrevTime = uint(t.Unix())
	j.task.ExecuteTimes++
	setErr := j.task.UpdateFields(nil, map[string]interface{}{
		`prev_time`:     j.task.PrevTime,
		`execute_times`: j.task.ExecuteTimes,
	}, `id`, j.task.Id)
	if setErr != nil {
		log.Errorf(`failed to cron.task.UpdateFields: %v`, setErr)
	}

	// 发送邮件通知
	switch j.task.EnableNotify {
	case NotifyIfFail:
		if err == nil {
			return
		}
		fallthrough
	case NotifyIfEnd:
		out := cmdErr
		if len(out) == 0 {
			out = cmdOut
		}
		err := j.send(elapsed, t, err, out, isTimeout, timeout)
		if err != nil {
			log.Errorf(`failed to send cron job notify: %v`, err)
		}
	case NotifyDisabled:
		return
	}
}

const (
	//NotifyDisabled 不通知
	NotifyDisabled = iota
	//NotifyIfEnd 执行结束时通知
	NotifyIfEnd
	//NotifyIfFail 执行失败时通知
	NotifyIfFail
)
