package extractous

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/native/linux_amd64 -lextractous_ffi -Wl,-rpath,'$$ORIGIN'
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/native/linux_arm64 -lextractous_ffi -Wl,-rpath,'$$ORIGIN'
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/native/darwin_amd64 -lextractous_ffi -Wl,-rpath,@loader_path
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/native/darwin_arm64 -lextractous_ffi -Wl,-rpath,@loader_path
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/native/windows_amd64 -lextractous_ffi

#include "extractous.h"
#include <stdlib.h>
#include <stdbool.h>
*/
import "C"
