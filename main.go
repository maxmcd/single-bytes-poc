package main

import (
	"github.com/bytecodealliance/wasmtime-go"
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	engine := wasmtime.NewEngine()
	store := wasmtime.NewStore(engine)

	linker := wasmtime.NewLinker(store)
	// Configure WASI imports to write stdout into a file.
	wasiConfig := wasmtime.NewWasiConfig()
	wasiConfig.InheritStderr()
	wasiConfig.InheritStdout()

	// Set the version to the same as in the WAT.
	wasi, err := wasmtime.NewWasiInstance(store, wasiConfig, "wasi_snapshot_preview1")
	check(err)

	// Link WASI
	err = linker.DefineWasi(wasi)
	check(err)

	module, err := wasmtime.NewModuleFromFile(store.Engine, "target/wasm32-wasi/debug/single-bytes-poc.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	// Run the function
	nom := instance.GetExport("_start").Func()
	_, err = nom.Call()
	check(err)
}
