package backend

import (
	"io"
	"os"

	"github.com/admpub/log"
	"github.com/coscms/webcore/dbschema"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/library/sessionguard"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func User(ctx echo.Context) *dbschema.NgingUser {
	user, ok := ctx.Internal().Get(`user`).(*dbschema.NgingUser)
	if ok && user != nil {
		return user
	}
	user, ok = ctx.Session().Get(`user`).(*dbschema.NgingUser)
	if ok {
		if !sessionguard.Validate(ctx, user.LastIp, `user`, uint64(user.Id)) {
			log.Warn(ctx.T(`用户“%s”的会话环境发生改变，需要重新登录`, user.Username))
			ctx.Session().Delete(`user`)
			return nil
		}
		ctx.Internal().Set(`user`, user)
	}
	return user
}

func NoticeWriter(ctx echo.Context, noticeType string) (wOut io.Writer, wErr io.Writer, err error) {
	user := User(ctx)
	if user == nil {
		return nil, nil, ctx.Redirect(URLFor(`/login`))
	}
	typ := `service:` + noticeType
	notice.OpenMessage(user.Username, typ)

	wOut = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stdout.Write(b)
		notice.Send(user.Username, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Succeed))
		return nil
	}}
	wErr = &com.CmdResultCapturer{Do: func(b []byte) error {
		os.Stderr.Write(b)
		notice.Send(user.Username, notice.NewMessageWithValue(typ, noticeType, string(b), notice.Failed))
		return nil
	}}
	return
}

func IsBackendAdmin(c echo.Context) bool {
	if user := User(c); user != nil {
		return user.Id > 0
	}
	return false
}
