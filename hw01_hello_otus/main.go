package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	src := "Hello, OTUS!"
	rev := stringutil.Reverse(src)
	fmt.Println(rev)
}
