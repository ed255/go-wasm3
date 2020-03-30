package main

import (
	"io/ioutil"
	"log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename = "mycircuit.wasm"
)

func main() {
	log.Print("Initializing WASM3")

	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	log.Println("Runtime ok")
	// err := runtime.ResizeMemory(16)
	// if err != nil {
	// 	panic(err)
	// }

	// log.Println("Runtime Memory len: ", len(runtime.Memory()))

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	if err != nil {
		panic(err)
	}
	log.Printf("Read WASM module (%d bytes)\n", len(wasmBytes))

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		panic(err)
	}
	module, err = runtime.LoadModule(module)
	if err != nil {
		panic(err)
	}
	log.Print("Loaded module")

	fnGetNVars := "getNVars"
	fn, err := runtime.FindFunction(fnGetNVars)
	if err != nil {
		panic(err)
	}
	log.Printf("Found '%s' function (using runtime.FindFunction)", fnGetNVars)
	result, _ := fn()
	log.Print("Result is: ", result)
}
