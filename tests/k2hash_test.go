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

package k2hashtest

import (
	"fmt"
	"os"
	"runtime"
	"testing"
	"unsafe"
)

// TestMain is the main function of testing.
func TestMain(m *testing.M) {
	if runtime.GOOS != "linux" && runtime.GOARCH != "amd64" {
		fmt.Fprintf(os.Stderr, "k2hash currently works on linux only")
		os.Exit(-1)
	}
	i := uint32(1)
	b := (*[4]byte)(unsafe.Pointer(&i))
	if b[0] != 1 {
		fmt.Fprintf(os.Stderr, "k2hash_go currently works on little endian alignment only")
		os.Exit(-1)
	}
	// rpm or deb install the k2hash.so in /usr/lib or /usr/lib64
	if _, err := os.Stat("/usr/lib/libk2hash.so"); err != nil {
		if _, err := os.Stat("/usr/lib64/libk2hash.so"); err != nil {
			fmt.Fprintf(os.Stderr, "Please install the k2hash package at first")
			os.Exit(-1)
		}
	}
	setUp(m)
	status := m.Run()
	tearDown(m)
	os.Exit(status)
}

func setUp(m *testing.M) {
	//fmt.Println("test.go setUp")
}

func tearDown(m *testing.M) {
	//fmt.Println("test.go tearDown")
}
func TestBumpDebugLevel(t *testing.T)       { testBumpDebugLevel(t) }
func TestSetDebugLevelSilent(t *testing.T)  { testSetDebugLevelSilent(t) }
func TestSetDebugLevelError(t *testing.T)   { testSetDebugLevelError(t) }
func TestSetDebugLevelWarning(t *testing.T) { testSetDebugLevelWarning(t) }
func TestSetDebugLevelMessage(t *testing.T) { testSetDebugLevelMessage(t) }
func TestSetDebugFile(t *testing.T)         { testSetDebugFile(t) }
func TestUnsetDebugFile(t *testing.T)       { testUnsetDebugFile(t) }
func TestLoadDebugEnv(t *testing.T)         { testLoadDebugEnv(t) }
func TestSetSignalUser1(t *testing.T)       { testSetSignalUser1(t) }
func TestDumpHead(t *testing.T)             { testDumpHead(t) }
func TestDumpKeyTable(t *testing.T)         { testDumpKeyTable(t) }
func TestDumpElementTable(t *testing.T)     { testDumpElementTable(t) }
func TestDumpFull(t *testing.T)             { testDumpFull(t) }
func TestPrintState(t *testing.T)           { testPrintState(t) }
func TestPrintVersion(t *testing.T)         { testPrintVersion(t) }
func TestLoadHashLibrary(t *testing.T)      { testLoadHashLibrary(t) }
func TestUnloadHashLibrary(t *testing.T)    { testUnloadHashLibrary(t) }
func TestLoadTxLibrary(t *testing.T)        { testLoadTxLibrary(t) }
func TestUnloadTxLibrary(t *testing.T)      { testUnloadTxLibrary(t) }

func TestSet(t *testing.T)       { testSet(t) }
func TestGet(t *testing.T)       { testGet(t) }
func TestRemove(t *testing.T)    { testRemove(t) }
func TestAddSubKey(t *testing.T) { testAddSubKey(t) }

func TestEnableMtime(t *testing.T)           { testEnableMtime(t) }
func TestEnableEncryption(t *testing.T)      { testEnableEncryption(t) }
func TestEnableHistory(t *testing.T)         { testEnableHistory(t) }
func TestSetExpirationDuration(t *testing.T) { testSetExpirationDuration(t) }

func TestSetDefaultEncryptionPassword(t *testing.T) { testSetDefaultEncryptionPassword(t) }
func TestPrintAttrVersion(t *testing.T)             { testPrintAttrVersion(t) }
func TestPrintAttrInformation(t *testing.T)         { testPrintAttrInformation(t) }

func TestBeginTx(t *testing.T)               { testBeginTx(t) }
func TestStopTx(t *testing.T)                { testStopTx(t) }
func TestGetTxFileFD(t *testing.T)           { testGetTxFileFD(t) }
func TestGetTxThreadPoolSize(t *testing.T)   { testGetTxThreadPoolSize(t) }
func TestSetTxThreadPoolSize(t *testing.T)   { testSetTxThreadPoolSize(t) }
func TestUnsetTxThreadPoolSize(t *testing.T) { testUnsetTxThreadPoolSize(t) }
func TestLoadFromFile(t *testing.T)          { testLoadFromFile(t) }
func TestDumpToFile(t *testing.T)            { testDumpToFile(t) }

// Local Variables:
// c-basic-offset: 4
// tab-width: 4
// indent-tabs-mode: t
// End:
// vim600: noexpandtab sw=4 ts=4 fdm=marker
// vim<600: noexpandtab sw=4 ts=4
