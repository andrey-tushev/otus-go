package hw02unpackstring

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var char rune
	var num int

	src := []rune(str)
	var dst strings.Builder
	for i := 0; i < len(src); {
		if i == len(src)-1 && !unicode.IsDigit(src[i]) {
			char = src[i]
			num = 1
			i++
		} else if i < len(src)-1 && !unicode.IsDigit(src[i]) && !unicode.IsDigit(src[i+1]) {
			char = src[i]
			num = 1
			i++
		} else if i < len(src)-1 && !unicode.IsDigit(src[i]) && unicode.IsDigit(src[i+1]) {
			char = src[i]
			num, _ = strconv.Atoi(string(src[i+1]))
			i += 2
		} else {
			return "", ErrInvalidString
		}

		fmt.Fprint(&dst, strings.Repeat(string(char), num))
	}

	return dst.String(), nil
}
