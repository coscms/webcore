package sessionstore

import (
	_ "github.com/coscms/webcore/library/config/sessionstore/file"

	_ "github.com/coscms/webcore/library/config/sessionstore/redis"

	_ "github.com/coscms/webcore/library/config/sessionstore/bolt"

	_ "github.com/coscms/webcore/library/config/sessionstore/mysql"

	_ "github.com/coscms/webcore/library/config/sessionstore/sqlite"
)
