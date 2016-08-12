/*
Copyright Mojing Inc. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package store

import (
	"errors"
	"os"
	"path"

	"github.com/tecbot/gorocksdb"
)

type RocksdbStorage struct {
	db         *gorocksdb.DB
	defaultCFH *gorocksdb.ColumnFamilyHandle
}

func NewRocksdbStorage(dbpath string) (*RocksdbStorage, error) {
	missing, err := dirMissingOrEmpty(dbpath)
	if err != nil {
		return nil, err
	}

	if missing {
		err = os.MkdirAll(path.Dir(dbpath), 0755)
		if err != nil {
			return nil, err
		}
	}

	opts := gorocksdb.NewDefaultOptions()
	defer opts.Destroy()

	opts.SetCreateIfMissing(missing)
	opts.SetCreateIfMissingColumnFamilies(true)

	cfNames := []string{"default"}
	var cfOpts []*gorocksdb.Options
	for range cfNames {
		cfOpts = append(cfOpts, opts)
	}

	db, cfHandlers, err := gorocksdb.OpenDbColumnFamilies(opts, dbpath, cfNames, cfOpts)
	if err != nil {
		return nil, err
	}

	//db, err := gorocksdb.OpenDb(opts, "default")
	//if err != nil {
	//	return nil, err
	//}

	return &RocksdbStorage{
		db:         db,
		defaultCFH: cfHandlers[0],
	}, nil
}

func (this *RocksdbStorage) Get(key []byte) ([]byte, error) {
	opt := gorocksdb.NewDefaultReadOptions()
	defer opt.Destroy()

	slice, err := this.db.GetCF(opt, this.defaultCFH, key)
	if err != nil {
		return nil, err
	}
	defer slice.Free()

	if slice.Data() == nil || len(slice.Data()) == 0 {
		return nil, errors.New("no data found")
	}

	data := makeCopy(slice.Data())

	//data, err := this.db.GetBytes(opt, key)
	//if err != nil {
	//	return nil, err
	//}
	//if data == nil || len(data) == 0 {
	//	return nil, errors.New("no data found")
	//}

	return data, nil
}

func (this *RocksdbStorage) Set(key []byte, value []byte) error {
	opt := gorocksdb.NewDefaultWriteOptions()
	defer opt.Destroy()

	return this.db.PutCF(opt, this.defaultCFH, key, value)
	//return this.db.Put(opt, key, value)
}

func (this *RocksdbStorage) Del(key []byte) error {
	opt := gorocksdb.NewDefaultWriteOptions()
	defer opt.Destroy()

	return this.db.DeleteCF(opt, this.defaultCFH, key)
	//return this.db.Delete(opt, key)
}

func (this *RocksdbStorage) Close() error {
	this.defaultCFH.Destroy()
	this.db.Close()

	return nil
}
