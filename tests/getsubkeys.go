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
	"strings"
	"testing"

	"github.com/yahoojapan/k2hash_go/k2hash"
)

// The actual test functions are in non-_test.go files
// so that they can use cgo (import "C").
// These wrappers are here for gotest to find.

// testGetSubKeys tests k2hash.GetSubKeys method.
func testGetSubKeys(t *testing.T) {
	// 1. define test data.
	testData := []sks{
		{
			d: "binary",
			key: kv{
				k: []byte("getsubkeys1_parent"),
				v: []byte("p1"),
				p: "",
			},
			keys: []kv{
				{
					k: []byte("getsubkeys1_sub1"),
					v: []byte("s1"),
					p: "",
				},
				{
					k: []byte("getsubkeys1_sub2"),
					v: []byte("s2"),
					p: "",
				},
			},
		},
	}
	// 2. exec tests
	for _, d := range testData {
		if ok, err := clearIfExists("/tmp/test.k2h", string(d.key.k)); !ok {
			t.Errorf("clearIfExists(%v, %v) = (%v, %v)", "/tmp/test.k2h", string(d.key.k), ok, err)
		}
		file, _ := k2hash.NewK2hash("/tmp/test.k2h")
		defer file.Close()
		if ok, err := file.Set(string(d.key.k), string(d.key.v)); !ok {
			t.Errorf("file.Set(%v, %v) = (%v, %v)", string(d.key.k), string(d.key.v), ok, err)
		}
		var skeys []string
		for i, key := range d.keys {
			if ok, err := file.Set(key.k, key.v); !ok {
				t.Errorf("file.Set(%v, %v) = (%v, %v)", key.k, key.v, ok, err)
			}
			skeys[i] = string(key.k)
		}
		if ok, err := file.SetSubKeys(d.key.k, skeys); !ok {
			t.Errorf("file.Set(%v, %v) = (%v, %v)", d.key.k, skeys, ok, err)
		}
		file.Close()
		testGetSubKeysArgs(string(d.key.k), skeys, t)
	}
}

func testGetSubKeysArgs(key string, skeys []string, t *testing.T) {
	// 1. Instantiate K2hash class
	file, _ := k2hash.NewK2hash("/tmp/test.k2h")
	defer file.Close()

	// 2. Get
	// 2.1. no args
	rskeys, err := file.GetSubKeys(key)
	if err != nil {
		t.Errorf("GetSubKeys(%v) = (%v, %v)", key, rskeys, err)
	}
	for i, key := range skeys {
		if strings.Compare(key, rskeys[i]) != 0 {
			t.Errorf("%v = %v", key, rskeys[i])
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
