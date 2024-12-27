package prepare

import (
	"errors"
	"fmt"
	"io"
	"path"

	"github.com/webx-top/com"
	"github.com/webx-top/echo"
	"github.com/webx-top/echo/code"

	uploadLibrary "github.com/coscms/webcore/library/upload"
	modelFile "github.com/coscms/webcore/model/file"
	storerUtils "github.com/coscms/webcore/model/file/storer"
	"github.com/coscms/webcore/registry/upload"
	"github.com/coscms/webcore/registry/upload/checker"
	uploadChunk "github.com/coscms/webcore/registry/upload/chunk"
	"github.com/coscms/webcore/registry/upload/dbsaver"
	"github.com/coscms/webcore/registry/upload/driver"
	uploadClient "github.com/webx-top/client/upload"
)

func NewModel(ctx echo.Context, ownerType string, ownerID uint64, subdir string, fileType string, storerInfos ...storerUtils.Info) *modelFile.File {
	var storerInfo storerUtils.Info
	if len(storerInfos) > 0 {
		storerInfo = storerInfos[0]
	} else {
		storerInfo = storerUtils.Get()
	}
	fileM := modelFile.NewFile(ctx)
	fileM.StorerName = storerInfo.Name
	fileM.StorerId = storerInfo.ID
	fileM.OwnerId = ownerID
	fileM.OwnerType = ownerType
	fileM.Type = fileType
	fileM.Subdir = subdir
	return fileM
}

func NewClientWithModel(fileM *modelFile.File, clientName string, result *uploadClient.Result) uploadClient.Client {
	return NewClientWithResult(fileM.Context(), fileM.OwnerType, fileM.OwnerId, clientName, result)
}

func NewClient(ctx echo.Context, ownerType string, ownerID uint64, clientName string, fileType string) uploadClient.Client {
	result := &uploadClient.Result{
		FileType: uploadClient.FileType(fileType),
	}
	return NewClientWithResult(ctx, ownerType, ownerID, clientName, result)
}

func NewClientWithResult(ctx echo.Context, ownerType string, ownerID uint64, clientName string, result *uploadClient.Result) uploadClient.Client {
	client := uploadClient.Get(clientName)
	client.Init(ctx, result)
	cu := uploadChunk.NewUploader(fmt.Sprintf(`%s/%d`, ownerType, ownerID))
	client.SetChunkUpload(cu)
	uploadCfg := uploadLibrary.Get()
	client.SetReadBeforeHook(func(result *uploadClient.Result) error {
		extension := path.Ext(result.FileName)
		result.FileType = uploadClient.FileType(uploadCfg.DetectType(extension))
		return nil
	})
	return client
}

func NewStorer(ctx echo.Context, subdir string, storerInfos ...storerUtils.Info) (driver.Storer, error) {
	if len(subdir) == 0 {
		subdir = `default`
	}
	if !upload.AllowedSubdir(subdir) {
		return nil, ctx.NewError(code.InvalidParameter, `%s参数值“%s”未被登记`, `subdir`, subdir)
	}
	var storerInfo storerUtils.Info
	if len(storerInfos) > 0 {
		storerInfo = storerInfos[0]
	} else {
		storerInfo = storerUtils.Get()
	}
	//echo.Dump(ctx.Forms())
	newStore := driver.Get(storerInfo.Name)
	if newStore == nil {
		return nil, ctx.NewError(code.InvalidParameter, `存储引擎“%s”未被登记`, storerInfo.Name)
	}
	return newStore(ctx, subdir)
}

func DBSave(fileM *modelFile.File, subdir string, result *uploadClient.Result, originalReader io.Reader) error {
	dbSaverFn := dbsaver.Get(subdir)
	fileM.Id = 0
	fileM.SetByUploadResult(result)
	return dbSaverFn(fileM, result, originalReader)
}

type QuickConfig struct {
	FileType     string `validate:"required"`
	Subdir       string `validate:"required"`
	OwnerID      uint64 `validate:"required"`
	OwnerType    string `validate:"required"`
	Filename     string `validate:"required"`
	SaveFilename string
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

func CheckFileTypeString(fileType string) error {
	return CheckFileType(uploadClient.FileType(fileType))
}

var ErrInvalidFileType = errors.New(`invalid parameter: fileType`)

func CheckFileType(fileType uploadClient.FileType) error {
	_, ok := uploadClient.FileTypeExts[fileType]
	if ok {
		return nil
	}
	return fmt.Errorf(`%w (%s)`, ErrInvalidFileType, fileType)
}

// QuickStore 快存
func QuickStore(ctx echo.Context, cfg QuickConfig, r io.Reader, size int64) error {
	if err := cfg.Validate(ctx); err != nil {
		return err
	}
	stor, err := NewStorer(ctx, cfg.Subdir)
	if err != nil {
		return err
	}
	m := NewModel(ctx, cfg.OwnerType, cfg.OwnerID, cfg.Subdir, cfg.FileType)
	var subdir string
	subdir, _, err = checker.DefaultNoCheck(ctx)
	if err != nil {
		return err
	}
	var dstFile string
	dstFile, err = storerUtils.SaveFilename(subdir, cfg.SaveFilename, cfg.Filename)
	if err != nil {
		return err
	}
	result := uploadClient.Result{
		FileType: uploadClient.FileType(cfg.FileType),
		FileName: cfg.Filename,
		FileSize: size,
	}
	result.SavePath, result.FileURL, err = stor.Put(dstFile, r, size)
	if err != nil {
		return err
	}
	sk, ok := r.(io.Seeker)
	if ok {
		sk.Seek(0, 0)
	}
	result.Md5, err = com.Md5Reader(r)
	if err != nil {
		return err
	}
	if ok {
		sk.Seek(0, 0)
	}
	return DBSave(m, cfg.Subdir, &result, r)
}
