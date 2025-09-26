package leveldbpool

import (
	"os"
	"path/filepath"
	"regexp"
	"sync"

	"github.com/syndtr/goleveldb/leveldb"
	levelDBErrors "github.com/syndtr/goleveldb/leveldb/errors"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/webx-top/com"
	"github.com/webx-top/echo"
)

func New[T comparable](dir string) *dbPool[T] {
	return &dbPool[T]{mp: map[T]*leveldb.DB{}, dir: dir}
}

type dbPool[T comparable] struct {
	mu  sync.RWMutex
	mp  map[T]*leveldb.DB
	dir string
}

func (t *dbPool[T]) OpenDB(taskId T) (*leveldb.DB, error) {
	t.mu.RLock()
	db := t.mp[taskId]
	t.mu.RUnlock()

	if db == nil {
		t.mu.Lock()
		defer t.mu.Unlock()
		var err error
		db, err = t.openLevelDB(taskId)
		if err != nil {
			return nil, err
		}
		t.mp[taskId] = db
	}
	return db, nil
}

var regexBadManifest = regexp.MustCompile(`manifest corrupted(?:.*): missing`)

func (t *dbPool[T]) openLevelDB(taskId T) (*leveldb.DB, error) {
	idKey := com.String(taskId)
	dbDir := filepath.Join(echo.Wd(), t.dir)
	err := com.MkdirAll(dbDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	dbFile := filepath.Join(dbDir, idKey)
	var db *leveldb.DB
	var options *opt.Options
	db, err = leveldb.OpenFile(dbFile, options)
	if err != nil && (levelDBErrors.IsCorrupted(err) || regexBadManifest.MatchString(err.Error())) {
		db, err = leveldb.RecoverFile(dbFile, options)
	}
	return db, err
}

func (t *dbPool[T]) removeLevelDB(taskId T) error {
	idKey := com.String(taskId)
	dbFile := filepath.Join(echo.Wd(), t.dir, idKey)
	if !com.FileExists(dbFile) {
		return nil
	}
	return os.RemoveAll(dbFile)
}

func (t *dbPool[T]) CloseDB(taskId T) {
	t.mu.Lock()
	if db, ok := t.mp[taskId]; ok {
		db.Close()
		delete(t.mp, taskId)
	}
	t.mu.Unlock()
}

func (t *dbPool[T]) RemoveDB(taskId T) error {
	t.mu.Lock()
	if db, ok := t.mp[taskId]; ok {
		db.Close()
		delete(t.mp, taskId)
	}
	err := t.removeLevelDB(taskId)
	t.mu.Unlock()
	return err
}

func (t *dbPool[T]) CloseAllDB() {
	t.mu.Lock()
	for _, db := range t.mp {
		db.Close()
	}
	t.mu.Unlock()
}
