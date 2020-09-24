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

type qdata struct {
	p []byte // prefix
	v []byte // val
	k []byte // key
	f bool   // fifo
	c bool   // copy attr if true
	s bool   // string or binary
	w string // password
	e int64  // duration of expiration(unit seconds)
}

// testKeyQueuePush tests K2hashQueue.Push method.
func testKeyQueuePush(t *testing.T) {
	// 1. define test data.
	testData := []qdata{
		// 1. default
		{
			p: []byte("push_prefix_2"),
			v: []byte("push_prefix_2_val"),
			k: nil,
			f: true,
			c: true,
			w: "",
			e: 0,
			s: false,
		},
		// 2. default + key
		{
			p: []byte("push_prefix_2"),
			v: []byte("push_prefix_2_val"),
			k: []byte("push_prefix_2_key"),
			f: true,
			c: true,
			w: "",
			e: 0,
			s: false,
		},
	}
	k, err := k2hash.NewK2hash("/tmp/test.k2h")
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) return err %v", err)
	}
	defer k.Close()

	for _, d := range testData {
		if ok, err := clearIfExists("/tmp/test.k2h", string(d.k)); !ok {
			t.Errorf("clearIfExists(%v, %v = (%v, %v)", "/tmp/test.k2h", string(d.k), ok, err)
		}
		testKeyQPushArgs(d, t)
	}
}

func testKeyQPushArgs(d qdata, t *testing.T) {
	k, err := k2hash.NewK2hash("/tmp/test.k2h")
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) return err %v", err)
	}
	defer k.Close()
	q, err := k2hash.NewKeyQueue(k)
	if err != nil {
		t.Errorf("k2hash.NewKeyQueue(%v) return err %v", k, err)
	}
	defer q.Free()
	// 1. push
	if ok, err := q.Push(string(d.v)); !ok {
		t.Errorf("KeyQueue.Push(%v) return false. wants true. err %v", string(d.k), err)
	}
	// 2. count
	if c, _ := q.Count(); c != 1 {
		t.Errorf("KeyQueue.Count() return not 1. wants 1")
	}
	// 3. pop
	if s, err := q.Pop(); s != string(d.v) {
		t.Errorf("KeyQueue.Pop( return %v. wants true. err %v", string(d.k), err)
	}

}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
