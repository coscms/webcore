package prepare

import (
	"io"

	storerUtils "github.com/coscms/webcore/model/file/storer"
	"github.com/coscms/webcore/registry/upload"
	"github.com/coscms/webcore/registry/upload/checker"
	uploadClient "github.com/webx-top/client/upload"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"
)

type QuickConfig struct {
	FileType     string `validate:"required"`
	Subdir       string `validate:"required"`
	OwnerID      uint64 `validate:"required"`
	OwnerType    string `validate:"required"`
	Filename     string `validate:"required"`
	SaveFilename string
	result       *uploadClient.Result
}

func (qc QuickConfig) Validate(ctx echo.Context) error {
	if !upload.AllowedSubdir(qc.Subdir) {
		return ctx.NewError(code.InvalidParameter, `%s参数值“%s”未被登记`, `subdir`, qc.Subdir)
	}
	err := CheckFileTypeString(qc.FileType)
	if err != nil {
		return err
	}
	if len(qc.OwnerType) == 0 {
		return ctx.NewError(code.InvalidParameter, `%s参数值“%s”无效`, `ownerType`, qc.OwnerType)
	}
	if qc.OwnerID == 0 {
		return ctx.NewError(code.InvalidParameter, `%s参数值“%s”无效`, `ownerID`, qc.OwnerID)
	}
	if len(qc.Filename) == 0 {
		return ctx.NewError(code.InvalidParameter, `%s参数值“%s”无效`, `filename`, qc.Filename)
	}
	return err
}

// Save 快存
func (qc *QuickConfig) Save(ctx echo.Context, r io.ReadSeeker, size int64) error {
	if err := qc.Validate(ctx); err != nil {
		return err
	}
	stor, err := NewStorer(ctx, qc.Subdir)
	if err != nil {
		return err
	}
	m := NewModel(ctx, qc.OwnerType, qc.OwnerID, qc.Subdir, qc.FileType)
	var subdir string
	subdir, _, err = checker.DefaultNoCheck(ctx)
	if err != nil {
		return err
	}
	var dstFile string
	dstFile, err = storerUtils.SaveFilename(subdir, qc.SaveFilename, qc.Filename)
	if err != nil {
		return err
	}
	result := uploadClient.Result{
		FileType: uploadClient.FileType(qc.FileType),
		FileName: qc.Filename,
		FileSize: size,
	}
	result.SavePath, result.FileURL, err = stor.Put(dstFile, r, size)
	if err != nil {
		return err
	}
	r.Seek(0, 0)
	result.Md5, err = com.Md5Reader(r)
	if err != nil {
		return err
	}
	r.Seek(0, 0)
	qc.result = &result
	return DBSave(m, qc.Subdir, qc.result, r)
}

func (qc QuickConfig) Result() *uploadClient.Result {
	return qc.result
}
