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
	"unsafe"

	"github.com/yahoojapan/k2hash_go/k2hash"
)

// The actual test functions are in non-_test.go files
// so that they can use cgo (import "C").
// These wrappers are here for gotest to find.

// testSet tests k2hash.Set method.
func testSetSubKeys(t *testing.T) {
	// 1. define test data.
	testData := []sks{
		{
			d: "binary",
			key: kv{
				k: []byte("setsubkeys1_parent"),
				v: []byte("p1"),
				p: "",
			},
			keys: []kv{
				{
					k: []byte("setsubkeys1_sub1"),
					v: []byte("s1"),
					p: "",
				},
				{
					k: []byte("setsubkeys1_sub2"),
					v: []byte("s2"),
					p: "",
				},
			},
		},
	}

	k, err := k2hash.NewK2hash("/tmp/test.k2h")
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) return err %v", err)
	}
	defer k.Close()
	for _, d := range testData {
		// 1. removes the current parent key if exists
		if ok, err := clearIfExists("/tmp/test.k2h", string(d.key.k)); !ok {
			t.Errorf("clearIfExists(%v, %v) = (%v, %v)", "/tmp/test.k2h", string(d.key.k), ok, err)
		}
		// 2. make a new parent key
		if ok, err := setKey("/tmp/test.k2h", d.key); !ok {
			t.Errorf("setKey(%v, %v) = (%v, %v)", "/tmp/test.k2h", d.key, ok, err)
		}
		// 3. make a child key
		skeys := make([]*C.char, len(d.keys)+1, len(d.keys)+1)
		for i, key := range d.keys {
			if ok, err := setKey("/tmp/test.k2h", key); !ok {
				t.Errorf("saveData(%q, %q, %q) = (%v, %v)", key.k, key.v, key.p, ok, err)
			}
			skeys[i] = C.CString(string(key.k))
			defer C.free(unsafe.Pointer(skeys[i]))
		}
		skeys[len(d.keys)] = nil
		// 4. link the child key to the parent key as a child key
		if ok, err := k.SetSubKeys(d.key, skeys); !ok {
			t.Errorf("SetSubKeys(%v, %v) = (%v, %v)", d.key, skeys, ok, err)
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
