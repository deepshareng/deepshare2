package tokenutil

import (
	"math"
	"strings"
)

const (
	// symbols used for short-urls
	symbols = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// someone set us up the bomb !!
	base      = int64(len(symbols))
	floatbase = float64(base)
)

// encodes a number into our *base* representation
func Encode(number int64) string {
	rest := number % base
	// strings are a bit weird in go...
	result := string(symbols[rest])
	if number-rest != 0 {
		newnumber := (number - rest) / base
		result = Encode(newnumber) + result
	}
	return result
}

// Decodes a string given in our encoding and returns the decimal
// integer.
func Decode(input string) int64 {
	l := len(input)
	var sum int = 0
	for index := l - 1; index > -1; index -= 1 {
		current := string(input[index])
		pos := strings.Index(symbols, current)
		sum = sum + (pos * int(math.Pow(floatbase, float64((l-index-1)))))
	}
	return int64(sum)
}
