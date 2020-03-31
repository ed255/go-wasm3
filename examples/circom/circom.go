package main

/*
#cgo CFLAGS: -Iinclude

#include <stdio.h>

#include "wasm3.h"
#include "m3_env.h"
#include "m3_api_defs.h"
#include "m3_api_libc.h"
#include "m3_api_wasi.h"
#include "extra/wasi_core.h"

m3ApiRawFunction(m3_wasm3_raw_error)
{
    m3ApiGetArg      (int32_t, code)
    m3ApiGetArg      (int32_t, pStr)
    m3ApiGetArg      (int32_t, param1)
    m3ApiGetArg      (int32_t, param2)
    m3ApiGetArg      (int32_t, param3)
    m3ApiGetArg      (int32_t, param4)

    return NULL;
}

m3ApiRawFunction(m3_wasm3_raw_logSetSignal)
{
    m3ApiGetArg      (int32_t, signal)
    m3ApiGetArg      (int32_t, val)

    return NULL;
}

m3ApiRawFunction(m3_wasm3_raw_logGetSignal)
{
    m3ApiGetArg      (int32_t, signal)
    m3ApiGetArg      (int32_t, val)

    return NULL;
}

m3ApiRawFunction(m3_wasm3_raw_logFinishComponent)
{
    m3ApiGetArg      (int32_t, cIdx)

    return NULL;
}

m3ApiRawFunction(m3_wasm3_raw_logStartComponent)
{
    m3ApiGetArg      (int32_t, cIdx)

    return NULL;
}

m3ApiRawFunction(m3_wasm3_raw_log)
{
    m3ApiGetArg      (int32_t, code)

    return NULL;
}

*/
import "C"

import (
	"fmt"
	"io/ioutil"
	"log"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename = "mycircuit.wasm"
)

type WitnessCalcFns struct {
	getFrLen          func() (int32, error)
	getPRawPrime      func() (int32, error)
	getNVars          func() (int32, error)
	init              func(sanityCheck int32) error
	getSignalOffset32 func(pR, component, hashMSB, hashLSB int32) error
	setSignal         func(cIdx, component, signal, pVal int32) error
	getPWitness       func(w int32) (int32, error)
	getWitnessBuffer  func() (int32, error)
}

func NewWitnessCalcFns(r *wasm3.Runtime) (*WitnessCalcFns, error) {
	_getFrLen, err := r.FindFunction("getFrLen")
	if err != nil {
		return nil, err
	}
	getFrLen := func() (int32, error) {
		res, err := _getFrLen()
		if err != nil {
			return 0, err
		}
		return res.(int32), nil
	}
	_getPRawPrime, err := r.FindFunction("getPRawPrime")
	if err != nil {
		return nil, err
	}
	getPRawPrime := func() (int32, error) {
		res, err := _getPRawPrime()
		if err != nil {
			return 0, err
		}
		return res.(int32), nil
	}
	_getNVars, err := r.FindFunction("getNVars")
	if err != nil {
		return nil, err
	}
	getNVars := func() (int32, error) {
		res, err := _getNVars()
		if err != nil {
			return 0, err
		}
		return res.(int32), nil
	}
	_init, err := r.FindFunction("init")
	if err != nil {
		return nil, err
	}
	init := func(sanityCheck int32) error {
		_, err := _init(sanityCheck)
		if err != nil {
			return err
		}
		return nil
	}
	_getSignalOffset32, err := r.FindFunction("getSignalOffset32")
	if err != nil {
		return nil, err
	}
	getSignalOffset32 := func(pR, component, hashMSB, hashLSB int32) error {
		_, err := _getSignalOffset32(pR, component, hashMSB, hashLSB)
		if err != nil {
			return err
		}
		return nil
	}
	// _setSignal, err := r.FindFunctionByModule(".unnamed", "setSignal")
	_setSignal, err := r.FindFunction("setSignal")
	if err != nil {
		println("B")
		return nil, err
	}
	setSignal := func(cIdx, component, signal, pVal int32) error {
		_, err := _setSignal(cIdx, component, signal, pVal)
		if err != nil {
			return err
		}
		return nil
	}
	_getPWitness, err := r.FindFunction("getPWitness")
	if err != nil {
		return nil, err
	}
	getPWitness := func(w int32) (int32, error) {
		res, err := _getPWitness(w)
		if err != nil {
			return 0, err
		}
		return res.(int32), nil
	}
	_getWitnessBuffer, err := r.FindFunction("getWitnessBuffer")
	if err != nil {
		return nil, err
	}
	getWitnessBuffer := func() (int32, error) {
		res, err := _getWitnessBuffer()
		if err != nil {
			return 0, err
		}
		return res.(int32), nil
	}

	return &WitnessCalcFns{
		getFrLen:          getFrLen,
		getPRawPrime:      getPRawPrime,
		getNVars:          getNVars,
		init:              init,
		getSignalOffset32: getSignalOffset32,
		setSignal:         setSignal,
		getPWitness:       getPWitness,
		getWitnessBuffer:  getWitnessBuffer,
	}, nil
}

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

	fmt.Printf("NumImports: %v\n", module.NumImports())

	// return

	// for _, funName := range module.FunctionNames() {
	// 	fmt.Println(funName)
	// }

	err = module.LinkRawFunction("runtime", "error", "v(iiiiii)", C.m3_wasm3_raw_error)
	if err != nil {
		panic(err)
	}
	err = module.LinkRawFunction("runtime", "logSetSignal", "v(ii)", C.m3_wasm3_raw_logSetSignal)
	if err != nil {
		panic(err)
	}
	err = module.LinkRawFunction("runtime", "logGetSignal", "v(ii)", C.m3_wasm3_raw_logGetSignal)
	if err != nil {
		panic(err)
	}
	err = module.LinkRawFunction("runtime", "logFinishComponent", "v(i)", C.m3_wasm3_raw_logFinishComponent)
	if err != nil {
		panic(err)
	}
	err = module.LinkRawFunction("runtime", "logStartComponent", "v(i)", C.m3_wasm3_raw_logStartComponent)
	if err != nil {
		panic(err)
	}
	err = module.LinkRawFunction("runtime", "log", "v(i)", C.m3_wasm3_raw_log)
	if err != nil {
		panic(err)
	}

	fns, err := NewWitnessCalcFns(runtime)
	if err != nil {
		panic(err)
	}

	// fnGetNVars := "getNVars"
	// fn, err := runtime.FindFunction(fnGetNVars)
	// if err != nil {
	// 	panic(err)
	// }
	// log.Printf("Found '%s' function (using runtime.FindFunction)", fnGetNVars)
	// result, _ := fn()
	// log.Print("Result is: ", result)

	err = fns.init(1)
	if err != nil {
		panic(err)
	}
}
