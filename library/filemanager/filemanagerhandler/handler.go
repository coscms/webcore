package filemanagerhandler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"

	uploadClient "github.com/webx-top/client/upload"
	uploadDropzone "github.com/webx-top/client/upload/driver/dropzone"

	"github.com/coscms/webcore/library/backend"
	"github.com/coscms/webcore/library/common"
	"github.com/coscms/webcore/library/config"
	"github.com/coscms/webcore/library/filemanager"
	"github.com/coscms/webcore/library/notice"
	"github.com/coscms/webcore/library/respond"
	"github.com/coscms/webcore/registry/upload/chunk"
)

func New(root, urlPrefix string) *FileManagerHandler {
	return &FileManagerHandler{
		root:      root,
		canUpload: true,
		canEdit:   true,
		canDelete: true,
		urlPrefix: urlPrefix,
	}
}

type FileManagerHandler struct {
	root      string
	canUpload bool
	canEdit   bool
	canDelete bool
	canChmod  bool
	canChown  bool
	urlPrefix string
}

func (h *FileManagerHandler) SetCanUpload(can bool) {
	h.canUpload = can
}
func (h *FileManagerHandler) SetCanEdit(can bool) {
	h.canEdit = can
}
func (h *FileManagerHandler) SetCanDelete(can bool) {
	h.canDelete = can
}
func (h *FileManagerHandler) SetCanChmod(can bool) {
	h.canChmod = can
}
func (h *FileManagerHandler) SetCanChown(can bool) {
	h.canChown = can
}

