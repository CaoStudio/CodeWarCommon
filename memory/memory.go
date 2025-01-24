package memory

import (
	"reflect"
	"unsafe"
)

// Size returns the size of 'v' in bytes.
// If there is an error during calculation, Of returns -1.
func Size(v interface{}) uint64 {
	// Cache with every visited pointer, so we don't count two pointers
	// to the same memory twice.
	cache := make(map[unsafe.Pointer]bool)
	return sizeOf(reflect.Indirect(reflect.ValueOf(v)), cache)
}

// sizeOf returns the number of bytes the actual data represented by v occupies in memory.
// If there is an error, sizeOf returns -1.
func sizeOf(v reflect.Value, cache map[unsafe.Pointer]bool) uint64 {
	switch v.Kind() {

	case reflect.Array:
		var sum uint64 = 0
		for i := 0; i < v.Len(); i++ {
			s := sizeOf(v.Index(i), cache)
			if s < 0 {
				return 0
			}
			sum += s
		}

		return sum + uint64(v.Cap()-v.Len())*uint64(v.Type().Elem().Size())

	case reflect.Slice:
		// return 0 if this node has been visited already
		if cache[v.UnsafePointer()] {
			return 0
		}
		cache[v.UnsafePointer()] = true

		var sum uint64 = 0
		for i := 0; i < v.Len(); i++ {
			s := sizeOf(v.Index(i), cache)
			if s < 0 {
				return 0
			}
			sum += s
		}

		sum += uint64(v.Cap()-v.Len()) * uint64(v.Type().Elem().Size())

		return sum + uint64(v.Type().Size())

	case reflect.Struct:
		var sum uint64 = 0
		for i, n := 0, v.NumField(); i < n; i++ {
			s := sizeOf(v.Field(i), cache)
			if s < 0 {
				return 0
			}
			sum += s
		}

		// Look for struct padding.
		padding := uint64(v.Type().Size())
		for i, n := 0, v.NumField(); i < n; i++ {
			padding -= uint64(v.Field(i).Type().Size())
		}

		return sum + padding

	case reflect.String:
		s := v.String()
		//hdr := (*reflect.StringHeader)(unsafe.Pointer(&s))
		hdr := v.UnsafePointer()
		//unsafe.StringData(s)
		if cache[hdr] {
			return uint64(v.Type().Size())
		}
		cache[hdr] = true
		return uint64(len(s)) + uint64(v.Type().Size())

	case reflect.Ptr:
		// return Ptr size if this node has been visited already (infinite recursion)
		if cache[v.UnsafePointer()] {
			return uint64(v.Type().Size())
		}
		cache[v.UnsafePointer()] = true
		if v.IsNil() {
			return uint64(reflect.New(v.Type()).Type().Size())
		}
		s := sizeOf(reflect.Indirect(v), cache)
		if s < 0 {
			return 0
		}
		return s + uint64(v.Type().Size())

	case reflect.Bool,
		reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Int, reflect.Uint,
		reflect.Chan,
		reflect.Uintptr,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.Func:
		return uint64(v.Type().Size())

	case reflect.Map:
		// return 0 if this node has been visited already (infinite recursion)
		if cache[v.UnsafePointer()] {
			return 0
		}
		cache[v.UnsafePointer()] = true
		var sum uint64 = 0
		keys := v.MapKeys()
		for i := range keys {
			val := v.MapIndex(keys[i])
			// calculate size of key and value separately
			sv := sizeOf(val, cache)
			if sv < 0 {
				return 0
			}
			sum += sv
			sk := sizeOf(keys[i], cache)
			if sk < 0 {
				return 0
			}
			sum += sk
		}
		// Include overhead due to unused map buckets.  10.79 comes
		// from https://golang.org/src/runtime/map.go.
		return sum + uint64(v.Type().Size()) + uint64(float64(len(keys))*10.79)

	case reflect.Interface:
		return sizeOf(v.Elem(), cache) + uint64(v.Type().Size())

	}

	return 0
}
