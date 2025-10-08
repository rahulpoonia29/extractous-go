package extractous

import "unsafe"
import "C"

// Metadata represents document metadata as key-value pairs
type Metadata map[string][]string

func cMetadataToGo(cMeta *C.CMetadata) Metadata {
	if cMeta == nil {
		return nil
	}

	meta := make(Metadata)

	for i := 0; i < int(cMeta.len); i++ {
		keyPtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cMeta.keys)) +
			uintptr(i)*unsafe.Sizeof(*cMeta.keys)))
		valuePtr := (**C.char)(unsafe.Pointer(uintptr(unsafe.Pointer(cMeta.values)) +
			uintptr(i)*unsafe.Sizeof(*cMeta.values)))

		key := C.GoString(*keyPtr)
		value := C.GoString(*valuePtr)

		// Values are comma-separated in C, split them back
		// TODO: Handle proper CSV parsing if values contain commas
		meta[key] = []string{value}
	}

	return meta
}
