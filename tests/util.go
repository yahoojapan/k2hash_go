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
	"errors"
	"fmt"
	"unsafe"
)

type kv struct {
	d string // description
	k []byte // key
	v []byte // val
	r bool   // remove subkeys if true.
	s bool   // true if use string format when saving the data.
	p string // password
	e int64  // expire
}
type kk struct {
	d string // description
	o []byte // new
	v []byte // key
	n []byte // old
}
type sk struct {
	d    string
	key  kv // key(parent)
	skey kv // subkey
}
type sks struct {
	d    string // description
	key  kv     // parentkey
	keys []kv   // subkeys
}

func clearIfExists(f string, key string) (bool, error) {
	cK2h := C.CString(f)
	defer C.free(unsafe.Pointer(cK2h))
	// 1. open
	handle := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handle == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handle)
	}

	// 2. remove
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ok := C.k2h_remove_str_all(
		handle,
		(*C.char)(cKey))
	if ok == false {
		return false, errors.New("C.k2h_remove_str_all returned false")
	}

	// 3. close
	ok = C.k2h_close(handle)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	handle = C.K2H_INVALID_HANDLE
	return true, nil
}

func setKey(f string, d kv) (bool, error) {
	if len(d.k) == 0 {
		return false, errors.New("len(k) == 0")
	}

	// 1. open
	cK2h := C.CString(f)
	defer C.free(unsafe.Pointer(cK2h))
	handle := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handle == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handle)
	}

	// 2. set_value
	cKey := C.CString(string(d.k))
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(string(d.v))
	defer C.free(unsafe.Pointer(cVal))
	cPass := C.CString(d.p)
	defer C.free(unsafe.Pointer(cPass))
	var expire *C.time_t
	// WARNING: You can't set zero expire.
	if d.e != 0 {
		expire = (*C.time_t)(&d.e)
	}
	ok := C.k2h_set_str_value_wa(handle, (*C.char)(cKey), (*C.char)(cVal), cPass, expire)
	if ok == false {
		return false, errors.New("C.k2h_set_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handle)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	handle = C.K2H_INVALID_HANDLE
	return true, nil
}

func getKey(f string, k string) (bool, []byte, error) {
	cK2h := C.CString(f)
	defer C.free(unsafe.Pointer(cK2h))

	// 1. open
	handle := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handle == C.K2H_INVALID_HANDLE {
		return false, nil, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handle)
	}

	// 2. get
	cKey := C.CString(k)
	defer C.free(unsafe.Pointer(cKey))
	var cRetValue *C.char
	defer C.free(unsafe.Pointer(cRetValue))
	//cRetValue = C.k2h_get_str_direct_value(k2h.handle, (*C.char)(cKey))
	ok := C.k2h_get_str_value(handle, (*C.char)(cKey), &cRetValue)
	rv := C.GoString(cRetValue)
	if ok == false {
		return false, nil, errors.New("C.k2h_set_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handle)
	if ok != true {
		return false, nil, fmt.Errorf("k2h_close() returns false")
	}
	handle = C.K2H_INVALID_HANDLE
	return true, []byte(rv), nil
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
