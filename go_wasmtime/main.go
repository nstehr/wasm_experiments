package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"

	"github.com/bytecodealliance/wasmtime-go"
)

func main() {

	type Nic struct {
		Mac string
		IP  string
	}
	type Device struct {
		Id   string
		Nics []Nic
	}

	var mem *wasmtime.Memory
	var alloc *wasmtime.Func

	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(wasmtime.NewEngine())
	linker := wasmtime.NewLinker(store)

	//bytes, _ := ioutil.ReadFile("../as/build/optimized.wasm")
	bytes, _ := ioutil.ReadFile("../tinygo/wasm.wasm")
	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModule(store.Engine, bytes)
	check(err)

	strFromGoFn := func() int32 {
		strFromGo := []byte("This is a string passing from go to WASM")

		ptr, err := alloc.Call(len(strFromGo))
		check(err)

		memory := mem.UnsafeData()[int(ptr.(int32)):]

		// copy our array into that memory
		for i := 0; i < len(strFromGo); i++ {
			memory[i] = strFromGo[i]
		}
		return ptr.(int32)
	}

	wasiConfig := wasmtime.NewWasiConfig()
	// these will allow our WASM code to access the host stdin, stdout and stderr
	wasiConfig.InheritStdin()
	wasiConfig.InheritStdout()
	wasiConfig.InheritStderr()

	// assemblyscript generates the wasi interface wasi_snapshot_preview1
	//wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	// tinygo generates the wasi interface with `wasi_unstable`
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_unstable")
	check(err)
	err = linker.DefineWasi(wasi)
	check(err)
	// assemblyscript (atleast my sample code) generates the import against `index`
	linker.DefineFunc("index", "stringFromHost", strFromGoFn)
	linker.DefineFunc("env", "stringFromHost", strFromGoFn)

	instance, err := linker.Instantiate(module)
	check(err)

	mem = instance.GetExport("memory").Memory()

	nom := instance.GetExport("add").Func()
	result, err := nom.Call(1, 2)
	check(err)

	log.Println(result)

	tstArr := []uint8{9, 9, 9}
	alloc = instance.GetExport("alloc").Func()
	// ptr is a pointer to a block of memory we've allocated on the wasm side
	ptr, err := alloc.Call(3)
	check(err)
	if err != nil {
		log.Println(err)
		return
	}

	memory := mem.UnsafeData()[int(ptr.(int32)):]

	// copy our array into that memory
	for i := 0; i < len(tstArr); i++ {
		memory[i] = tstArr[i]
	}

	sumArr := instance.GetExport("array_sum").Func()
	result, err = sumArr.Call(ptr, 3)
	check(err)
	// enjoy the results
	log.Println(result)

	d := Device{}
	n1 := Nic{Mac: "abs:123", IP: "1.2.2.3"}
	n2 := Nic{Mac: "abs:444", IP: "2.2.2.3"}
	d.Nics = []Nic{n1, n2}
	d.Id = "Device1"

	data, err := json.Marshal(&d)
	check(err)
	ptr, err = alloc.Call(len(data))
	check(err)
	memory = mem.UnsafeData()[int(ptr.(int32)):]

	// copy our array into that memory
	for i := 0; i < len(data); i++ {
		memory[i] = data[i]
	}
	execute := instance.GetExport("execute").Func()
	execute.Call(ptr, len(data))

	returnStr := instance.GetExport("returnString").Func()
	ptr, err = returnStr.Call()

	i := 0
	memory = mem.UnsafeData()[int(ptr.(int32)):]
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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
