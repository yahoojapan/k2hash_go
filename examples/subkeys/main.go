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
	// #cgo CFLAGS: -g -O2 -Wall -Wextra -Wno-unused-variable -Wno-unused-parameter -I. -I/usr/include/k2hash
	// #cgo LDFLAGS: -L/usr/lib -lk2hash
	// #include <stdlib.h>
	// #include <string.h>
	// #include "k2hash.h"
	//
	// static PK2HKEYPCK my_get_subkeypack(k2h_h handle, const char* pkey) {
	//      PK2HKEYPCK pskeypck = NULL;
	//      size_t len = strlen(pkey) + 1; // plus one for a null char
	//      int* pskeypckcnt = NULL;
	//      bool b = false;
	//      b = k2h_get_subkeys(handle, (const unsigned char*)pkey, len, &pskeypck, pskeypckcnt);
	//      if (b == false) {
	//          return NULL; // error
	//      }
	//      return pskeypck;
	// }
	"C"
)

import (
	"flag"
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

// K2hashParams stores parameters for k2hash C API
type K2hashParams struct {
	password           string
	expirationDuration int64
}

// String returns a text representation of the object.
func (p *K2hashParams) String() string {
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
	ok := C.k2h_close(k2h.handle)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	k2h.handle = C.K2H_INVALID_HANDLE
	return true, nil
}

// CloseWait closes a k2hash file after designated sleep time.
func (k2h *K2hash) CloseWait() (bool, error) {
	ok := C.k2h_close_wait(k2h.handle, C.long(k2h.waitms))
	if ok != true {
		return false, fmt.Errorf("k2h_close_wait() returns false")
	}
	k2h.handle = C.K2H_INVALID_HANDLE
	return true, nil
}

// GetResult holds the result of Get.
type GetResult struct {
	val []byte // text or binary
	ok  bool   // true if success
}

// Bytes returns the value in binary format.
func (r *GetResult) Bytes() []byte {
	return r.val
}

// String returns the value in text format.
func (r *GetResult) String() string {
	return string(r.val)
}

// Get returns data from a k2hash file.
func (k2h *K2hash) Get(k interface{}, options ...func(*K2hashParams)) (*GetResult, error) {
	// 1. binary or text
	var key string
	switch k.(type) {
	default:
		return nil, fmt.Errorf("unsupported key data format %T", key)
	case string:
		key = k.(string)
	}
	fmt.Printf("key %v\n", key)

	// 2. set params
	params := K2hashParams{
		password:           "",
		expirationDuration: 0,
	}

	for _, option := range options {
		option(&params)
	}

	// 3. get from k2hash get API
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var cRetValue *C.char
	defer C.free(unsafe.Pointer(cRetValue))
	//cRetValue = C.k2h_get_str_direct_value(k2h.handle, (*C.char)(cKey))
	ok := C.k2h_get_str_value(k2h.handle, (*C.char)(cKey), &cRetValue)
	rv := C.GoString(cRetValue)
	fmt.Printf("ok %v val %v\n", ok, rv)
	r := &GetResult{
		val: []byte(rv),
		ok:  (bool)(ok),
	}
	return r, nil
}

// Set returns true if successfully set a key with a value.
func (k2h *K2hash) Set(k interface{}, v interface{}, options ...func(*K2hashParams)) (bool, error) {
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
	fmt.Printf("key %v val %v\n", key, val)

	// 2. set params
	params := K2hashParams{
		password:           "",
		expirationDuration: 0,
	}

	for _, option := range options {
		option(&params)
	}

	// 3. Sets a text(string) type variable and retrives it by using k2h_get_str_value_wp and k2h_set_str_value_wa
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	cPass := C.CString(params.password)
	defer C.free(unsafe.Pointer(cPass))
	var expire *C.time_t
	// WARNING: You can't set zero expire.
	if params.expirationDuration != 0 {
		expire = (*C.time_t)(&params.expirationDuration)
	}
	ok := C.k2h_set_str_value_wa(k2h.handle, (*C.char)(cKey), (*C.char)(cVal), cPass, expire)
	fmt.Printf("ok %v\n", ok)
	if ok != true {
		return false, fmt.Errorf("C.k2h_set_str_value_wa return false")
	}
	return true, nil
}

// AddSubKey add a subkey to a key.
func (k2h *K2hash) AddSubKey(k interface{}, s interface{}, v interface{}, options ...func(*K2hashParams)) (bool, error) {
	// 1. binary or text
	var key string
	var subkey string
	var val string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}
	switch s.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", s)
	case string:
		subkey = s.(string)
	}
	switch v.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", v)
	case string:
		val = v.(string)
	}
	//fmt.Printf("key %v val %v\n", key, val)

	// 2. set params
	params := K2hashParams{
		password:           "",
		expirationDuration: 0,
	}

	for _, option := range options {
		option(&params)
	}

	// 3. get from k2hash get API
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cSubKey := C.CString(subkey)
	defer C.free(unsafe.Pointer(cSubKey))
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	cPass := C.CString(params.password)
	defer C.free(unsafe.Pointer(cPass))
	var expire *C.time_t
	// WARNING: You can't set zero expire.
	if params.expirationDuration != 0 {
		expire = (*C.time_t)(&params.expirationDuration)
	}
	ok := C.k2h_add_str_subkey_wa(k2h.handle, (*C.char)(cKey), (*C.char)(cSubKey), (*C.char)(cVal), cPass, expire)
	//fmt.Printf("ok %v\n", ok)
	if ok != true {
		return false, fmt.Errorf("C.C.k2h_add_str_subkey_wa return false")
	}
	return true, nil
}

