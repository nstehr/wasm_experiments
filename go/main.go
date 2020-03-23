package main

import (
	"fmt"
	"log"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func main() {
	// Reads the WebAssembly module as bytes.
	bytes, _ := wasm.ReadBytes("../as/build/untouched.wasm")

	module, _ := wasm.Compile(bytes)

	// Read the WASI version required for this module
	wasiVersion := wasm.WasiGetVersion(module)

	log.Println(wasiVersion)
	importObject := wasm.NewDefaultWasiImportObjectForVersion(wasiVersion)
	instance, _ := module.InstantiateWithImportObject(importObject)
	log.Println(instance)
	add := instance.Exports["add"]
	log.Println(add)
	result, _ := add(5, 37)
	fmt.Println(result)
}
