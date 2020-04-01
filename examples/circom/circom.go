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
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io/ioutil"
	"log"
	"math"
	"math/big"
	"reflect"

	wasm3 "github.com/matiasinsaurralde/go-wasm3"
)

const (
	wasmFilename   = "mycircuit2.wasm"
	inputsFilename = "mycircuit2-input.json"
)

func LinkImports(m *wasm3.Module) error {
	err := m.LinkRawFunction("runtime", "error", "v(iiiiii)", C.m3_wasm3_raw_error)
	if err != nil {
		return err
	}
	err = m.LinkRawFunction("runtime", "logSetSignal", "v(ii)", C.m3_wasm3_raw_logSetSignal)
	if err != nil {
		return err
	}
	err = m.LinkRawFunction("runtime", "logGetSignal", "v(ii)", C.m3_wasm3_raw_logGetSignal)
	if err != nil {
		return err
	}
	err = m.LinkRawFunction("runtime", "logFinishComponent", "v(i)", C.m3_wasm3_raw_logFinishComponent)
	if err != nil {
		return err
	}
	err = m.LinkRawFunction("runtime", "logStartComponent", "v(i)", C.m3_wasm3_raw_logStartComponent)
	if err != nil {
		return err
	}
	err = m.LinkRawFunction("runtime", "log", "v(i)", C.m3_wasm3_raw_log)
	if err != nil {
		return err
	}
	return nil
}

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

func NewWitnessCalcFns(r *wasm3.Runtime, m *wasm3.Module) (*WitnessCalcFns, error) {
	if err := LinkImports(m); err != nil {
		return nil, err
	}
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

	// fmt.Printf("NumImports: %v\n", module.NumImports())
	fns, err := NewWitnessCalcFns(runtime, module)
	if err != nil {
		panic(err)
	}

	inputsBytes, err := ioutil.ReadFile(inputsFilename)
	if err != nil {
		panic(err)
	}
	inputs, err := parseInputs(inputsBytes)
	if err != nil {
		panic(err)
	}
	log.Print("Inputs: ", inputs)

	witnessCalculator, err := NewWitnessCalculator(runtime, fns)
	if err != nil {
		panic(err)
	}
	log.Print(witnessCalculator)
	w, err := witnessCalculator.CalculateWitness(inputs, false)
	if err != nil {
		panic(err)
	}
	log.Print("Witness: ", w)
}

type WitnessCalculator struct {
	n32    int32
	prime  *big.Int
	mask32 *big.Int
	nVars  int32
	n64    uint
	r      *big.Int
	rInv   *big.Int

	shortMax *big.Int
	shortMin *big.Int

	runtime *wasm3.Runtime
	fns     *WitnessCalcFns
}

func loadBigInt(runtime *wasm3.Runtime, p int32, n int32) *big.Int {
	bigIntBytes := make([]byte, n)
	copy(bigIntBytes, runtime.Memory()[p:p+n])
	return new(big.Int).SetBytes(swap(bigIntBytes))
}

func storeBigInt(runtime *wasm3.Runtime, p int32, v *big.Int) {
	bigIntBytes := swap(v.Bytes())
	copy(runtime.Memory()[p:p+int32(len(bigIntBytes))], bigIntBytes)
}

func NewWitnessCalculator(runtime *wasm3.Runtime, fns *WitnessCalcFns) (*WitnessCalculator, error) {
	n32, err := fns.getFrLen()
	if err != nil {
		return nil, err
	}
	// n32 = (n32 >> 2) - 2
	n32 = n32 - 8
	log.Print("n32: ", n32)

	pRawPrime, err := fns.getPRawPrime()
	if err != nil {
		return nil, err
	}
	log.Print("pRawPrime: ", pRawPrime)

	prime := loadBigInt(runtime, pRawPrime, n32)
	log.Print("prime: ", prime)

	mask32 := new(big.Int).SetUint64(0xFFFFFFFF)
	log.Print("mask32: ", mask32)
	nVars, err := fns.getNVars()
	if err != nil {
		return nil, err
	}
	log.Print("nVars: ", nVars)

	n64 := uint(((prime.BitLen() - 1) / 64) + 1)
	log.Print("n64: ", n64)
	r := new(big.Int).SetInt64(1)
	r.Lsh(r, n64*64)
	log.Print("r: ", r)
	rInv := new(big.Int).ModInverse(r, prime)
	log.Print("rInv: ", rInv)

	shortMax, ok := new(big.Int).SetString("0x80000000", 0)
	if !ok {
		return nil, fmt.Errorf("unable to set shortMax from string")
	}
	shortMin := new(big.Int).Set(prime)
	shortMin.Sub(shortMin, shortMax)

	return &WitnessCalculator{
		n32:    n32,
		prime:  prime,
		mask32: mask32,
		nVars:  nVars,
		n64:    n64,
		r:      r,
		rInv:   rInv,

		shortMin: shortMin,
		shortMax: shortMax,

		runtime: runtime,
		fns:     fns,
	}, nil
}

