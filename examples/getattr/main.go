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

	// 2. Gets a text(string) type variable and retrives it by using k2h_get_str_value_wp and k2h_set_str_value_wa
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

	// 2. Sets a text(string) type variable and retrives it by using k2h_get_str_value_wp and k2h_set_str_value_wa
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

// AddAttr adds an attribute with a value to a key.
func (k2h *K2hash) AddAttr(k interface{}, ak interface{}, av interface{}) (bool, error) {
	// 1. binary or text
	var key string
	var attrkey string
	var attrval string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}
	switch ak.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", ak)
	case string:
		attrkey = ak.(string)
	}
	switch av.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", av)
	case string:
		attrval = av.(string)
	}

	// 2. get from k2hash get API
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cAttrKey := C.CString(attrkey)
	defer C.free(unsafe.Pointer(cAttrKey))
	cAttrVal := C.CString(attrval)
	defer C.free(unsafe.Pointer(cAttrVal))
	ok := C.k2h_add_str_attr(k2h.handle, (*C.char)(cKey), (*C.char)(cAttrKey), (*C.char)(cAttrVal))
	//fmt.Printf("ok %v\n", ok)
	if ok != true {
		return false, fmt.Errorf("C.k2h_add_str_attr return false")
	}
	return true, nil
}

// Attr holds attribute names and values.
type Attr struct {
	key string
	val string
}

// GetAttrs returns a slice of string
func (k2h *K2hash) GetAttrs(k string) ([]Attr, error) {
	// 1. retrieve an attribute using k2h_get_attrs
	// bool k2h_get_attrs(k2h_h handle, const unsigned char* pkey, size_t keylength, PK2HATTRPCK* ppattrspck, int* pattrspckcnt)
	cKey := C.CBytes([]byte(k))
	defer C.free(unsafe.Pointer(cKey))
	var attrpack C.PK2HATTRPCK
	var attrpackCnt C.int
	ok := C.k2h_get_attrs(
		k2h.handle,
		(*C.uchar)(cKey),
		C.size_t(len([]byte(k))+1), // plus one for a null termination
		&attrpack,
		&attrpackCnt,
	)
	defer C.k2h_free_attrpack(attrpack, attrpackCnt) // free the memory for the keypack for myself(GC doesn't know the area)

	if ok == false {
		fmt.Println("C.k2h_get_attrs returns false")
		return []Attr{}, fmt.Errorf("C.k2h_get_attrs() = %v", ok)
	} else if attrpackCnt == 0 {
		fmt.Printf("attrpackLen is zero")
		return []Attr{}, nil
	} else {
		fmt.Printf("attrpackLen is %v\n", attrpackCnt)
	}
	// 2. copy an attribute data to a slice
	var CAttrs C.PK2HATTRPCK = attrpack
	count := (int)(attrpackCnt)
	slice := (*[1 << 28]C.K2HATTRPCK)(unsafe.Pointer(CAttrs))[:count:count]
	fmt.Printf("slice size is %v\n", len(slice))
	//
	attrs := make([]Attr, count) // copy
	for i, data := range slice {
		// copy the data with len-1 length, which exclude a null termination.
		attrkey := C.GoBytes(unsafe.Pointer(data.pkey), (C.int)(data.keylength-1))
		fmt.Printf("i %v data %T pkey %v length %v attrkey %v\n", i, data, data.pkey, data.keylength, string(attrkey))
		attrval := C.GoBytes(unsafe.Pointer(data.pval), (C.int)(data.vallength-1))
		fmt.Printf("i %v data %T pval %v length %v attrval %v\n", i, data, data.pval, data.vallength, string(attrval))
		// cast bytes to a string
		attrs[i].key = string(attrkey)
		attrs[i].val = string(attrval)
	}
	return attrs, nil
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

	ok, err = k.AddAttr("set1", "attrkey1", "attrval1")
	fmt.Printf("AddAttr(sey1, attrkey1, attrval1) ok %v err %v\n", ok, err)

	attrs, err := k.GetAttrs("set1")
	if err == nil {
		for _, data := range attrs {
			fmt.Printf("GetAttrs(key1) ok key %v val %v\n", data.key, data.val)
		}
	} else {
		fmt.Printf("GetAttrs(key1) not ok %v\n", err)
	}
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