func (h FileManagerHandler) Handle(ctx echo.Context) error {
	filePath := ctx.Form(`path`)
	do := ctx.Form(`do`)
	mgr, err := filemanager.New(h.root, config.FromFile().Sys.EditableFileMaxBytes(), ctx)
	if err != nil {
		return err
	}
	if len(filePath) > 0 {
		filePath = filepath.Clean(`/` + filePath)
		filePath = strings.TrimPrefix(filePath, `/`)
	}

	user := backend.User(ctx)
	switch do {
	case `chmod`:
		if !h.canChmod {
			return echo.ErrNotFound
		}
		data := ctx.Data()
		perms := filemanager.Perms{}
		err = ctx.MustBind(&perms)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		err = mgr.Chmod(filePath, perms.Owner, perms.Group, perms.Other)
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		return ctx.JSON(data.SetCode(1))
	case `chown`:
		if !h.canChown {
			return echo.ErrNotFound
		}
		data := ctx.Data()
		uid := ctx.Formx(`uid`).Int()
		gid := ctx.Formx(`gid`).Int()
		if uid > 0 && gid > 0 {
			err = mgr.ChownByID(filePath, uid, gid)
		} else {
			username := ctx.Formx(`username`).String()
			usergroup := ctx.Formx(`usergroup`).String()
			if len(username) == 0 {
				return ctx.JSON(data.SetInfo(ctx.T(`请指定用户名`), 0))
			}
			err = mgr.Chown(filePath, username, usergroup)
		}
		if err != nil {
			return ctx.JSON(data.SetError(err))
		}
		return ctx.JSON(data.SetCode(1))
	case `edit`:
		if !h.canEdit {
			return echo.ErrNotFound
		}
		data := ctx.Data()
		if _, ok := Editable(filePath); !ok {
			data.SetInfo(ctx.T(`此文件不能在线编辑`), 0)
		} else {
			content := ctx.Form(`content`)
			encoding := ctx.Form(`encoding`)
			dat, err := mgr.Edit(filePath, content, encoding)
			if err != nil {
				data.SetInfo(err.Error(), 0)
			} else {
				data.SetData(dat, 1)
			}
		}
		return ctx.JSON(data)
	case `rename`:
		if !h.canEdit {
			return echo.ErrNotFound
		}
		data := ctx.Data()
		newName := ctx.Form(`name`)
		newName = echo.CleanFilePath(newName)
		newName = filepath.Join(filePath, newName)
		oldName := filePath
		err = mgr.Rename(oldName, newName)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetCode(1)
		}
		return ctx.JSON(data)
	case `mkdir`:
		if !h.canEdit {
			return echo.ErrNotFound
		}
		data := ctx.Data()
		newName := ctx.Form(`name`)
		newName = echo.CleanFilePath(newName)
		newName = filepath.Join(filePath, newName)
		err = mgr.Mkdir(newName, os.ModePerm)
		if err != nil {
			data.SetInfo(err.Error(), 0)
		} else {
			data.SetCode(1)
		}
		return ctx.JSON(data)
	case `delete`:
		if !h.canDelete {
			return echo.ErrNotFound
		}
		paths := ctx.FormValues(`path`)
		next := ctx.Referer()
		if len(next) == 0 {
			next = h.urlPrefix + com.URLEncode(filepath.Dir(filePath))
		}
		for _, filePath := range paths {
			filePath = strings.TrimSpace(filePath)
			if len(filePath) == 0 {
				continue
			}
			filePath = filepath.Clean(`/` + filePath)
			filePath = strings.TrimPrefix(filePath, `/`)
			err = mgr.Remove(filePath)
			if err != nil {
				common.SendFail(ctx, err.Error())
				return ctx.Redirect(next)
			}
		}
		return ctx.Redirect(next)
	case `upload`:
		if !h.canUpload {
			return echo.ErrNotFound
		}
		var cu *uploadClient.ChunkUpload
		var opts []uploadClient.ChunkInfoOpter
		if user != nil {
			cu = chunk.NewUploader(fmt.Sprintf(`user/%d`, user.Id))
			opts = append(opts, uploadClient.OptChunkInfoMapping(uploadDropzone.MappingChunkInfo))
		}
		err = mgr.Upload(filePath, cu, opts...)
		if err != nil {
			user := backend.User(ctx)
			if user != nil {
				notice.OpenMessage(user.Username, `upload`)
				notice.Send(user.Username, notice.NewMessageWithValue(`upload`, ctx.T(`文件上传出错`), err.Error()))
			}
		}
		return respond.Dropzone(ctx, err, nil)
	default:
		var dirs []os.FileInfo
		var exit bool
		err, exit, dirs = mgr.List(filePath)
		if exit {
			return err
		}
		ctx.Set(`dirs`, dirs)
	}
	if filePath == `.` {
		filePath = ``
	}
	pathSlice := strings.Split(strings.Trim(filePath, echo.FilePathSeparator), echo.FilePathSeparator)
	pathLinks := make(echo.KVList, len(pathSlice))
	encodedSep := filemanager.EncodedSepa
	urlPrefix := h.urlPrefix
	if !strings.HasSuffix(urlPrefix, encodedSep) {
		urlPrefix += encodedSep
	}
	for k, v := range pathSlice {
		urlPrefix += com.URLEncode(v)
		pathLinks[k] = &echo.KV{K: v, V: urlPrefix}
		urlPrefix += encodedSep
	}
	ctx.Set(`pathLinks`, pathLinks)
	ctx.Set(`rootPath`, strings.TrimSuffix(h.root, echo.FilePathSeparator))
	ctx.Set(`path`, filePath)
	ctx.Set(`absPath`, filepath.Join(mgr.Root.Name(), filePath))

	ctx.Set(`canUpload`, h.canUpload)
	ctx.Set(`canEdit`, h.canEdit)
	ctx.Set(`canDelete`, h.canDelete)

	ctx.SetFunc(`Editable`, func(fileName string) bool {
		if !h.canEdit {
			return false
		}
		_, ok := Editable(fileName)
		return ok
	})
	ctx.SetFunc(`Playable`, func(fileName string) string {
		mime, _ := Playable(fileName)
		return mime
	})
	return err
}
