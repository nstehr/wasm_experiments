package main

import (
	"fmt"
	"log"

	wasm "github.com/wasmerio/go-ext-wasm/wasmer"
)

func main() {
	// Reads the WebAssembly module as bytes.
	bytes, _ := wasm.ReadBytes("/Users/nstehr/wasm_test/as/build/untouched.wasm")

	// Instantiates the WebAssembly module.
	instance, _ := wasm.NewInstance(bytes)
	defer instance.Close()

	// instance 'Exports' map is empty, why?
	log.Println(instance)

	// Gets the `add` exported function from the WebAssembly instance.
	add := instance.Exports["add"]

	// Calls that exported function with Go standard values. The WebAssembly
	// types are inferred and values are casted automatically.
	result, _ := add(5, 37)

	fmt.Println(result) // 42!
}
