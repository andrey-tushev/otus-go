package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const (
		inFile  = "testdata/input.txt"
		outFile = "out.txt"
	)

	in, _ := ioutil.ReadFile(inFile)

	err := Copy(inFile, outFile, 0, 0)
	require.NoError(t, err)
	out, _ := ioutil.ReadFile(outFile)
	require.Equal(t, in, out)

	err = Copy(inFile, outFile, 1000, 0)
	require.NoError(t, err)
	out, _ = ioutil.ReadFile(outFile)
	require.Equal(t, in[1000:], out)

	err = Copy(inFile, outFile, 1000, 100000)
	require.NoError(t, err)
	out, _ = ioutil.ReadFile(outFile)
	require.Equal(t, in[1000:], out)

	err = Copy(inFile, outFile, 1000, 500)
	require.NoError(t, err)
	out, _ = ioutil.ReadFile(outFile)
	require.Equal(t, in[1000:1500], out)

	err = Copy(inFile, outFile, 1000000, 500)
	require.Error(t, err)

	os.Remove(outFile)
}