// K2hashRemoveParams holds parameters to be removed.
type K2hashRemoveParams struct {
	all    bool
	subkey string
}

// Remove removes a key using a K2hashRemoveParams.
func (k2h *K2hash) Remove(k interface{}, options ...func(*K2hashRemoveParams)) (bool, error) {
	// 1. binary or text
	var key string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}
	fmt.Printf("key %v\n", key)

	// 2. remove params
	params := K2hashRemoveParams{
		all:    false,
		subkey: "",
	}

	for _, option := range options {
		option(&params)
	}

	// 3. get from k2hash get API
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	var ok C._Bool = false
	if params.all == true {
		fmt.Printf("C.k2h_remove_str_all\n")
		ok = C.k2h_remove_str_all(k2h.handle, (*C.char)(cKey))
	} else if params.subkey != "" {
		cSubKey := C.CString(params.subkey)
		defer C.free(unsafe.Pointer(cSubKey))
		fmt.Printf("C.k2h_remove_str_subkey\n")
		ok = C.k2h_remove_str_subkey(k2h.handle, (*C.char)(cKey), (*C.char)(cSubKey))
	} else {
		fmt.Printf("C.k2h_remove_str\n")
		ok = C.k2h_remove_str(k2h.handle, (*C.char)(cKey))
	}
	fmt.Printf("ok %v\n", ok)
	if ok != true {
		return false, fmt.Errorf("C.k2h_set_str_value_wa return false")
	}
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

// GetSubKeys retrieves subkeys to a key.
func (k2h *K2hash) GetSubKeys(k string) (bool, error) {
	cKey := C.CString(k)
	defer C.free(unsafe.Pointer(cKey))
	keypack := C.my_get_subkeypack(k2h.handle, (*C.char)(cKey))
	if keypack != nil {
		return false, fmt.Errorf("C.C.k2h_set_str_str_wa return false")
	}
	return true, nil
}

