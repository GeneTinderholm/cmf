package arena

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"unsafe"
)

func TestAlloc(t *testing.T) {
	t.Run("mutating the returned pointers should mutate the backing memory", func(t *testing.T) {
		a := New()
		intPtr := Alloc[int64](a)
		*intPtr = 0x1234
		assert.Equal(t, int64(0x1234), *(*int64)(a.currentRef))
		intPtr2 := Alloc[int64](a)
		*intPtr2 = 0x5678
		assert.Equal(t, int64(0x5678), *(*int64)(unsafe.Pointer(uintptr(a.currentRef) + unsafe.Sizeof(int64(1)))))
		intSlice := AllocSlice[int64](a, 2)
		intSlice[0] = 0x1111
		assert.Equal(t, int64(0x1111), *(*int64)(
			unsafe.Pointer(
				uintptr(a.currentRef) +
					unsafe.Sizeof(int64(1))*2,
			),
		))

		intSlice[1] = 0x2222
		assert.Equal(t, int64(0x2222), *(*int64)(
			unsafe.Pointer(
				uintptr(a.currentRef) +
					unsafe.Sizeof(int64(1))*3,
			),
		))
	})

	t.Run("should not have an off by one in the bounds checking", func(t *testing.T) {
		a := New()
		AllocSlice[int64](a, int(pageSize)/8)
		assert.Nil(t, a.refs)
	})

	t.Run("should allocate a one off page if too big of a slice is asked for", func(t *testing.T) {
		a := New()
		prevRef := a.currentRef
		AllocSlice[int64](a, int(pageSize)/8+1)
		assert.Len(t, a.refs, 1)
		assert.Equal(t, prevRef, a.currentRef)
	})
	t.Run("should allocate a new page if old one fills up", func(t *testing.T) {
		a := New()
		prevRef := a.currentRef
		AllocSlice[int64](a, int(pageSize)/8-1)
		assert.Len(t, a.refs, 0)
		Alloc[int64](a)
		assert.Len(t, a.refs, 0)
		Alloc[int64](a)
		assert.Len(t, a.refs, 1)
		assert.NotEqual(t, prevRef, a.currentRef)
	})
}

func TestReset(t *testing.T) {
	t.Run("should mark all bytes to memory as free", func(t *testing.T) {
		a := New()
		Alloc[int64](a)
		Alloc[int64](a)
		AllocSlice[int64](a, 2)

		assert.NotEqual(t, uint(0x0), *(*uint)(a.currentRef))
		Reset(a)
		assert.Equal(t, uint(0x0), *(*uint)(a.currentRef))
	})
}

func BenchmarkArena(b *testing.B) {
	//BenchmarkArena/slice_comparison-10         	 3582571	       334.4 ns/op
	b.Run("slice comparison", func(b *testing.B) {
		a := make([]int, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				a[j] = j
			}
		}
	})
	//BenchmarkArena/map_comparison-10           	  177844	      6781 ns/op
	b.Run("map comparison", func(b *testing.B) {
		a := map[int]int{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				a[j] = j
			}
		}
	})
	//BenchmarkArena/heap_comparison-10          	  136845	      8747 ns/op
	b.Run("heap comparison", func(b *testing.B) {
		a := make([]*int, 1024)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				j := j
				a[j] = &j
			}
		}
	})

	//BenchmarkArena/arena_comparison-10         	  356715	      3352 ns/op
	b.Run("arena comparison", func(b *testing.B) {
		a := make([]*int, 1024)
		ar := New()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			for j := 0; j < 1024; j++ {
				intPtr := Alloc[int](ar)
				*intPtr = j
				a[j] = intPtr
			}
		}
	})
}
