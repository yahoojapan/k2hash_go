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

// RemoveParams is a parameter set of remove.
type RemoveParams struct {
	all    bool
	subkey string
}

// Remove removes a key or a subkey.
func (k2h *K2hash) Remove(k interface{}, options ...func(*RemoveParams)) (bool, error) {
	// 1. binary or text
	var key string
	switch k.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", k)
	case string:
		key = k.(string)
	}

	// 2. remove params
	params := RemoveParams{
		all:    false,
		subkey: "",
	}

	for _, option := range options {
		option(&params)
	}

	// 3. remove
	cKey := C.CString(key)
	defer C.free(unsafe.Pointer(cKey))
	ok := C._Bool(false)
	if params.all == true {
		ok = C.k2h_remove_str_all(k2h.handle, (*C.char)(cKey))
	} else if params.subkey != "" {
		cSubKey := C.CString(params.subkey)
		defer C.free(unsafe.Pointer(cSubKey))
		ok = C.k2h_remove_str_subkey(k2h.handle, (*C.char)(cKey), (*C.char)(cSubKey))
	} else {
		ok = C.k2h_remove_str(k2h.handle, (*C.char)(cKey))
	}
	if ok != true {
		return false, fmt.Errorf("C.k2h_set_str_value_wa return false")
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
