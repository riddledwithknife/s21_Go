package ex00

import (
	"errors"
	"unsafe"
)

func getElement(arr []int, idx int) (int, error) {
	if len(arr) == 0 {
		return 0, errors.New("empty slice")
	}

	if idx < 0 || idx >= len(arr) {
		return 0, errors.New("index out of bounds")
	}

	element := *(*int)(unsafe.Pointer(uintptr(unsafe.Pointer(&arr[0])) + uintptr(idx)*unsafe.Sizeof(arr[0])))
	return element, nil
}
