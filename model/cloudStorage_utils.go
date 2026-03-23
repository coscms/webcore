package model

import (
	"context"
	"errors"

	"github.com/coscms/webcore/model/file/storer"
	"github.com/webx-top/echo/defaults"
)

var (
	ErrNoSetStorerID = errors.New(`no set storerID`)
	ErrNoSetStorer   = errors.New(`no set storer`)
)

func GetCloudStorage(ctx context.Context, storerID string) (*CloudStorage, error) {
	eCtx := defaults.MustGetContext(ctx)
	m := NewCloudStorage(eCtx)
	cloudAccountID, err := storer.StorerIDToNumber(storerID)
	if err != nil {
		return nil, err
	}
	if cloudAccountID <= 0 {
		storerConfig, ok := storer.GetOk()
		if !ok {
			return nil, ErrNoSetStorer
		}
		cloudAccountID, err = storer.StorerIDToNumber(storerConfig.ID)
		if err != nil {
			return nil, err
		}
		if cloudAccountID <= 0 {
			return nil, ErrNoSetStorerID
		}
	}
	err = m.Get(nil, `id`, cloudAccountID)
	return m, err
}
