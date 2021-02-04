package main

import (
	"log"
	"strings"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func main() {
	// Reads the WebAssembly module as bytes.
	bytes, _ := wasm.ReadBytes("../as/build/optimized.wasm")

	module, _ := wasm.Compile(bytes)

	// Read the WASI version required for this module
	wasiVersion := wasm.WasiGetVersion(module)

	log.Println(wasiVersion)
	if wasiVersion <= 0 {
		return
	}
	importObject := wasm.NewDefaultWasiImportObjectForVersion(wasiVersion)
	instance, _ := module.InstantiateWithImportObject(importObject)
	add := instance.Exports["add"]
	result, _ := add(5, 37)
	log.Println(result)

	// test of allocating memory, writing an array and having the wasm function sum the elements in the array
	tstArr := []uint8{2, 4, 6}
	alloc := instance.Exports["alloc"]
	// ptr is a pointer to a block of memory we've allocated on the wasm side
	ptr, err := alloc(3)

	if err != nil {
		log.Println(err)
		return
	}
	// index into the block of memory we've reserved
	memory := instance.Memory.Data()[int(ptr.ToI32()):]

	// copy our array into that memory
	for i := 0; i < len(tstArr); i++ {
		memory[i] = tstArr[i]
	}

	sumArr := instance.Exports["array_sum"]
	result, err = sumArr(ptr, 3)

	if err != nil {
		log.Println(err)
		return
	}
	// enjoy the results
	log.Println(result)

	// test of returning a string from WASM to go
	returnStr := instance.Exports["returnString"]
	ptr, err = returnStr()

	i := 0
	memory = instance.Memory.Data()[int(ptr.ToI32()):]
	var output strings.Builder

	for {
		// on the wasm side, we null terminated the string, so we can rely on this C style :)
		if memory[i] == 0 {
			break
		}

		output.WriteByte(memory[i])
		i++
	}

	log.Println(output.String())

}
