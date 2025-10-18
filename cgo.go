//go:build windows || darwin || linux
// +build windows darwin linux

//go:generate go run check_native_libs.go

package extractous

/*
// Platform-Specific CGO Directives

// Linux
#cgo linux,amd64 CFLAGS: -I${SRCDIR}/native/linux_amd64/include
#cgo linux,amd64 LDFLAGS: -L${SRCDIR}/native/linux_amd64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,${SRCDIR}/native/linux_amd64/lib
#cgo linux,arm64 CFLAGS: -I${SRCDIR}/native/linux_arm64/include
#cgo linux,arm64 LDFLAGS: -L${SRCDIR}/native/linux_arm64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,${SRCDIR}/native/linux_arm64/lib

// macOS
#cgo darwin,amd64 CFLAGS: -I${SRCDIR}/native/darwin_amd64/include
#cgo darwin,amd64 LDFLAGS: -L${SRCDIR}/native/darwin_amd64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,${SRCDIR}/native/darwin_amd64/lib
#cgo darwin,arm64 CFLAGS: -I${SRCDIR}/native/darwin_arm64/include
#cgo darwin,arm64 LDFLAGS: -L${SRCDIR}/native/darwin_arm64/lib -lextractous_ffi -ldl -lm -lpthread -Wl,-rpath,${SRCDIR}/native/darwin_arm64/lib

// Windows
#cgo windows,amd64 CFLAGS: -I${SRCDIR}/native/windows_amd64/include
#cgo windows,amd64 LDFLAGS: -L${SRCDIR}/native/windows_amd64/lib -lextractous_ffi

// Platform-Specific Headers and Constructor Implementation
#include <stdlib.h>
#include <string.h>
#include <stdio.h>

// Windows-specific implementation
#if defined(_WIN32) || defined(_WIN64)
    #include <windows.h>

    // This runs BEFORE Go runtime initializes on Windows
    static void __attribute__((constructor)) setup_dll_search_path(void) {
        char dll_dir[MAX_PATH] = {0};
        int path_found = 0;

        // Priority 1: Check EXTRACTOUS_LIB_PATH environment variable
        const char* custom_path = getenv("EXTRACTOUS_LIB_PATH");
        if (custom_path && custom_path[0] != '\0') {
            strncpy(dll_dir, custom_path, MAX_PATH - 1);
            dll_dir[MAX_PATH - 1] = '\0';  // Ensure null termination
            path_found = 1;

            #ifdef DEBUG_EXTRACTOUS
            fprintf(stderr, "[extractous] Using EXTRACTOUS_LIB_PATH: %s\n", dll_dir);
            #endif
        }

        // Priority 2: Try to find relative to the current module
        if (!path_found) {
            HMODULE hModule = NULL;
            if (GetModuleHandleExA(
                GET_MODULE_HANDLE_EX_FLAG_FROM_ADDRESS | GET_MODULE_HANDLE_EX_FLAG_UNCHANGED_REFCOUNT,
                (LPCSTR)&setup_dll_search_path,
                &hModule
            ) && hModule) {
                char module_path[MAX_PATH];
                if (GetModuleFileNameA(hModule, module_path, MAX_PATH) > 0) {
                    // Strip filename to get directory
                    char* last_slash = strrchr(module_path, '\\');
                    if (last_slash) {
                        *last_slash = '\0';
                    }
                    // Append native library path
                    snprintf(dll_dir, sizeof(dll_dir), "%s\\native\\windows_amd64\\lib", module_path);
                    path_found = 1;

                    #ifdef DEBUG_EXTRACTOUS
                    fprintf(stderr, "[extractous] Using relative path: %s\n", dll_dir);
                    #endif
                }
            }
        }

        // Set the DLL search directory if we found a path
        if (path_found && dll_dir[0] != '\0') {
            // SetDllDirectory adds the directory to search path
            if (!SetDllDirectoryA(dll_dir)) {
                DWORD error = GetLastError();
                fprintf(stderr, "[extractous] Warning: SetDllDirectory failed for: %s (error: %lu)\n", dll_dir, error);
            } else {
                #ifdef DEBUG_EXTRACTOUS
                fprintf(stderr, "[extractous] Successfully set DLL directory\n");
                #endif
            }

            // Also use AddDllDirectory for Windows 7+ (more secure)
            HMODULE kernel32 = GetModuleHandleA("kernel32.dll");
            if (kernel32) {
                typedef DLL_DIRECTORY_COOKIE (WINAPI *AddDllDirectoryFunc)(PCWSTR);
                AddDllDirectoryFunc addDllDir =
                    (AddDllDirectoryFunc)GetProcAddress(kernel32, "AddDllDirectory");

                if (addDllDir) {
                    wchar_t wide_path[MAX_PATH];
                    if (MultiByteToWideChar(CP_UTF8, 0, dll_dir, -1, wide_path, MAX_PATH) > 0) {
                        DLL_DIRECTORY_COOKIE cookie = addDllDir(wide_path);
                        if (cookie == NULL) {
                            fprintf(stderr, "[extractous] Warning: AddDllDirectory failed\n");
                        }
                    }
                }
            }
        } else if (!path_found) {
            fprintf(stderr, "[extractous] Warning: Could not determine library path. Relying on system PATH.\n");
        }
    }

// Unix-specific implementation (Linux and macOS)
#else
    #include <dlfcn.h>

    // This runs BEFORE Go runtime initializes on Unix
    static void __attribute__((constructor)) setup_library_path(void) {
        const char* custom_path = getenv("EXTRACTOUS_LIB_PATH");

        // Only intervene if EXTRACTOUS_LIB_PATH is set
        // Otherwise, RPATH from LDFLAGS handles library loading automatically
        if (custom_path && custom_path[0] != '\0') {
            char lib_path[4096];

            // Construct full library path based on platform
            #if defined(__APPLE__) || defined(__MACH__)
                snprintf(lib_path, sizeof(lib_path), "%s/libextractous_ffi.dylib", custom_path);
            #else
                snprintf(lib_path, sizeof(lib_path), "%s/libextractous_ffi.so", custom_path);
            #endif

            #ifdef DEBUG_EXTRACTOUS
            fprintf(stderr, "[extractous] Attempting to preload from: %s\n", lib_path);
            #endif

            // Pre-load the library with RTLD_NOW | RTLD_GLOBAL
            // RTLD_NOW: resolve all symbols immediately
            // RTLD_GLOBAL: make symbols available for subsequently loaded libraries
            void* handle = dlopen(lib_path, RTLD_NOW | RTLD_GLOBAL);
            if (!handle) {
                fprintf(stderr, "[extractous] Error: Failed to preload library from EXTRACTOUS_LIB_PATH: %s\n", dlerror());
                fprintf(stderr, "[extractous] The library will not be available. Please check the path.\n");
                // Don't exit - let the linker try RPATH as fallback
            } else {
                #ifdef DEBUG_EXTRACTOUS
                fprintf(stderr, "[extractous] Successfully preloaded library from custom path\n");
                #endif
            }
        }
        // If EXTRACTOUS_LIB_PATH is not set, RPATH handles everything automatically
    }
#endif

// Include the generated header
#include "extractous.h"
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
