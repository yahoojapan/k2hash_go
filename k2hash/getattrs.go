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

// Attr holds attribute names and values.
type Attr struct {
	key string
	val string
}

// String returns a text representation of the object.
func (r *Attr) String() string {
	return fmt.Sprintf("[%v, %v]", r.key, r.val)
}

// GetAttrs returns a slice of string.
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
		return []Attr{}, fmt.Errorf("C.k2h_get_attrs() = %v", ok)
	} else if attrpackCnt == 0 {
		return []Attr{}, nil
	}
	// 2. copy an attribute data to a slice
	var CAttrs C.PK2HATTRPCK = attrpack
	count := (int)(attrpackCnt)
	slice := (*[1 << 28]C.K2HATTRPCK)(unsafe.Pointer(CAttrs))[:count:count]
	attrs := make([]Attr, count) // copy
	for i, data := range slice {
		// copy the data with len-1 length, which exclude a null termination.
		attrkey := C.GoBytes(unsafe.Pointer(data.pkey), (C.int)(data.keylength-1))
		attrval := C.GoBytes(unsafe.Pointer(data.pval), (C.int)(data.vallength-1))
		// cast bytes to a string
		attrs[i].key = string(attrkey)
		attrs[i].val = string(attrval)
	}
	return attrs, nil
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
