package main

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFlatSlice(t *testing.T) {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	three := new(big.Int).SetInt64(3)
	four := new(big.Int).SetInt64(4)

	a := one
	fa := flatSlice(a)
	require.Equal(t, []*big.Int{one}, fa)

	b := []*big.Int{one, two}
	fb := flatSlice(b)
	require.Equal(t, []*big.Int{one, two}, fb)

	c := []interface{}{one, []*big.Int{two, three}}
	fc := flatSlice(c)
	require.Equal(t, []*big.Int{one, two, three}, fc)

	d := []interface{}{[]*big.Int{one, two}, []*big.Int{three, four}}
	fd := flatSlice(d)
	require.Equal(t, []*big.Int{one, two, three, four}, fd)
}

func TestParseInputs(t *testing.T) {
	one := new(big.Int).SetInt64(1)
	two := new(big.Int).SetInt64(2)
	three := new(big.Int).SetInt64(3)
	four := new(big.Int).SetInt64(4)

	a, err := parseInputs([]byte(`{"a": 1, "b": "2"}`))
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{"a": one, "b": two}, a)

	b, err := parseInputs([]byte(`{"a": 1, "b": [2, 3]}`))
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{"a": one, "b": []interface{}{two, three}}, b)

	c, err := parseInputs([]byte(`{"a": 1, "b": [[1, 2], [3, 4]]}`))
	require.Nil(t, err)
	require.Equal(t, map[string]interface{}{"a": one, "b": []interface{}{[]interface{}{one, two}, []interface{}{three, four}}}, c)
}
