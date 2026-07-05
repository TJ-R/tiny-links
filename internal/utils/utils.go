package utils

import (
	"strings"
)

func Base62_encoding(id int, base62_map map[int8]rune) string {
	stringBuilder := strings.Builder{}
	quotient := id
	remainder := -1
	for {
		remainder = quotient % 62
		quotient /= 62

		// if remainder is 0 we are done
		if remainder == 0 {
			break
		} else {
			stringBuilder.WriteRune(base62_map[int8(remainder)])
		}
	}

	return reverseString(stringBuilder.String())
}

func Make_base62_map() map[int8]rune {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	m := make(map[int8]rune)

	for i := range int8(len(alphabet)) {
		m[i] = rune(alphabet[i])
	}

	return m
}

func reverseString(s string) string {
	stringBuilder := strings.Builder{}

	for i := len(s) - 1; i >= 0; i-- {
		stringBuilder.WriteByte(s[i])
	}

	return stringBuilder.String()
}
