//go:build js && wasm
// +build js,wasm

package main

import (
	"log"
	"syscall/js"

)

func main() {
	log.Println("Hello, WebAssembly!")
	js.Global().Set("wasmGenerateWallet", js.FuncOf(func(this js.Value, p []js.Value) interface{} {
		go func() {
			say()
		}()
		return js.Undefined()
	}));
	select {};
}

func say() {
	log.Println("Generated wallet")
}
