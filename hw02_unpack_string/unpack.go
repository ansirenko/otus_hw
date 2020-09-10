package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString is Unpack error that happened when string is incorrect.
var ErrInvalidString = errors.New("invalid string")

// Unpack is function for unpack string.
func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}

	isDigit := false
	var resp strings.Builder
	var previousLetter rune

	for i := 0; i < len(str); i++ {
		if str[i] == '\\' {
			i++
			if i == len(str) {
				return "", ErrInvalidString
			}
			resp.WriteRune(rune(str[i]))
			previousLetter = rune(str[i])
			isDigit = false
			continue
		}
		if unicode.IsDigit(rune(str[i])) {
			if i == 0 {
				return "", ErrInvalidString
			}
			if isDigit {
				return "", ErrInvalidString
			}
			isDigit = true

			digit, err := strconv.Atoi(string(str[i]))
			if err != nil {
				return "", err
			}
			if digit == 0 {
				buf := resp.String()[:resp.Len()-1]
				resp.Reset()
				resp.WriteString(buf)
				continue
			}
			resp.WriteString(strings.Repeat(string(previousLetter), digit-1))
		} else {
			resp.WriteRune(rune(str[i]))
			previousLetter = rune(str[i])

			isDigit = false
		}
	}

	return resp.String(), nil
}