func main() {

	var loopCount int
	flag.IntVar(&loopCount, "loopCount", 1, "loopCount to detect memory leak")
	flag.Parse()
	fmt.Printf("loopCount %v\n", loopCount)

	// 1. define test data.
	k, err := NewK2hash("/tmp/sample.k2h")
	if err != nil {
		fmt.Printf("err %v¥n", err)
		os.Exit(1)
	}
	defer k.Close()

	// 1. key1
	ok, err := k.Set("key1", "val1")
	// 2. subkey1
	ok, err = k.Set("s1", "subval1")
	// 3. subkey2
	ok, err = k.Set("s2", "subval2")
	// 4. subkey3
	ok, err = k.Set("s3", "subval3")
	//for i := 0; i < loopCount; i++ { // memory leak check
	a := []string{"s1", "s2"}
	ok, err = k.SetSubKeys("key1", a)
	//}
	fmt.Printf("SetSubKeys(key1, []string{subkey1, subkey2}) ok %v err %v\n", ok, err)
	// 5. get
	val1, err := k.Get("key1")
	fmt.Printf("val1 should be val1(%v) err %v\n", val1.String(), err)
}

func main2() {
	// 1. define test data.
	k, err := NewK2hash("/tmp/sample.k2h")
	if err != nil {
		fmt.Printf("err %v¥n", err)
		os.Exit(1)
	}
	defer func() {
		ok, error := k.Close()
		if ok != true {
			fmt.Printf("Close() error %v\n", error)
		}
	}()

	// 1. set
	ok, err := k.Set("key1", "val1")
	fmt.Printf("Set() ok %v err %v\n", ok, err)
	// 2. get
	val1, err := k.Get("key1")
	fmt.Printf("val1 should be val1(%v) err %v\n", val1.String(), err)

	// 3. add a subkey
	ok, err = k.AddSubKey("key1", "sub1", "subval1")
	fmt.Printf("AddSubKey() ok %v err %v\n", ok, err)

	// 3.1. get the value of sub1
	subval1, err := k.Get("sub1")
	fmt.Printf("subval1 should be subval1(%v) err %v\n", subval1.String(), err)

	// 4. remove
	// 4.1. remove key1 and keep sub1
	ok, err = k.Remove("key1")
	fmt.Printf("Remove() ok %v err %v\n", ok, err)

	// 4.1.1. confirm key1 is removed
	val1, err = k.Get("key1")
	fmt.Printf("val1 should be removed(%v) err %v\n", val1.String(), err)

	// 4.1.2. confirm sub1 exists
	subval1, err = k.Get("sub1")
	fmt.Printf("subval1 should be subval1(%v) err %v\n", subval1.String(), err)

	// 4.2. remove all
	ok, err = k.Set("key1", "val1")
	ok, err = k.AddSubKey("key1", "sub1", "subval")
	opts := func(params *K2hashRemoveParams) {
		params.all = true
	}
	ok, err = k.Remove("key1", opts)
	fmt.Printf("Remove() ok %v err %v\n", ok, err)

	// 4.2.1. confirm key1 is removed
	val1, err = k.Get("key1")
	fmt.Printf("val1 should be removed(%v) err %v\n", val1.String(), err)

	// 4.2.2. confirm sub1 is removed
	subval1, err = k.Get("sub1")
	fmt.Printf("subval1 should be removed(%v) err %v\n", subval1.String(), err)

	// 4.3. remove subkey
	ok, err = k.Set("key1", "val1")
	ok, err = k.AddSubKey("key1", "sub1", "subval")
	opts = func(params *K2hashRemoveParams) {
		params.subkey = "sub1"
	}
	ok, err = k.Remove("key1", opts)
	fmt.Printf("ok %v err %v\n", ok, err)

	// 4.3.1. confirm key1 exists
	val1, err = k.Get("key1")
	fmt.Printf("val1 should be val1(%v) err %v\n", val1.String(), err)

	// 4.3.2. confirm sub1 is removed
	subval1, err = k.Get("sub1")
	fmt.Printf("subval1 should be removed(%v) err %v\n", subval1.String(), err)
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
