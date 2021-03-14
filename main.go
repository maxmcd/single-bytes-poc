package main

import (
	"encoding/binary"
	"fmt"

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

	_ = wasi
	err = linker.DefineWasi(wasi)
	check(err)

	fmt.Println()
	linker.AllowShadowing(true)
	var memory *wasmtime.Memory
	err = linker.DefineFunc("wasi_snapshot_preview1", "fd_write", func(a int32, b int32, c int32, d int32) int32 {
		fmt.Println("called wasi_snapshot_preview1.fd_write()")
		buf := memory.UnsafeData()
		idx := binary.LittleEndian.Uint32(buf[b : b+4])
		len := binary.LittleEndian.Uint32(buf[b+4 : b+8])
		fmt.Println(string(buf[idx : idx+len]))
		binary.LittleEndian.PutUint32(buf[d:d+4], len)
		return 0
	})
	check(err)

	module, err := wasmtime.NewModuleFromFile(store.Engine, "target/wasm32-wasi/debug/single-bytes-poc.wasm")
	check(err)
	instance, err := linker.Instantiate(module)
	check(err)

	memory = instance.GetExport("memory").Memory()

	// Run the function
	nom := instance.GetExport("_start").Func()
	_, err = nom.Call()
	check(err)
}
