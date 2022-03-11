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
	var char rune // Символ
	var num int   // Количество повторов

	src := []rune(str)
	var dst strings.Builder

	for i := 0; i < len(src); {
		// Берем символ
		char = src[i]
		i++

		// Если вместо символа оказалась цифра, то значит нарушен формат строки
		if unicode.IsDigit(char) {
			return "", ErrInvalidString
		}

		// Если это было начала экранирования
		// тогда возьмем заэкранированный символ
		if char == '\\' {
			// Нельзя чтобы символ экранирования был последним
			if i >= len(src) {
				return "", ErrInvalidString
			}

			char = src[i]
			i++
		}

		// А нет ли там цифрового суффикса?
		if i < len(src) && unicode.IsDigit(src[i]) {
			num, _ = strconv.Atoi(string(src[i]))
			i++
		} else {
			num = 1
		}

		// Добавляем символы к результату
		fmt.Fprint(&dst, strings.Repeat(string(char), num))
	}

	return dst.String(), nil
}
