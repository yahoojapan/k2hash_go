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

// TxParams is a parameter set of transaction.
type TxParams struct {
	prefix             string
	params             string
	expirationDuration int64
}

// BeginTx enables transaction.
func (k2h *K2hash) BeginTx(file string, options ...func(*TxParams)) (bool, error) {
	params := TxParams{
		prefix:             "",
		params:             "",
		expirationDuration: 0,
	}

	for _, option := range options {
		option(&params)
	}

	// 3. remove
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	cPrefix := C.CBytes([]byte(params.prefix))
	defer C.free(unsafe.Pointer(cPrefix))
	cParams := C.CBytes([]byte(params.params))
	defer C.free(unsafe.Pointer(cParams))
	var expire *C.time_t
	// WARNING: You can't set zero expire.
	if params.expirationDuration != 0 {
		expire = (*C.time_t)(&params.expirationDuration)
	}
	ok := C.k2h_enable_transaction_param_we(k2h.handle, cFile, (*C.uchar)(cPrefix), (C.size_t)(len(params.prefix)+1), (*C.uchar)(cParams), (C.size_t)(len(params.params)+1), expire)
	if ok != true {
		return false, fmt.Errorf("C.k2h_enable_transaction_param_we return false")
	}
	return true, nil
}

// StopTx disables transaction.
func (k2h *K2hash) StopTx() (bool, error) {
	ok := C.k2h_disable_transaction(k2h.handle)
	if ok != true {
		return false, fmt.Errorf("C.k2h_disable_transaction return false")
	}
	return true, nil
}

// GetTxFileFD returns the file descriptor for transaction file.
func (k2h *K2hash) GetTxFileFD() int32 {
	fd := C.k2h_get_transaction_archive_fd(k2h.handle)
	return (int32)(fd)
}

// GetTxThreadPoolSize returns the number of thread pools.
func (k2h *K2hash) GetTxThreadPoolSize() (int32, error) {
	pool := C.k2h_get_transaction_thread_pool()
	if pool < 0 {
		return 0, fmt.Errorf("C.k2h_get_transaction_thread_pool return error")
	}
	return (int32)(pool), nil
}

// SetTxThreadPoolSize set the number of thread pools.
func (k2h *K2hash) SetTxThreadPoolSize(pool int32) (bool, error) {
	ok := C.k2h_set_transaction_thread_pool((C.int)(pool))
	if ok == false {
		return false, fmt.Errorf("C.k2h_set_transaction_thread_pool return false")
	}
	return true, nil
}

// UnsetTxThreadPoolSize set the number of thread pools zero.
func (k2h *K2hash) UnsetTxThreadPoolSize() (bool, error) {
	ok := C.k2h_unset_transaction_thread_pool()
	if ok == false {
		return false, fmt.Errorf("C.k2h_unset_transaction_thread_pool return false")
	}
	return true, nil
}

// LoadFromFile loads data from archive files.
func (k2h *K2hash) LoadFromFile(file string, ignoreError bool) (bool, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	ok := C.k2h_load_archive(k2h.handle, cFile, (C._Bool)(ignoreError))
	if ok == false {
		return false, fmt.Errorf("C.k2h_load_archive return false")
	}
	return true, nil
}

// DumpToFile saves data from a file as a serialized data.
func (k2h *K2hash) DumpToFile(file string, ignoreError bool) (bool, error) {
	cFile := C.CString(file)
	defer C.free(unsafe.Pointer(cFile))
	ok := C.k2h_put_archive(k2h.handle, cFile, (C._Bool)(ignoreError))
	if ok != true {
		return false, fmt.Errorf("k2h_put_archive() returns false")
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
