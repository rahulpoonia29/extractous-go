package src

/*
#include <extractous.h>
#include <stdlib.h>
*/
import "C"
import (
	"runtime"
	"strings"
	"unsafe"
)

// Metadata represents document metadata as key-value pairs.
// Values are slices because some metadata fields can have multiple values.
type Metadata map[string][]string

// metadataWrapper wraps C metadata for proper cleanup
type metadataWrapper struct {
	ptr *C.struct_CMetadata
}

// newMetadata converts C metadata to Go and sets up cleanup
func newMetadata(cMeta *C.struct_CMetadata) Metadata {
	if cMeta == nil {
		return make(Metadata)
	}

	// Create wrapper for proper cleanup
	wrapper := &metadataWrapper{ptr: cMeta}
	runtime.SetFinalizer(wrapper, (*metadataWrapper).free)

	// Convert to Go map
	result := make(Metadata, int(cMeta.len))

	if cMeta.len == 0 {
		return result
	}

	// Convert C arrays to Go slices
	keys := unsafe.Slice(cMeta.keys, cMeta.len)
	values := unsafe.Slice(cMeta.values, cMeta.len)

	for i := 0; i < int(cMeta.len); i++ {
		key := C.GoString(keys[i])
		value := C.GoString(values[i])

		// Values are comma-separated in C, split them into individual values
		valueSlice := strings.Split(value, ",")
		// Trim whitespace from each value
		for j := range valueSlice {
			valueSlice[j] = strings.TrimSpace(valueSlice[j])
		}
		result[key] = valueSlice
	}

	return result
}

// free releases C metadata resources
func (m *metadataWrapper) free() {
	if m.ptr != nil {
		C.extractous_metadata_free(m.ptr)
		m.ptr = nil
	}
}

// Get returns the first value for a metadata key, or empty string if not found
func (m Metadata) Get(key string) string {
	if vals, ok := m[key]; ok && len(vals) > 0 {
		return vals[0]
	}
	return ""
}

// GetAll returns all values for a metadata key
func (m Metadata) GetAll(key string) []string {
	return m[key]
}

// Has checks if a metadata key exists
func (m Metadata) Has(key string) bool {
	_, ok := m[key]
	return ok
}

// Keys returns all metadata keys
func (m Metadata) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
