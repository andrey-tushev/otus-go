package main

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	env := Environment{
		"FOO": {
			Value:      "foo",
			NeedRemove: false,
		},
		"BAR": {
			Value:      "",
			NeedRemove: false,
		},
		"USER": {
			NeedRemove: true,
		},
	}

	require.Equal(t, "foo", checkVar("FOO", env))
	require.Equal(t, "", checkVar("BAR", env))
	require.Equal(t, "", checkVar("USER", env))
}

func TestRetCode(t *testing.T) {
	env := Environment{}

	code := RunCmd("ls", []string{"."}, env)
	require.Equal(t, 0, code)

	code = RunCmd("ls", []string{"/bad-dir"}, env)
	require.NotEqual(t, 0, code)
}

func checkVar(name string, env Environment) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	args := []string{`-c`, "echo $" + name}
	RunCmd("bash", args, env)

	w.Close()
	capture, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return strings.TrimRight(string(capture), "\n")
}
