package cayley

import (
	"os"
	"path/filepath"

	"github.com/cayleygraph/cayley/clog"
	"github.com/cayleygraph/cayley/graph"
	hkv "github.com/hidal-go/hidalgo/kv"
	"github.com/hidal-go/hidalgo/kv/bolt"
)

const Type = bolt.Name

func boltFilePath(path, filename string) string {
	return filepath.Join(path, filename)
}

func boltCreate(path string, opt graph.Options) (hkv.KV, error) {
	filename, err := opt.StringKey("filename", "indexes.bolt")
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(boltFilePath(path, filename), nil)
	if err != nil {
		clog.Errorf("Error: couldn't create Bolt database: %v", err)
		return nil, err
	}

	return db, nil
}

func boltOpen(path string, opt graph.Options) (hkv.KV, error) {
	filename, err := opt.StringKey("filename", "indexes.bolt")
	if err != nil {
		return nil, err
	}

	db, err := bolt.Open(boltFilePath(path, filename), nil)
	if err != nil {
		clog.Errorf("Error, couldn't open! %v", err)
		return nil, err
	}

	bdb := db.DB()
	bdb.NoSync, err = opt.BoolKey("nosync", false)
	if err != nil {
		db.Close()
		return nil, err
	}

	bdb.NoGrowSync = bdb.NoSync
	if bdb.NoSync {
		clog.Infof("Running in nosync mode")
	}

	return db, nil
}
