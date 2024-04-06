package internal

import "unsafe"

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
