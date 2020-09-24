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

// BumpDebugLevel changes the log level.
func (k2h *K2hash) BumpDebugLevel() {
	C.k2h_bump_debug_level()
}

// SetDebugLevelSilent disables logging.
func (k2h *K2hash) SetDebugLevelSilent() {
	C.k2h_set_debug_level_silent()
}

// SetDebugLevelError logs error level messages.
func (k2h *K2hash) SetDebugLevelError() {
	C.k2h_set_debug_level_error()
}

// SetDebugLevelWarning logs warning or higher level messages.
func (k2h *K2hash) SetDebugLevelWarning() {
	C.k2h_set_debug_level_warning()
}

// SetDebugLevelMessage logs info or higher level messages.
func (k2h *K2hash) SetDebugLevelMessage() {
	C.k2h_set_debug_level_message()
}

// SetDebugFile defines the file where log messages are saved.
func (k2h *K2hash) SetDebugFile(filepath string) bool {
	cFilePath := C.CString(filepath)
	ok := C.k2h_set_debug_file(cFilePath) // ok is a C._Bool type, which is aliased by bool.
	if ok == true {
		return true
	}
	return false
}

// UnsetDebugFile disables saving log messages to a file.
func (k2h *K2hash) UnsetDebugFile() bool {
	ok := C.k2h_unset_debug_file()
	if ok == true {
		return true
	}
	return false
}

// LoadDebugEnv defines the log level and the file by using K2HDBGMODE, K2HDBGFILE.
func (k2h *K2hash) LoadDebugEnv() bool {
	ok := C.k2h_load_debug_env()
	if ok == true {
		return true
	}
	return false
}

// SetSignalUser1 changes log level by receiving SIGUSR1 signal.
func (k2h *K2hash) SetSignalUser1() bool {
	ok := C.k2h_set_bumpup_debug_signal_user1()
	if ok == true {
		return true
	}
	return false
}

// DumpHead dumps K2HASH header information to a file referred by FILE pointer.
func (k2h *K2hash) DumpHead() bool {
	ok := C.k2h_dump_head(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// DumpKeyTable dumps K2HASH's hash key table information to a file referred by FILE pointer.
func (k2h *K2hash) DumpKeyTable() bool {
	ok := C.k2h_dump_keytable(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// DumpFullKeyTable dumps K2HASH's hash key and subkey information to a file referred by FILE pointer.
func (k2h *K2hash) DumpFullKeyTable() bool {
	ok := C.k2h_dump_full_keytable(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// DumpElementTable dumps K2HASH's Element table information to a file referred by FILE pointer.
func (k2h *K2hash) DumpElementTable() bool {
	ok := C.k2h_dump_elementtable(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// DumpFull dumps K2HASH's all information.
func (k2h *K2hash) DumpFull() bool {
	ok := C.k2h_dump_full(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// PrintState prints k2hash file stats.
func (k2h *K2hash) PrintState() bool {
	ok := C.k2h_print_state(k2h.handle, nil)
	if ok == true {
		return true
	}
	return false
}

// PrintVersion prints the k2hash library version and the credit.
func (k2h *K2hash) PrintVersion() {
	C.k2h_print_version(nil)
}

// LoadHashLibrary loads the hash library.
func (k2h *K2hash) LoadHashLibrary(filepath string) bool {
	cFilePath := C.CString(filepath)
	ok := C.k2h_load_hash_library(cFilePath) // ok is a C._Bool type, which is aliased by bool.
	if ok == true {
		return true
	}
	return false
}

// UnloadHashLibrary unloads the loaded hash library.
func (k2h *K2hash) UnloadHashLibrary() bool {
	ok := C.k2h_unload_hash_library()
	if ok == true {
		return true
	}
	return false
}

// LoadTxLibrary loads the hash library.
func (k2h *K2hash) LoadTxLibrary(filepath string) bool {
	cFilePath := C.CString(filepath)
	ok := C.k2h_load_transaction_library(cFilePath) // ok is a C._Bool type, which is aliased by bool.
	if ok == true {
		return true
	}
	return false
}

// UnloadTxLibrary unloads the loaded hash library.
func (k2h *K2hash) UnloadTxLibrary() bool {
	ok := C.k2h_unload_transaction_library()
	if ok == true {
		return true
	}
	return false
}

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
