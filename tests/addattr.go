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

// testAddAttr tests k2hash.AddAttr method
func testAddAttr(t *testing.T) {
	// 1. define test data.
	pkey := kv{
		d: "parent key",
		k: []byte("addsubkeys1_parent"),
		v: []byte("p1"),
		s: true,
		p: "",
		e: 0,
	}

	// 2. open a k2h file
	k, err := k2hash.NewK2hash("/tmp/test.k2h")
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) return err %v", err)
	}
	defer k.Close()

	// 3. set a key
	ok, err := k.Set(string(pkey.k), string(pkey.v))
	if ok == false {
		t.Errorf("k2hash.Set(%v, %v) return false. want true. err %v", pkey.k, pkey.v, err)
	}
	ok, val, err := getKey("/tmp/test.k2h", string(pkey.k))
	if ok != true {
		t.Errorf("getKey(%v, %v) = (%v, %v, %v)", "/tmp/test.k2h", string(pkey.k), ok, val, err)
	}
	if string(val) != string(pkey.v) {
		t.Errorf("getKey(%v, %v) = (%v, %v, %v), want %v", "/tmp/test.k2h", string(pkey.k), ok, val, err, string(pkey.v))
	}
	// 4. add a subkey
	ok, err = k.AddAttr(string(pkey.k), "sample_attrkey", "sample_attrval")
	if ok != true {
		t.Errorf("AddAttr(%v) = (%v, %v)", string(pkey.k), ok, err)
	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
