package main

// import (
// 	"io/ioutil"
// 	"testing"
//
// 	wasm3 "github.com/matiasinsaurralde/go-wasm3"
// )
//
// var (
// 	wasmBytes []byte
// )
//
// func init() {
// 	var err error
// 	wasmBytes, err = ioutil.ReadFile(wasmFilename)
// 	if err != nil {
// 		panic(err)
// 	}
// }
//
// func TestSum(t *testing.T) {
// 	runtime := wasm3.NewRuntime(&wasm3.Config{
// 		Environment: wasm3.NewEnvironment(),
// 		StackSize:   64 * 1024,
// 	})
// 	defer runtime.Destroy()
// 	_, err := runtime.Load(wasmBytes)
// 	if err != nil {
// 		t.Fatal(err)
// 	}
// 	{
// 		fn, err := runtime.FindFunction(fnName)
// 		if err != nil {
// 			t.Fatal(err)
// 		}
// 		result, _ := fn(1, 1)
// 		if result.(int32) != 2 {
// 			t.Fatal("Result doesn't match")
// 		}
// 	}
// 	{
// 		// fn, err := module.GetFunctionByName(fnName)
// 		// if err != nil {
// 		// 	t.Fatal(err)
// 		// }
// 	}
// }
