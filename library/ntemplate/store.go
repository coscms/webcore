package ntemplate

import (
	"encoding/json"

	"github.com/coscms/webcore/library/filecache"
	"github.com/webx-top/echo"
)

type Storer interface {
	Put(echo.Context, string, *ThemeInfo) error
	Get(echo.Context, string) (*ThemeInfo, error)
}

func NewFileStore(kind string) Storer {
	return &fileStore{kind: kind}
}

type fileStore struct {
	kind string
}

func (f *fileStore) Put(_ echo.Context, name string, v *ThemeInfo) error {
	b, _ := json.MarshalIndent(v, ``, `  `)
	return filecache.WriteCache(`themeinfo`, f.kind+`_`+name+`.json`, b)
}

func (f *fileStore) Get(_ echo.Context, name string) (*ThemeInfo, error) {
	themeInfo := NewThemeInfo()
	b, err := filecache.ReadCache(`themeinfo`, f.kind+`_`+name+`.json`)
	if err != nil {
		return themeInfo, err
	}
	err = json.Unmarshal(b, themeInfo)
	return themeInfo, err
}
