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

// testBeginTx tests k2hash.BeginTx method
func testBeginTx(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. set
	f.BeginTx("/tmp/transaction.log")
}

// testStopTx tests k2hash.StopTx method
func testStopTx(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 3. print
	//f.PrintAttrInformation()
	// 2. set
	_, error := f.StopTx()
	if error != nil {
		t.Errorf("k2hash.StopTx() = ()")
	}
}

// testGetTxFileFD tests k2hash.GetTxFileFD method
func testGetTxFileFD(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.GetTxFileFD()
}

// testGetTxThreadPoolSize tests k2hash.GetTxThreadPoolSize method
func testGetTxThreadPoolSize(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.GetTxThreadPoolSize()
}

// testSetTxThreadPoolSize tests k2hash.SetTxThreadPoolSize method
func testSetTxThreadPoolSize(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.SetTxThreadPoolSize(1)
}

// testUnsetTxThreadPoolSize tests k2hash.UnsetTxThreadPoolSize method
func testUnsetTxThreadPoolSize(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.UnsetTxThreadPoolSize()
}

// testLoadFromFile tests k2hash.LoadFromFile method
func testLoadFromFile(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.LoadFromFile("/tmp/testarchive.k2h", true)
}

// testDumpToFile tests k2hash.DumpToFile method
func testDumpToFile(t *testing.T) {
	// 1. Instantiate K2hash class
	f, err := k2hash.NewK2hash("/tmp/test.k2h")
	if f == nil {
		defer f.Close()
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) = (nil, error), want not nil")
	}
	if err != nil {
		t.Errorf("k2hash.NewK2hash(/tmp/test.k2h) error %v", err)
	}

	// 2. print
	f.DumpToFile("/tmp/testarchive.k2h", true)
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
