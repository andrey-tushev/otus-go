package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	env, err := ReadDir("testdata/env")
	require.NoError(t, err)

	require.Equal(t, `"hello"`, env["HELLO"].Value)
	require.False(t, env["HELLO"].NeedRemove)

	require.Equal(t, "bar", env["BAR"].Value)
	require.False(t, env["BAR"].NeedRemove)

	require.Equal(t, "   foo\nwith new line", env["FOO"].Value)
	require.False(t, env["FOO"].NeedRemove)

	require.True(t, env["UNSET"].NeedRemove)

	require.Equal(t, "", env["EMPTY"].Value)
	require.False(t, env["EMPTY"].NeedRemove)
}

func TestTrim(t *testing.T) {
	require.Equal(t, "HELLO WORLD", trim("HELLO WORLD"))

	require.Equal(t, "  HELLO", trim("  HELLO  "))
	require.Equal(t, "  HELLO", trim("  HELLO \t "))

	require.Equal(t, "HELLO\nWORLD", trim("HELLO\x00WORLD"))
}
