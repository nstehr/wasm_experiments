package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

// This function is imported from host
//export stringFromHost
func StringFromHost() *byte

func main() {
	// return_string()
	// log.Println(buffer)
}

//go:export add
func add(x, y int32) int32 {
	return x + y
}

//go:export array_sum
func array_sum(bufPtr *byte, len int) uint8 {
	tst := rawBytePtrToByteSlice(bufPtr, len)
	result := uint8(0)
	for i := 0; i < len; i++ {
		result += tst[i]
	}
	return result
}

//go:export returnString
func returnString() *uint8 {

	// sticking this here to test calling a host provided function
	// that modifies memory
	ptr := StringFromHost()
	val := rawBytePtrToByteSlice(ptr, 40)
	fmt.Println(string(val))
	test := "This is a test of the emergency broadcast system"
	arr := make([]byte, len(test)+1)
	for i := 0; i < len(test); i++ {
		arr[i] = test[i]
	}
	return &arr[0]
}

//go:export alloc
func alloc(length int) *uint8 {
	buf := make([]byte, length)
	return &buf[0]
}

// https://github.com/tetratelabs/proxy-wasm-go-sdk/blob/d91a09b6807cd8a9726a9653478428067e413d8e/proxywasm/hostcall_utils_go.go#L36
func rawBytePtrToByteSlice(raw *byte, size int) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(raw)),
		Len:  uintptr(size),
		Cap:  uintptr(size),
	}))
}
