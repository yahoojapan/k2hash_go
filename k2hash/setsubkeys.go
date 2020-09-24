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

// SetSubKeys links another key as a child to the key.
func (k2h *K2hash) SetSubKeys(k interface{}, sk interface{}) (bool, error) {
	// 1. binary or text
	var key string
	var skeys []string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}
	switch sk.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", sk)
	case []string:
		skeys = sk.([]string)
	}

	// 2. set subkeys
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	cSkeys := make([]*C.char, len(skeys)+1, len(skeys)+1)
	for i, s := range skeys {
		cSkeys[i] = C.CString(s)
		defer C.free(unsafe.Pointer(cSkeys[i]))
	}
	cSkeys[len(skeys)] = nil

	ok := C.k2h_set_str_subkeys(k2h.handle, (*C.char)(cKey), (**C.char)(&cSkeys[0]))
	if ok != true {
		return false, fmt.Errorf("C.k2h_set_str_subkeys return false")
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
