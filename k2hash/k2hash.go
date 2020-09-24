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

package k2hash

import (
	// #cgo CFLAGS: -g -O2 -Wall -Wextra -Wno-unused-variable -Wno-unused-parameter -I. -I/usr/include/k2hash
	// #cgo LDFLAGS: -L/usr/lib -lk2hash
	// #include <stdlib.h>
	// #include "k2hash.h"
	"C"
)

import (
	"fmt"
	"unsafe"
)

// Params stores parameters for k2hash C API.
type Params struct {
	password           string
	expirationDuration int64
}

// QueueParams stores parameters for k2hash queue C API.
type QueueParams struct {
	password           string
	expirationDuration int64
	attrs              C.PK2HATTRPCK
}

// String returns a text representation of the object.
func (p *Params) String() string {
	return fmt.Sprintf("[%v, %v]", p.password, p.expirationDuration)
}

// K2hash keeps configurations, and it is responsible for creating request handles with a k2hash database files and closing them.
type K2hash struct {
	// filepath is a path to K2HASH file.
	filepath string
	// readonly enables read only file access if true. default is false.
	readonly bool
	// removefile enables automatic file deletion if no process attaches the file. default is false.
	removefile bool
	// fullmap enables memory mapped io for a whole data in a file. default is false.
	fullmap bool
	// key mask bit. default is 8.
	maskbitcnt int
	// key collision mask bit. default is 4.
	cmaskbitcnt int
	// maxelementcnt is a max number of duplicated elements if a hash collision occurs. default is 1024(bytes).
	maxelementcnt int
	// pagesize is a block size of data. default is 512(bytes).
	pagesize int
	// waitms is a time to wait until a transaction is completed. default is -1.
	waitms int
	// handle is a file descriptor to a K2HASH file.
	handle C.k2h_h
}

// String returns a text representation of the object.
func (k2h *K2hash) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v, %v, %v, %v, %v, %v, %v]",
		k2h.filepath, k2h.readonly, k2h.removefile, k2h.fullmap, k2h.maskbitcnt, k2h.cmaskbitcnt, k2h.maxelementcnt, k2h.pagesize, k2h.waitms, k2h.handle)
}

// NewK2hash returns a new k2hash instance.
func NewK2hash(f string, options ...func(*K2hash)) (*K2hash, error) {
	// 1. set defaults
	k2h := K2hash{
		filepath:      f,
		readonly:      false,
		removefile:    false,
		fullmap:       false,
		maskbitcnt:    8,
		cmaskbitcnt:   4,
		maxelementcnt: 1024,
		pagesize:      512,
		waitms:        0,
		handle:        0,
	}
	// 2. set options
	for _, option := range options {
		option(&k2h)
	}
	// 3. open
	ok, err := k2h.Open()
	if ok == false {
		return nil, err
	}
	return &k2h, nil
}

// Create creates a k2hash file. It returns a error if the file already exists.
func (k2h *K2hash) Create() bool {
	cK2h := C.CString(k2h.filepath)
	defer C.free(unsafe.Pointer(cK2h))
	ok := C.k2h_create(cK2h, C.int(k2h.maskbitcnt), C.int(k2h.cmaskbitcnt), C.int(k2h.maxelementcnt), C.ulong(k2h.pagesize))
	if ok == true {
		return true
	}
	return false
}

// Open opens a k2hash file.
func (k2h *K2hash) Open() (bool, error) {
	cK2h := C.CString(k2h.filepath)
	defer C.free(unsafe.Pointer(cK2h))
	k2h.handle = C.k2h_open(
		cK2h,
		C._Bool(k2h.readonly),
		C._Bool(k2h.removefile),
		C._Bool(k2h.fullmap),
		C.int(k2h.maskbitcnt),
		C.int(k2h.cmaskbitcnt),
		C.int(k2h.maxelementcnt),
		C.size_t(k2h.pagesize))

	if k2h.handle == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", k2h.handle)
	}
	return true, nil
}

// Close closes a k2hash file.
func (k2h *K2hash) Close() (bool, error) {
	ok := C.k2h_close_wait(k2h.handle, C.long(k2h.waitms))
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	k2h.handle = C.K2H_INVALID_HANDLE
	return true, nil
}

// GetHandle returns a k2hash file handle.
func (k2h *K2hash) GetHandle() C.k2h_h {
	return k2h.handle
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
