//go:build windows || darwin || linux
// +build windows darwin linux

//go:generate go run check_native_libs.go

package extractous

/*
// Linux
// #cgo linux,amd64 CFLAGS: -I${SRCDIR}
// #cgo linux,amd64 LDFLAGS: -lextractous_ffi -ldl -lm -lpthread
// #cgo linux,arm64 CFLAGS: -I${SRCDIR}
// #cgo linux,arm64 LDFLAGS: -lextractous_ffi -ldl -lm -lpthread

// // macOS
// #cgo darwin,amd64 CFLAGS: -I${SRCDIR}
// #cgo darwin,amd64 LDFLAGS: -lextractous_ffi -ldl -lm -lpthread
// #cgo darwin,arm64 CFLAGS: -I${SRCDIR}
// #cgo darwin,arm64 LDFLAGS: -lextractous_ffi -ldl -lm -lpthread

// // Windows
// #cgo windows,amd64 CFLAGS: -I${SRCDIR}
// #cgo windows,amd64 LDFLAGS: -lextractous_ffi


// Include the generated header
// #include "extractous.h"
*/
/*
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"unsafe"
)

// init locks the OS thread for JNI compatibility and library initialization.
// The constructor functions above run BEFORE this init() is called.
func init() {
	runtime.LockOSThread()
}

// Helper Functions for C Interop
func cString(s string) *C.char {
	return C.CString(s)
}

func goString(cs *C.char) string {
	return C.GoString(cs)
}

func freeString(cs *C.char) {
	C.free(unsafe.Pointer(cs))
}