func (wc *WitnessCalculator) memFreePos() int32 {
	return int32(binary.LittleEndian.Uint32(wc.runtime.Memory()[:4]))
}

func (wc *WitnessCalculator) setMemFreePos(p int32) {
	binary.LittleEndian.PutUint32(wc.runtime.Memory()[:4], uint32(p))
}

func (wc *WitnessCalculator) allocInt() int32 {
	p := wc.memFreePos()
	wc.setMemFreePos(p + 8)
	return p
}

func (wc *WitnessCalculator) allocFr() int32 {
	p := wc.memFreePos()
	wc.setMemFreePos(p + wc.n32*4 + 8)
	return p
}

func (wc *WitnessCalculator) getInt(p int32) int32 {
	return int32(binary.LittleEndian.Uint32(wc.runtime.Memory()[p : p+4]))
}

func (wc *WitnessCalculator) setInt(p, v int32) {
	binary.LittleEndian.PutUint32(wc.runtime.Memory()[p:p+4], uint32(v))
}

func (wc *WitnessCalculator) setShortPositive(p int32, v *big.Int) {
	if !v.IsInt64() || v.Int64() >= 0x80000000 {
		panic(fmt.Errorf("v should be < 0x80000000"))
	}
	wc.setInt(p, int32(v.Int64()))
	wc.setInt(p+4, 0)
}

func (wc *WitnessCalculator) setShortNegative(p int32, v *big.Int) {
	vNeg := new(big.Int).Set(wc.prime) // prime
	vNeg.Sub(vNeg, wc.shortMax)        // prime - max
	vNeg.Sub(v, vNeg)                  // v - (prime - max)
	vNeg.Add(wc.shortMax, vNeg)        // max + (v - (prime - max))
	if !vNeg.IsInt64() || vNeg.Int64() < 0x80000000 || vNeg.Int64() >= 0x80000000*2 {
		panic(fmt.Errorf("v should be < 0x80000000"))
	}
	wc.setInt(p, int32(vNeg.Int64()))
	wc.setInt(p+4, 0)
}

func (wc *WitnessCalculator) setLongNormal(p int32, v *big.Int) {
	wc.setInt(p, 0)
	wc.setInt(p+4, math.MinInt32) // math.MinInt32 = 0x80000000
	storeBigInt(wc.runtime, p+8, v)
}

func (wc *WitnessCalculator) setFr(p int32, v *big.Int) {
	if v.Cmp(wc.shortMax) == -1 {
		wc.setShortPositive(p, v)
	} else if v.Cmp(wc.shortMin) >= 0 {
		wc.setShortNegative(p, v)
	} else {
		wc.setLongNormal(p, v)
	}
}

func (wc *WitnessCalculator) fromMontgomery(v *big.Int) *big.Int {
	res := new(big.Int).Set(v)
	res.Mul(res, wc.rInv)
	res.Mod(res, wc.prime)
	return res
}

func (wc *WitnessCalculator) getFr(p int32) *big.Int {
	m := wc.runtime.Memory()
	if (m[p+4+3] & 0x80) != 0 {
		res := loadBigInt(wc.runtime, p+4, wc.n32)
		if (m[p+4+3] & 0x40) != 0 {
			return wc.fromMontgomery(res)
		} else {
			return res
		}
	} else {
		if (m[p+3] & 0x40) != 0 {
			res := loadBigInt(wc.runtime, p, 4) // res
			res.Sub(res, wc.shortMax)           // res - max
			res.Add(wc.prime, res)              // res - max + prime
			res.Sub(res, wc.shortMax)           // res - max + (prime - max)
			return res
		} else {
			return loadBigInt(wc.runtime, p, 4)
		}
	}
}

