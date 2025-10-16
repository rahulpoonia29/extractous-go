//go:build windows || darwin || linux
// +build windows darwin linux

//go:generate go run check_native_libs.go

package extractous

/*
// Linux
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/native/linux_amd64/include
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/native/linux_amd64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,'$ORIGIN/../native/linux_amd64/lib'
#cgo linux,arm64 CFLAGS: -I${SRCDIR}/native/linux_arm64/include
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/native/linux_arm64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,'$ORIGIN/../native/linux_arm64/lib'

// macOS
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/native/darwin_amd64/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/native/darwin_amd64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,'@executable_path/../native/darwin_amd64/lib'
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/native/darwin_arm64/include
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/native/darwin_arm64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,'@executable_path/../native/darwin_arm64/lib'

// Windows
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/native/windows_amd64/include
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/native/windows_amd64/lib -lextractous_ffi

#include <stdlib.h>
#include <extractous.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// init locks the OS thread for JNI compatibility.
// The native libraries are expected to be in the ./native directory.
func init() {
	runtime.LockOSThread()
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
