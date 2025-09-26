package leveldbpool

import "github.com/syndtr/goleveldb/leveldb"

type LevelDB[T comparable] interface {
	OpenDB(taskId T) (*leveldb.DB, error)
	CloseDB(taskId T)
	RemoveDB(taskId T) error
	CloseAllDB()
}
