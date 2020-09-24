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

// Rename renames an old key with a new key.
func (k2h *K2hash) Rename(o interface{}, n interface{}) (bool, error) {
	// 1. binary or text
	var old string
	var new string
	switch o.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", o)
	case string:
		old = o.(string)
	}
	switch n.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", n)
	case string:
		new = n.(string)
	}

	// 3. get from k2hash get API
	cOld := C.CString(old)
	defer C.free(unsafe.Pointer(cOld))
	cNew := C.CString(new)
	defer C.free(unsafe.Pointer(cNew))
	ok := C.k2h_rename_str(k2h.handle, (*C.char)(cOld), (*C.char)(cNew))
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
