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
	// #cgo CFLAGS: -g -Wall -O2 -Wall -Wextra -Wno-unused-variable -Wno-unused-parameter -I. -I/usr/include/k2hash
	// #cgo LDFLAGS: -L/usr/lib -lk2hash
	// #include <stdlib.h>
	// #include <time.h>
	// #include "k2hash.h"
	// bool free_keypack(PK2HKEYPCK pkeys, int keycnt) {
	//     return k2h_free_keypack(pkeys, keycnt);
	// }
	// bool free_keypack(PK2HKEYPCK pkeys, int keycnt);
	"C"
)

import (
	"fmt"
	"os"
	"unsafe"
)

// K2hashParams stores parameters for k2hash C API
type K2hashParams struct {
	password           string
	expirationDuration int64
}

// K2hashQueue keeps a queue configurations.
type K2hashQueue struct {
	// K2HASH file handle
	handle C.k2h_h
	// K2HASH queue handle
	qhandle C.k2h_q_h
	// fifo
	fifo bool
	// prefix
	prefix string
}

// String returns a text representation of the object.
func (q *K2hashQueue) String() string {
	return fmt.Sprintf("[%v, %v, %v, %v]", q.handle, q.qhandle, q.fifo, q.prefix)
}

// NewK2hashQueue returns a new k2hash queue instance.
func NewK2hashQueue(h C.k2h_h, options ...func(*K2hashQueue)) (*K2hashQueue, error) {
	// 1. set defaults
	q := K2hashQueue{
		handle:  h,
		qhandle: C.K2H_INVALID_HANDLE,
		fifo:    true,
		prefix:  "",
	}
	// 2. set options
	for _, option := range options {
		option(&q)
	}
	// 3. open
	var qh C.k2h_q_h
	if q.prefix == "" {
		qh = C.k2h_q_handle(q.handle, C._Bool(q.fifo))
	} else {
		cPrefix := C.CBytes([]byte(q.prefix))
		defer C.free(unsafe.Pointer(cPrefix))
		qh = C.k2h_q_handle_prefix(q.handle, C._Bool(q.fifo), (*C.uchar)(cPrefix), C.size_t(len([]byte(q.prefix))))
	}
	// 4. check qeueu handle
	if qh == C.K2H_INVALID_HANDLE {
		return nil, fmt.Errorf("C.k2h_q_handle return false")
	}
	// 5. reset qhandle
	q.qhandle = qh
	return &q, nil
}

// Push adds a value to the queue.
func (q *K2hashQueue) Push(v interface{}, options ...func(*K2hashParams)) (bool, error) {
	// 1. binary or text
	var val string
	switch v.(type) {
	default:
		return false, fmt.Errorf("unsupported key data format %T", v)
	case string:
		val = v.(string)
	}
	// 2. set params
	params := K2hashParams{
		password:           "",
		expirationDuration: 0,
	}
	for _, option := range options {
		option(&params)
	}
	cPass := C.CString(params.password)
	defer C.free(unsafe.Pointer(cPass))
	var expire *C.time_t
	// WARNING: You can't set zero expire.
	if params.expirationDuration != 0 {
		expire = (*C.time_t)(&params.expirationDuration)
	}
	cVal := C.CString(val)
	defer C.free(unsafe.Pointer(cVal))
	if ok := C.k2h_q_str_push_wa(q.qhandle, (*C.char)(cVal), nil, 0, cPass, expire); !ok {
		return false, fmt.Errorf("C.k2h_q_str_push return false")
	}
	return true, nil
}

// Pop retrieves a value from the queue.
func (q *K2hashQueue) Pop(options ...func(*K2hashParams)) (string, error) {
	// 2. set params
	params := K2hashParams{
		password:           "",
		expirationDuration: 0,
	}
	for _, option := range options {
		option(&params)
	}
	cPass := C.CString(params.password)
	defer C.free(unsafe.Pointer(cPass))
	var cRetVal (*C.char)
	defer C.free(unsafe.Pointer(cRetVal))
	ok := C.k2h_q_str_pop_wa(q.qhandle, &cRetVal, nil, nil, cPass)
	defer C.free(unsafe.Pointer(cRetVal))
	if !ok {
		return "", fmt.Errorf("C.k2h_q_str_pop return false")
	}
	val := C.GoString(cRetVal)
	return val, nil
}

// Free destroys a k2hash queue handle.
func (q *K2hashQueue) Free() (bool, error) {
	ok := C.k2h_q_free(q.qhandle)
	if ok != true {
		return false, fmt.Errorf("k2h_close() returns false")
	}
	q.qhandle = C.K2H_INVALID_HANDLE
	return true, nil
}

// Count returns the number of k2hash queue elements.
func (q *K2hashQueue) Count() (int, error) {
	count := C.k2h_q_count(q.qhandle)
	return int(count), nil
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

func main() {
	// 1. define test data.
	k, err := NewK2hash("/tmp/sample.k2h")
	if err != nil {
		fmt.Printf("err %v¥n", err)
		os.Exit(1)
	}
	defer k.Close()

	q, err := NewK2hashQueue(k.handle)
	defer q.Free()
	fmt.Printf("q %v err %v\n", q.String(), err)
	// 1. count
	c, _ := q.Count()
	fmt.Printf("count %v\n", c)
	// 2. push
	q.Push("val")
	// 3. count
	c, _ = q.Count()
	fmt.Printf("count %v\n", c)
	// 4. pop
	s, _ := q.Pop()
	fmt.Printf("pop %v\n", s)
	// 5. count
	c, _ = q.Count()
	fmt.Printf("count %v\n", c)
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
