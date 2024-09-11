package file

import (
	"github.com/webx-top/com"
	"github.com/webx-top/echo/param"

	"github.com/coscms/webcore/dbschema"
)

func (f *Embedded) DeleteByInstance(m *dbschema.NgingFileEmbedded) error {
	err := f.Delete(nil, `id`, m.Id)
	if err != nil {
		return err
	}

	ids := param.Split(m.FileIds, ",").Uint64()
	return f.DeleteFileByIds(ids)
}

func (f *Embedded) AddFileByIds(fileIds []uint64, excludeFileIds ...uint64) error {
	var newIds []uint64
	if len(excludeFileIds) > 0 {
		for _, v := range fileIds {
			if !com.InUint64Slice(v, excludeFileIds) {
				newIds = append(newIds, v)
			}
		}
	} else {
		newIds = fileIds
	}
	if len(newIds) == 0 {
		return nil
	}
	return f.File.Incr(newIds...)
}

func (f *Embedded) DeleteFileByIds(fileIds []uint64, excludeFileIds ...uint64) error {
	var delIds []uint64
	if len(excludeFileIds) > 0 {
		for _, v := range fileIds {
			if !com.InUint64Slice(v, excludeFileIds) {
				delIds = append(delIds, v)
			}
		}
	} else {
		delIds = fileIds
	}
	if len(delIds) == 0 {
		return nil
	}
	err := f.File.Decr(delIds...)
	return err
}
