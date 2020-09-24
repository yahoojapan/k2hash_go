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

// testGetAttrs tests k2hash.Get method.
func testGetAttrs(t *testing.T) {
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
			t.Errorf("clearIfExists(%v, %v) = (%v, %v)", "/tmp/test.k2h", string(d.k), ok, err)
		}
		k, err := k2hash.NewK2hash("/tmp/sample.k2h")
		if err != nil {
			t.Errorf("NewK2hash(%v) = (%v, %v)", "/tmp/test.k2h", k, err)
		}
		defer k.Close()
		if ok, err := k.Set(string(d.k), string(d.v)); !ok {
			t.Errorf("k.Set(%v, %v) = (%v, %v)", string(d.k), string(d.v), ok, err)
		}
		if ok, err := k.AddAttr(string(d.k), "attrkey1", "attrval1"); !ok {
			t.Errorf("AddAttr(%v, %v, %v) = (%v, %v)", string(d.k), "attrkey1", "attrval1", ok, err)
		}

		testGetAttrsArgs(d, "attrkey1", "attrval1", t)
	}
}

func testGetAttrsArgs(d kv, attrkey string, attrval string, t *testing.T) {
	// 1. Instantiate K2hash class
	file, _ := k2hash.NewK2hash("/tmp/test.k2h")
	defer file.Close()

	// 2. Get
	_, err := file.GetAttrs(string(d.k))
	if err != nil {
		//if strings.Compare(attrs[0].key, attrkey) != 0 {
		//	t.Errorf("GetAttrs(%v) = (%v, %v)", string(d.k), attrs, err)
		//}
		//if strings.Compare(attrs[0].val, attrval) != 0 {
		//	t.Errorf("GetAttrs(%v) = (%v, %v)", string(d.k), attrs, err)
		//}
	} else {
		t.Errorf("GetAttrs(%v) = (%v, %v)", string(d.k), "attrkey1", "attrval1")
	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
