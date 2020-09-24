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
		return []string{""}, fmt.Errorf("C.k2h_get_subkeys() = %v", ok)
	} else if keypackLen == 0 {
		return []string{""}, nil
	}
	// 2. copy a subkey data to a slice
	var CSubKey C.PK2HKEYPCK = keypack
	length := (int)(keypackLen)
	slice := (*[1 << 28]C.K2HKEYPCK)(unsafe.Pointer(CSubKey))[:length:length]
	skeys := make([]string, length) // copy
	for i, data := range slice {
		// copy the data with len-1 length, which exclude a null termination.
		sk := C.GoBytes(unsafe.Pointer(data.pkey), (C.int)(data.length-1))
		// cast byte array to string
		skeys[i] = string(sk)
	}
	return skeys, nil
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
