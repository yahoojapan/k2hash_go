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

// testBumpDebugLevel tests k2hash.BumpDebugLevel method
func testBumpDebugLevel(t *testing.T) {
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
	f.BumpDebugLevel()
}

// testSetDebugLevelSilent tests k2hash.SetDebugLevelSilent method
func testSetDebugLevelSilent(t *testing.T) {
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
	f.SetDebugLevelSilent()
}

// testSetDebugLevelError tests k2hash.SetDebugLevelError method
func testSetDebugLevelError(t *testing.T) {
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
	f.SetDebugLevelError()
}

// testSetDebugLevelWarning tests k2hash.SetDebugLevelWarning method
func testSetDebugLevelWarning(t *testing.T) {
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
	f.SetDebugLevelWarning()
}

// testSetDebugLevelMessage tests k2hash.SetDebugLevelMessage method
func testSetDebugLevelMessage(t *testing.T) {
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
	f.SetDebugLevelMessage()
}

// testSetDebugFile tests k2hash.SetDebugFile method
func testSetDebugFile(t *testing.T) {
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
	f.SetDebugFile("/tmp/debug.log")
}

// testUnsetDebugFile tests k2hash.UnsetDebugFile method
func testUnsetDebugFile(t *testing.T) {
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
	f.UnsetDebugFile()
}

// testLoadDebugEnv tests k2hash.LoadDebugEnv method
func testLoadDebugEnv(t *testing.T) {
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
	f.LoadDebugEnv()
}

// testSetSignalUser1 tests k2hash.SetSignalUser1 method
func testSetSignalUser1(t *testing.T) {
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
	f.SetSignalUser1()
}

// testDumpHead tests k2hash.DumpHead method
func testDumpHead(t *testing.T) {
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
	f.DumpHead()
}

// testDumpKeyTable tests k2hash.DumpKeyTable method
func testDumpKeyTable(t *testing.T) {
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
	f.DumpKeyTable()
}

// testDumpFullKeyTable tests k2hash.DumpFullKeyTable method
func testDumpFullKeyTable(t *testing.T) {
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
	f.DumpFullKeyTable()
}

// testDumpElementTable tests k2hash.DumpElementTable method
func testDumpElementTable(t *testing.T) {
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
	f.DumpElementTable()
}

// testDumpFull tests k2hash.DumpFull method
func testDumpFull(t *testing.T) {
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
	f.DumpFull()
}

// testPrintState tests k2hash.PrintState method
func testPrintState(t *testing.T) {
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
	f.PrintState()
}

// testPrintVersion tests k2hash.PrintVersion method
func testPrintVersion(t *testing.T) {
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
	f.PrintVersion()
}

// testLoadHashLibrary tests k2hash.LoadHashLibrary method
func testLoadHashLibrary(t *testing.T) {
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
	f.LoadHashLibrary("/tmp/debug.log")
}

// testUnloadHashLibrary tests k2hash.UnloadHashLibrary method
func testUnloadHashLibrary(t *testing.T) {
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
	f.UnloadHashLibrary()
}

// testLoadTxLibrary tests k2hash.LoadTxLibrary method
func testLoadTxLibrary(t *testing.T) {
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
	f.LoadTxLibrary("/tmp/debug.log")
}

// testUnloadTxLibrary tests k2hash.UnloadTxLibrary method
func testUnloadTxLibrary(t *testing.T) {
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
	f.UnloadTxLibrary()
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
