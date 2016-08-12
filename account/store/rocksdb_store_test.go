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
	"fmt"
	"gopkg.in/check.v1"
	"os"
	"path/filepath"
	"testing"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

type RocksdbStorageTest struct {
	storage *RocksdbStorage
}

var _ = check.Suite(&RocksdbStorageTest{})

func (t *RocksdbStorageTest) SetUpSuite(c *check.C) {
	var err error
	t.storage, err = NewRocksdbStorage(filepath.Join(os.TempDir(), "/test.db"))
	c.Assert(err, check.IsNil)
}

func (t *RocksdbStorageTest) TearDownSuite(c *check.C) {
	t.storage.Close()
	os.Remove(filepath.Join(os.TempDir(), "/test.db"))
}

func (t *RocksdbStorageTest) TestRocksdbStorage_Set_Get(c *check.C) {
	c.Assert(t.storage.Set([]byte("abc"), []byte("abc")), check.IsNil)
	get, err := t.storage.Get([]byte("abc"))
	c.Assert(err, check.IsNil)
	c.Assert(string(get), check.Equals, "abc")
	get, err = t.storage.Get([]byte("adb"))
	c.Assert(err, check.NotNil)
}

func (t *RocksdbStorageTest) TestRocksdbStorage_Del(c *check.C) {
	c.Assert(t.storage.Set([]byte("abc"), []byte("abc")), check.IsNil)
	c.Assert(t.storage.Del([]byte("abc")), check.IsNil)
	get, err := t.storage.Get([]byte("abc"))
	c.Logf("after delete, get value: %v, err: %v", string(get), err)
	c.Assert(string(get), check.Equals, "")
}

func (t *RocksdbStorageTest) BenchmarkRocksdbStorage_Set(c *check.C) {
	for i := 0; i < c.N; i++ {
		val := []byte(fmt.Sprintf("benchmarkRocksdb_%v", i))
		t.storage.Set(val, val)
	}
}

func (t *RocksdbStorageTest) BenchmarkRocksdbStorage_Get(c *check.C) {
	val := []byte("benchmarkRocksdb")
	t.storage.Set(val, val)

	for i := 0; i < c.N; i++ {
		t.storage.Get(val)
	}
}

func (t *RocksdbStorageTest) BenchmarkRocksdbStorage_Del(c *check.C) {
	val := []byte("benchmarkRocksdb")
	t.storage.Set(val, val)

	for i := 0; i < c.N; i++ {
		t.storage.Del(val)
	}
}