// fnvHash returns (hMSB, hLSB)
func fnvHash(s string) (int32, int32) {
	hash := fnv.New64a()
	hash.Write([]byte(s))
	h := hash.Sum64()
	return int32(h >> 32), int32(h & 0xffffffff)
}

func (wc *WitnessCalculator) CalculateWitness(inputs map[string]interface{}, sanityCheck bool) ([]*big.Int, error) {
	oldMemFreePos := wc.memFreePos()

	if err := wc.doCalculateWitness(inputs, sanityCheck); err != nil {
		return nil, err
	}

	w := make([]*big.Int, wc.nVars)
	for i := int32(0); i < wc.nVars; i++ {
		pWitness, err := wc.fns.getPWitness(i)
		if err != nil {
			return nil, err
		}
		w[i] = wc.getFr(pWitness)
	}

	wc.setMemFreePos(oldMemFreePos)
	return w, nil
}

func (wc *WitnessCalculator) doCalculateWitness(inputs map[string]interface{}, sanityCheck bool) error {
	sanityCheckVal := int32(0)
	if sanityCheck {
		sanityCheckVal = 1
	}
	if err := wc.fns.init(sanityCheckVal); err != nil {
		return err
	}
	pSigOffset := wc.allocInt()
	log.Print("pSigOffset: ", pSigOffset)
	pFr := wc.allocFr()
	log.Print("pFr: ", pFr)

	for inputName, inputValue := range inputs {
		hMSB, hLSB := fnvHash(inputName)
		log.Printf("h(%v) = %v %v", inputName, uint32(hMSB), uint32(hLSB))
		wc.fns.getSignalOffset32(pSigOffset, 0, hMSB, hLSB)
		sigOffset := wc.getInt(pSigOffset)
		log.Print("sigOffset: ", sigOffset)
		fSlice := flatSlice(inputValue)
		for i, value := range fSlice {
			wc.setFr(pFr, value)
			wc.fns.setSignal(0, 0, sigOffset+int32(i), pFr)
		}
	}

	return nil
}

func swap(b []byte) []byte {
	bs := make([]byte, len(b))
	for i := 0; i < len(b); i++ {
		bs[len(b)-1-i] = b[i]
	}
	return bs
}

func parseInputs(inputsJSON []byte) (map[string]interface{}, error) {
	inputsRAW := make(map[string]interface{})
	if err := json.Unmarshal(inputsJSON, &inputsRAW); err != nil {
		return nil, err
	}
	inputs := make(map[string]interface{})
	for inputName, inputValue := range inputsRAW {
		v, err := parseInput(inputValue)
		if err != nil {
			return nil, err
		}
		inputs[inputName] = v
	}
	return inputs, nil
}

func parseInput(v interface{}) (interface{}, error) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		n, ok := new(big.Int).SetString(v.(string), 0)
		if !ok {
			return nil, fmt.Errorf("Error parsing input %v", v)
		}
		return n, nil
	case reflect.Float64:
		return new(big.Int).SetInt64(int64(v.(float64))), nil
	case reflect.Slice:
		res := make([]interface{}, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			var err error
			res[i], err = parseInput(rv.Index(i).Interface())
			if err != nil {
				return nil, fmt.Errorf("Error parsing input %v: %w", v, err)
			}
		}
		return res, nil
	default:
		return nil, fmt.Errorf("Unexpected type for input %v: %T", v, v)
	}
}

func _flatSlice(acc *[]*big.Int, v interface{}) {
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Slice:
		for i := 0; i < rv.Len(); i++ {
			_flatSlice(acc, rv.Index(i).Interface())
		}
	default:
		*acc = append(*acc, v.(*big.Int))
	}
}

func flatSlice(v interface{}) []*big.Int {
	res := make([]*big.Int, 0)
	_flatSlice(&res, v)
	return res
}
