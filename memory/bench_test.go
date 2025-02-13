package memory

import (
	"reflect"
	"testing"
	"unsafe"
)

func BenchmarkSize(b *testing.B) {
	var pointStr int64
	pointStr = 1
	ref := reflect.Indirect(reflect.ValueOf(pointStr))

	b.Run("unsafe.Sizeof", func(b *testing.B) {
		var sizeof uint64
		for i := 0; i < b.N; i++ {
			sizeof = uint64(unsafe.Sizeof(pointStr))
		}
		b.Log("unsafe.Sizeof:", sizeof)
	})
	b.Run("unsafe.Pointer.SizeBack", func(b *testing.B) {
		var sizeof uint64
		for i := 0; i < b.N; i++ {
			sizeof = uint64(ref.Type().Size())
		}
		b.Log("unsafe.Pointer.SizeBack:", sizeof)
	})
}
