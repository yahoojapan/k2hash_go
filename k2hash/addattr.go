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

// AddAttr adds an attribute with a value to a key.
func (k2h *K2hash) AddAttr(k interface{}, ak interface{}, av interface{}, options ...func(*Params)) (bool, error) {
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
	if ok != true {
		return false, fmt.Errorf("C.k2h_add_str_attr return false")
	}
	return true, nil
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
