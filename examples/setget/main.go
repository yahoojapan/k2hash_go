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

package main

import (
	// #cgo CFLAGS: -g -Wall -O2 -Wall -Wextra -Wno-unused-variable -Wno-unused-parameter -I. -I/usr/include/k2hash
	// #cgo LDFLAGS: -L/usr/lib -lk2hash
	// #include <stdlib.h>
	// #include <time.h>
	// #include "k2hash.h"
	"C"
)

import (
	"errors"
	"fmt"
	"os"
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

func clearIfExists(d kv) (bool, error) {
	filepath := "/tmp/test.k2h"
	cK2h := C.CString(filepath)
	defer C.free(unsafe.Pointer(cK2h))
	// 1. open
	handler := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handler == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handler)
	}

	// 2. remove
	cKey := C.CBytes(d.k)
	defer C.free(cKey)
	ok := C.k2h_remove_all(
		handler,
		(*C.uchar)(cKey),
		C.size_t(len(d.k)))
	if ok == false {
		return false, errors.New("C.k2h_remove_all returned false")
	}

	// 3. close
	ok = C.k2h_close(handler)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	handler = C.K2H_INVALID_HANDLE
	return true, nil
}

func setKey(d kv) (bool, error) {
	if len(d.k) == 0 {
		return false, errors.New("len(k) == 0")
	}

	// 1. open
	filepath := "/tmp/test.k2h"
	cK2h := C.CString(filepath)
	defer C.free(unsafe.Pointer(cK2h))
	handler := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handler == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handler)
	}

	// 2. set_value
	cKey := C.CBytes(d.k)
	defer C.free(cKey)
	cVal := C.CBytes(d.v)
	defer C.free(cVal)
	cPass := C.CString(d.p)
	defer C.free(unsafe.Pointer(cPass))
	expire := (C.time_t)(60) // expire in 60sec.
	ok := C.k2h_set_value_wa(handler, (*C.uchar)(cKey), C.size_t(len(d.k)), (*C.uchar)(cVal), C.size_t(len(d.v)), cPass, &expire)
	if ok == false {
		return false, errors.New("C.k2h_set_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handler)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	handler = C.K2H_INVALID_HANDLE
	return true, nil
}

func setKeyByString(d kv) (bool, error) {
	filepath := "/tmp/test.k2h"
	cK2h := C.CString(filepath)
	defer C.free(unsafe.Pointer(cK2h))

	// 1. open
	handler := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handler == C.K2H_INVALID_HANDLE {
		return false, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handler)
	}

	// 2. get
	cKey := C.CString(string(d.k)) // d.k's type is a Go's byte array. string(d.k)'s type is a Go String, C.CString(string(d.k))'s type is *C.char
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(string("text")) // Cast a Go string to a *C.char
	defer C.free(unsafe.Pointer(cVal))
	cPass := C.CString(d.p)
	defer C.free(unsafe.Pointer(cPass))
	expire := (C.time_t)(60) // expire in 60sec.
	ok := C.k2h_set_str_value_wa(handler, cKey, cVal, cPass, &expire)
	if ok == false {
		return false, errors.New("C.k2h_set_str_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handler)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	handler = C.K2H_INVALID_HANDLE
	return true, nil
}

func getKey(d kv) (bool, []byte, error) {
	filepath := "/tmp/test.k2h"
	cK2h := C.CString(filepath)
	defer C.free(unsafe.Pointer(cK2h))

	// 1. open
	handler := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handler == C.K2H_INVALID_HANDLE {
		return false, nil, fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handler)
	}

	// 2. get
	cKey := C.CBytes(d.k)
	defer C.free(cKey)
	var cRetValue *C.uchar  // value:(*main._Ctype_char)(nil) type:*main._Ctype_char
	var valLen C.size_t     // valLen value:0x0 type:main._Ctype_size_t
	cPass := C.CString(d.p) // func C.CString(string) *C.char
	defer C.free(unsafe.Pointer(cPass))
	ok := C.k2h_get_value_wp(
		handler,
		(*C.uchar)(cKey),
		C.size_t(len(d.k)),
		&cRetValue,
		&valLen,
		cPass)
	defer C.free(unsafe.Pointer(cRetValue))
	val := C.GoBytes(unsafe.Pointer(cRetValue), C.int(valLen))
	if ok == false {
		return false, nil, errors.New("C.k2h_set_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handler)
	if ok != true {
		return false, nil, fmt.Errorf("k2h_close() returns false")
	}
	handler = C.K2H_INVALID_HANDLE
	return true, val, nil
}

func getKeyByString(d kv) (bool, string, error) {
	filepath := "/tmp/test.k2h"
	cK2h := C.CString(filepath)
	defer C.free(unsafe.Pointer(cK2h))

	// 1. open
	handler := C.k2h_open(
		cK2h,
		C._Bool(false),
		C._Bool(false),
		C._Bool(false),
		C.int(8),
		C.int(4),
		C.int(1024),
		C.size_t(512))
	if handler == C.K2H_INVALID_HANDLE {
		return false, "", fmt.Errorf("k2h_open() returns K2H_INVALID_HANDLE(%v)", handler)
	}

	// 2. get
	cKey := C.CString(string(d.k))
	defer C.free(unsafe.Pointer(cKey))
	var cRetValue *C.char   // value:(*main._Ctype_char)(nil) type:*main._Ctype_char
	cPass := C.CString(d.p) // func C.CString(string) *C.char
	defer C.free(unsafe.Pointer(cPass))
	ok := C.k2h_get_str_value_wp(
		handler,
		cKey,
		&cRetValue,
		cPass)
	defer C.free(unsafe.Pointer(cRetValue))
	val := C.GoString(cRetValue)
	if ok == false {
		return false, "", errors.New("C.k2h_set_value_wa returned false")
	}

	// 3. close
	ok = C.k2h_close(handler)
	if ok != true {
		return false, "", fmt.Errorf("k2h_close() returns false")
	}
	handler = C.K2H_INVALID_HANDLE
	return true, val, nil
}

func main() {
	fmt.Printf("params\n")
	// 1. define test data.
	var testData = []kv{
		{d: "binary", k: []byte("set1"), v: []byte("bin"), p: "", e: 0, r: false, s: false},
		{d: "string", k: []byte("set2"), v: []byte("bin"), p: "", e: 0, r: false, s: true},
	}
	for _, d := range testData {
		if ok, err := clearIfExists(d); !ok {
			fmt.Printf("error %v\n", err)
			os.Exit(1)
		}
		if ok, err := setKey(d); !ok {
			fmt.Printf("error %v\n", err)
			os.Exit(1)
		}
		ok, val, err := getKey(d)
		if !ok {
			fmt.Printf("error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("getKey %v val %v\n", string(d.k), string(val))
		if ok, err = setKeyByString(d); !ok {
			fmt.Printf("error %v\n", err)
			os.Exit(1)
		}
		ok, valString, err := getKeyByString(d)
		if !ok {
			fmt.Printf("error %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("getKeyByString %v valString %v\n", string(d.k), valString)

	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
