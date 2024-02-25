package arena

import (
	"unsafe"
)

/*
Arena is terribly unsafe, and is really only in here because I'm bored.
Should probably not be used unless you're doing an optimization,
and even then only as a convenience, if you're using this you should be able to write
it yourself.
*/
type Arena struct {
	currentSpot uintptr
	currentRef  unsafe.Pointer
	refs        []unsafe.Pointer
}

// TODO let them pass in page size?
var pageSize = uintptr(4 * 1024)

func New() *Arena {
	return &Arena{currentRef: unsafe.Pointer(unsafe.SliceData(make([]byte, pageSize)))}
}
func Alloc[T any](a *Arena) *T {
	var t T
	sizeOfT := unsafe.Sizeof(t)
	if sizeOfT > pageSize {
		ptr := new(T)
		a.refs = append(a.refs, unsafe.Pointer(ptr))
		return ptr
	}
	if sizeOfT+a.currentSpot > pageSize {
		a.refs = append(a.refs, a.currentRef)
		a.currentSpot = 0
		a.currentRef = unsafe.Pointer(unsafe.SliceData(make([]byte, pageSize)))
	}
	ptr := (*T)(unsafe.Pointer(uintptr(a.currentRef) + a.currentSpot))
	a.currentSpot += sizeOfT
	return ptr
}

func AllocSlice[T any](a *Arena, length int) []T {
	var t T
	sizeOfT := unsafe.Sizeof(t) * uintptr(length)
	if sizeOfT > pageSize {
		ptr := unsafe.Pointer(unsafe.SliceData(make([]byte, sizeOfT)))
		a.refs = append(a.refs, ptr)
		return unsafe.Slice((*T)(ptr), length)
	}
	if sizeOfT+a.currentSpot > pageSize {
		a.refs = append(a.refs, a.currentRef)
		a.currentSpot = 0
		a.currentRef = unsafe.Pointer(unsafe.SliceData(make([]byte, pageSize)))
	}
	ptr := (*T)(unsafe.Pointer(uintptr(a.currentRef) + a.currentSpot))
	a.currentSpot += sizeOfT
	return unsafe.Slice(ptr, length)
}

func Reset(a *Arena) {
	*a = Arena{currentRef: unsafe.Pointer(unsafe.SliceData(make([]byte, pageSize)))}
}
