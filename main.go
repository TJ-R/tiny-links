package main

import (
	"fmt"
	"strings"
)

func main() {
	// Create server and listen for requests
	base62_map := make_base62_map()
	fmt.Println(base62_encoding(1000, base62_map))
}

func base62_encoding(id int, base62_map map[int8]rune) string {
	stringBuilder := strings.Builder{}
	quotient := id
	remainder := -1
	for {
		remainder = quotient % 62
		quotient /= 62

		stringBuilder.WriteRune(base62_map[int8(remainder)])
		// if remainder is 0 we are done
		if remainder == 0 {
			break
		}
	}

	return stringBuilder.String()
}

func make_base62_map() map[int8]rune {
	const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	m := make(map[int8]rune)

	for i := range int8(len(alphabet)) {
		m[i] = rune(alphabet[i])
	}

	return m
}
