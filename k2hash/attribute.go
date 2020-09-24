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

// EnableMtime enables the k2hash file attributes of modification time.
func (k2h *K2hash) EnableMtime(enable bool) (bool, error) {
	cBool := C._Bool(enable)
	ok := C.k2h_set_common_attr(k2h.handle, (*C._Bool)(&cBool), nil, nil, nil, nil)
	if ok != true {
		return false, fmt.Errorf("k2h_set_common_attr() returns false")
	}
	return true, nil
}

// EnableEncryption enables the k2hash file attributes of data encryption.
func (k2h *K2hash) EnableEncryption(enable bool, file string) (bool, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	cBool := C._Bool(enable)
	ok := C.k2h_set_common_attr(k2h.handle, nil, (*C._Bool)(&cBool), cFile, nil, nil)
	if ok != true {
		return false, fmt.Errorf("k2h_set_common_attr() returns false")
	}
	return true, nil
}

// EnableHistory enables the k2hash file attributes of history.
func (k2h *K2hash) EnableHistory(enable bool) (bool, error) {
	cBool := C._Bool(enable)
	ok := C.k2h_set_common_attr(k2h.handle, nil, nil, nil, (*C._Bool)(&cBool), nil)
	if ok != true {
		return false, fmt.Errorf("k2h_set_common_attr() returns false")
	}
	return true, nil
}

// SetExpirationDuration enables the k2hash file attributes of modification time.
func (k2h *K2hash) SetExpirationDuration(duration int) (bool, error) {
	cDuration := C.time_t(duration)
	ok := C.k2h_set_common_attr(k2h.handle, nil, nil, nil, nil, (*C.time_t)(&cDuration))
	if ok != true {
		return false, fmt.Errorf("k2h_set_common_attr() returns false")
	}
	return true, nil
}

// AddAttrPluginLibrary loads a shared library for processing attribute data.
func (k2h *K2hash) AddAttrPluginLibrary(file string) (bool, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	ok := C.k2h_add_attr_plugin_library(k2h.handle, cFile)
	if ok != true {
		return false, fmt.Errorf("k2h_add_attr_plugin_library() returns false")
	}
	return true, nil
}

// AddDecryptionPassword sets the decryption passphrase.
func (k2h *K2hash) AddDecryptionPassword(pass string) (bool, error) {
	cPass := C.CString(pass)
	defer C.free(unsafe.Pointer(cPass))
	ok := C.k2h_add_attr_crypt_pass(k2h.handle, cPass, false)
	if ok != true {
		return false, fmt.Errorf("k2h_add_attr_crypt_pass() returns false")
	}
	return true, nil
}

// SetDefaultEncryptionPassword sets the encryption passphrase.
func (k2h *K2hash) SetDefaultEncryptionPassword(pass string) (bool, error) {
	cPass := C.CString(pass)
	defer C.free(unsafe.Pointer(cPass))
	ok := C.k2h_add_attr_crypt_pass(k2h.handle, cPass, true)
	if ok != true {
		return false, fmt.Errorf("k2h_add_attr_crypt_pass() returns false")
	}
	return true, nil
}

// PrintAttrVersion prints attribute plugins to stderr.
func (k2h *K2hash) PrintAttrVersion() {
	C.k2h_print_attr_version(k2h.handle, nil)
}

// PrintAttrInformation prints attributes to stderr.
func (k2h *K2hash) PrintAttrInformation() {
	C.k2h_print_attr_information(k2h.handle, nil)
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
