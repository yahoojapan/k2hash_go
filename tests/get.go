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

// testGet tests k2hash.Get method.
func testGet(t *testing.T) {
	// 1. define test data.
	testData := []kv{
		{
			d: "test for string data",
			k: []byte("strkey"),
			v: []byte("strval"),
			s: true,
			p: "",
			e: 0,
		},
	}
	// 2. exec tests
	for _, d := range testData {
		if ok, err := clearIfExists("/tmp/test.k2h", string(d.k)); !ok {
			t.Errorf("clearIfExists(%v, %v = (%v, %v)", "/tmp/test.k2h", string(d.k), ok, err)
		}
		if ok, err := setKey("/tmp/test.k2h", d); !ok {
			t.Errorf("saveData(%v, %v) = (%v, %v)", "/tmp/test.k2h", d, ok, err)
		}
		testGetArgs(d, t)
	}
}

func testGetArgs(d kv, t *testing.T) {
	// 1. Instantiate K2hash class
	file, _ := k2hash.NewK2hash("/tmp/test.k2h")
	defer file.Close()

	// 2. Get
	if d.e == 0 && d.p == "" {
		// 2.1. no args
		val, err := file.Get(string(d.k))
		if err != nil {
			t.Errorf("Get(%v) = (%v, %v)", string(d.v), val.String(), err)
		}
		//fmt.Printf("val = %v, err = %v\n", val.String(), err)
		if val.String() != string(d.v) {
			t.Errorf("GetResult().String() = %v, want %v", val.String(), string(d.v))
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
