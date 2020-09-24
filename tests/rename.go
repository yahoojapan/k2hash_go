//
// k2hash_go
//
// Copyright 2018 Yahoo Japan Corporation.
//
// Go driver for k2hash that is a NoSQL Key Value Store(KVS) library.
// For k2hash, see https://github.com/yahoojapan/k2hash for the details.
//
// For the full copyright and license information, please view
// the license file that was distributed with this source code.
//
// AUTHOR:   Hirotaka Wakabayashi
// CREATE:   Fri, 14 Sep 2018
// REVISION:
//

package k2hashtest

import (
	// #cgo CFLAGS: -g -Wall -O2 -Wall -Wextra -Wno-unused-variable -Wno-unused-parameter -I. -I/usr/include/k2hash
	// #cgo LDFLAGS: -L/usr/lib -lk2hash
	// #include <stdlib.h>
	// #include "k2hash.h"
	"C"
)

import (
	"testing"

	"github.com/yahoojapan/k2hash_go/k2hash"
)

// The actual test functions are in non-_test.go files
// so that they can use cgo (import "C").
// These wrappers are here for gotest to find.

// testRename tests k2hash.Set method.
func testRename(t *testing.T) {
	// 1. define test data.
	testData := []kk{
		{
			d: "test for string data",
			o: []byte("oldkey"),
			v: []byte("value"),
			n: []byte("newkey"),
		},
	}
	k, err := k2hash.NewK2hash("/tmp/test.k2h")
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) return err %v", err)
	}
	defer k.Close()
	for _, d := range testData {
		// 1. set the old data.
		ok, err := k.Set(string(d.o), string(d.v))
		if ok == false {
			t.Errorf("k2hash.Set(%v, %v) return false. want true. err %v", string(d.o), string(d.v), err)
		}
		// 2. rename the old data with the new data
		ok, err = k.Rename(string(d.o), string(d.n))
		if ok == false {
			t.Errorf("k2hash.Rename(%v, %v) return false. want true. err %v", string(d.o), string(d.n), err)
		}
		ok, val, err := getKey("/tmp/test.k2h", string(d.n))
		if ok != true {
			t.Errorf("getKey(/tmp/test.k2h, %v) = (%v, %v, %v)", string(d.n), ok, val, err)
		}
	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
