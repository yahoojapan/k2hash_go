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

// GetResult holds the result of Get.
type GetResult struct {
	val []byte // text or binary
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
func (k2h *K2hash) Get(k interface{}, options ...func(*Params)) (*GetResult, error) {
	// 1. binary or text
	var key string
	switch k.(type) {
	default:
		return nil, fmt.Errorf("unsupported key data format %T", key)
	case string:
		key = k.(string)
	}

	// 2. set params
	params := Params{
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
	cPass := C.CString(params.password)
	defer C.free(unsafe.Pointer(cPass))
	//cRetValue = C.k2h_get_str_direct_value(k2h.handle, (*C.char)(cKey))
	ok := C.k2h_get_str_value_wp(k2h.handle, (*C.char)(cKey), &cRetValue, cPass)
	if ok != true {
		return nil, fmt.Errorf("C.k2h_get_str_value return false")
	}
	rv := C.GoString(cRetValue)
	r := &GetResult{
		val: []byte(rv),
	}
	return r, nil
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
