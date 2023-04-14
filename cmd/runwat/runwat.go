package main

import (
	"fmt"
	"os"

	"github.com/bytecodealliance/wasmtime-go/v7"
)

func main() {
	// Almost all operations in wasmtime require a contextual `store`
	// argument to share, so create that first
	store := wasmtime.NewStore(wasmtime.NewEngine())

	// Try to load the first argument as a string
	filename := os.Args[1]
	wat, err := os.ReadFile(filename)
	check(err)

	// Compiling modules requires WebAssembly binary input, but the wasmtime
	// package also supports converting the WebAssembly text format to the
	// binary format.
	wasm, err := wasmtime.Wat2Wasm(string(wat))
	check(err)

	// Once we have our binary `wasm` we can compile that into a `*Module`
	// which represents compiled JIT code.
	module, err := wasmtime.NewModule(store.Engine, wasm)
	check(err)

	// Create shared memory
	memory, err := wasmtime.NewMemory(store, wasmtime.NewMemoryType(1, true, 1))
	check(err)

	// Our `hello.wat` file imports one item, so we create that function
	// here.
	item := wasmtime.WrapFunc(store, func(offset, length int32) {
		message := memory.UnsafeData(store)
		fmt.Println(string(message[offset:length]))
	})

	// Next up we instantiate a module which is where we link in all our
	// imports. We've got one import so we pass that in here.
	instance, err := wasmtime.NewInstance(store, module, []wasmtime.AsExtern{item, memory})
	check(err)

	// After we've instantiated we can lookup our `run` function and call
	// it.
	run := instance.GetFunc(store, "run")
	if run == nil {
		panic("not a function")
	}
	_, err = run.Call(store)
	check(err)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
