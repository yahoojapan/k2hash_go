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

// testEnableMtime tests k2hash.EnableMtime method
func testEnableMtime(t *testing.T) {
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
	f.EnableMtime(true)
	// 3. print
	//f.PrintAttrInformation()
}

// testEnableEncryption tests k2hash.EnableEncryption method
func testEnableEncryption(t *testing.T) {
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
	f.EnableEncryption(true, "/tmp/pass.txt")
	// 3. print
	//f.PrintAttrInformation()
}

// testEnableHistory tests k2hash.EnableHistory method
func testEnableHistory(t *testing.T) {
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
	f.EnableHistory(true)
	// 3. print
	//f.PrintAttrInformation()
}

// testSetExpirationDuration tests k2hash.SetExpirationDuration method
func testSetExpirationDuration(t *testing.T) {
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
	f.SetExpirationDuration(1)
	// 3. print
	//f.PrintAttrInformation()
}

// testAddDecryptionPassword tests k2hash.AddDecryptionPassword method
func testAddDecryptionPassword(t *testing.T) {
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
	f.AddDecryptionPassword("password")
	// 3. print
	//f.PrintAttrInformation()
}

// testSetDefaultEncryptionPassword tests k2hash.SetDefaultEncryptionPassword method
func testSetDefaultEncryptionPassword(t *testing.T) {
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
	f.SetDefaultEncryptionPassword("password")
	// 3. print
	//f.PrintAttrInformation()
}

// testPrintAttrVersion tests k2hash.PrintAttrVersion method
func testPrintAttrVersion(t *testing.T) {
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
	f.PrintAttrVersion()
}

// testPrintAttrInformation tests k2hash.PrintAttrInformation method
func testPrintAttrInformation(t *testing.T) {
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
	f.PrintAttrInformation()
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
