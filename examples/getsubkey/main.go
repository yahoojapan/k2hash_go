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
	// bool free_keypack(PK2HKEYPCK pkeys, int keycnt) {
	//     return k2h_free_keypack(pkeys, keycnt);
	// }
	// bool free_keypack(PK2HKEYPCK pkeys, int keycnt);
	"C"
)

import (
	"errors"
	"fmt"
	"os"
	"unsafe"
)

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
	ok := C.k2h_close(k2h.handle)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	k2h.handle = C.K2H_INVALID_HANDLE
	return true, nil
}

// Get returns data from a k2hash file.
func (k2h *K2hash) Get(k interface{}) string {
	// 1. binary or text
	var key string
	switch k.(type) {
	default:
		return ""
	case string:
		key = k.(string)
	}

	// 2. Sets a text(string) type variable and retrives it by using k2h_get_str_value_wp and k2h_set_str_value_wa
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var cRetValue *C.char
	defer C.free(unsafe.Pointer(cRetValue))
	ok := C.k2h_get_str_value(k2h.handle, (*C.char)(cKey), &cRetValue)
	if ok {
		return C.GoString(cRetValue)
	}
	return ""
}

// Set returns true if successfully set a key with a value.
func (k2h *K2hash) Set(k interface{}, v interface{}) (bool, error) {
	// 1. binary or text
	var key string
	var val string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}
	switch v.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", v)
	case string:
		val = v.(string)
	}

	// 2. Gets a text(string) type variable and retrives it by using k2h_get_str_value_wp and k2h_set_str_value_wa
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	ok := C.k2h_set_str_value_wa(k2h.handle, (*C.char)(cKey), (*C.char)(cVal), nil, nil)
	fmt.Printf("ok %v\n", ok)
	if ok != true {
		return false, fmt.Errorf("C.k2h_set_str_value_wa return false")
	}
	return true, nil
}

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

// SetSubKeys add subkeys to a key.
func (k2h *K2hash) SetSubKeys(k string, sk []string) (bool, error) {
	skeys := make([]*C.char, len(sk)+1, len(sk)+1)
	for i, s := range sk {
		fmt.Printf("i %v s %v\n", i, s)
		skeys[i] = C.CString(s)
		defer C.free(unsafe.Pointer(skeys[i]))
	}
	skeys[len(sk)] = nil
	// 3. get from k2hash get API
	cKey := C.CString(k)
	defer C.free(unsafe.Pointer(cKey))
	//ok := true
	if &skeys[0] != nil {
		ok := C.k2h_set_str_subkeys(k2h.handle, (*C.char)(cKey), (**C.char)(&skeys[0]))
		if ok != true {
			return false, fmt.Errorf("C.C.k2h_set_str_str_wa return false")
		}
	}
	return true, nil
}

// GetSubKeys returns subkeys to a key.
func (k2h *K2hash) GetSubKeys(k string) ([]string, error) {
	// 1. retrieve subkeys using k2h_get_subkeys
	cKey := C.CBytes([]byte(k))
	defer C.free(unsafe.Pointer(cKey))
	var keypack C.PK2HKEYPCK
	var keypackLen C.int
	ok := C.k2h_get_subkeys(
		k2h.handle,
		(*C.uchar)(cKey),
		C.size_t(len([]byte(k))+1), // plus one for a null termination
		&keypack,
		&keypackLen,
	)
	defer C.k2h_free_keypack(keypack, keypackLen) // free the memory for the keypack for myself(GC doesn't know the area)

	if ok == false {
		fmt.Println("C.k2h_get_subkeys returns false")
		return []string{""}, fmt.Errorf("C.k2h_get_subkeys() = %v", ok)
	} else if keypackLen == 0 {
		fmt.Printf("keypackLen is zero")
		return []string{""}, nil
	} else {
		fmt.Printf("keypackLen is %v\n", keypackLen)
	}
	// 2. copy a subkey data to a slice
	var CSubKey C.PK2HKEYPCK = keypack
	length := (int)(keypackLen)
	slice := (*[1 << 28]C.K2HKEYPCK)(unsafe.Pointer(CSubKey))[:length:length]
	fmt.Printf("slice size is %v\n", len(slice))
	skeys := make([]string, length) // copy
	for i, data := range slice {
		// copy the data with len-1 length, which exclude a null termination.
		sk := C.GoBytes(unsafe.Pointer(data.pkey), (C.int)(data.length-1))
		fmt.Printf("i %v data %T pkey %v length %v sk %v\n", i, data, data.pkey, data.length, string(sk))
		skeys[i] = string(sk)
	}
	return skeys, nil
}

func main() {
	// 1. define test data.
	k, err := NewK2hash("/tmp/sample.k2h")
	if err != nil {
		fmt.Printf("err %v¥n", err)
		os.Exit(1)
	}

	defer k.Close()
	ok, err := k.Set("key1", "val1")
	fmt.Printf("Set(key1) ok %v err %v\n", ok, err)

	ok, err = k.Set("set1", "setval1")
	fmt.Printf("Set(sey1) ok %v err %v\n", ok, err)

	ok, err = k.Set("set2", "setval2")
	fmt.Printf("Set(sey2) ok %v err %v\n", ok, err)

	ok, err = k.SetSubKeys("key1", []string{"set1", "set2"})
	fmt.Printf("SetSubKeys(key1, []string{subkey1, subkey2}) ok %v err %v\n", ok, err)

	skeys, err := k.GetSubKeys("key1")
	for i, key := range skeys {
		fmt.Printf("i %v key %v err %v\n", i, key, err)
	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
