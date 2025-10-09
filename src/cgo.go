//go:build windows || darwin || linux
// +build windows darwin linux

package src

/*
#cgo CFLAGS: -I${SRCDIR}/../include

// Linux
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/../native/linux_amd64 -lextractous_ffi -Wl,-rpath,${SRCDIR}/../native/linux_amd64
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/../native/linux_arm64 -lextractous_ffi -Wl,-rpath,${SRCDIR}/../native/linux_arm64

// macOS
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/../native/darwin_amd64 -lextractous_ffi -Wl,-rpath,${SRCDIR}/../native/darwin_amd64
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/../native/darwin_arm64 -lextractous_ffi -Wl,-rpath,${SRCDIR}/../native/darwin_arm64

// Windows (RPATH not used)
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/../native/windows_amd64 -lextractous_ffi

#include <stdlib.h>
#include <extractous.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// init sets up library paths for runtime linking
func init() {
	// Ensure native libraries are in the search path
	setupLibraryPaths()
}

// setupLibraryPaths configures OS-specific library search paths
func setupLibraryPaths() {
	// Platform-specific library loading is handled by RPATH on Unix
	// and by PATH on Windows (see internal/loader package)
	runtime.LockOSThread() // Lock to ensure JNI thread attachment works
}

// Helper functions for C interop

func cString(s string) *C.char {
	return C.CString(s)
}

func goString(cs *C.char) string {
	return C.GoString(cs)
}

func freeString(cs *C.char) {
	C.free(unsafe.Pointer(cs))
}
