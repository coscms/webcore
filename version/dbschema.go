package version

import "github.com/coscms/webcore/library/config"

const (
	// 数据表结构版本
	DBSCHEMA = 6.700
)

func init() {
	config.Version.SetPkgDBSchemas(`webcore`, DBSCHEMA)
}
