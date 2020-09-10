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

	var response strings.Builder

	isSafe := false
	isNumber := false
	for index, val := range str {
		if val == '\\' && !isSafe {
			isSafe = true
			isNumber = false
			continue
		}
		if unicode.IsDigit(val) && !isSafe {
			if index == 0 {
				return "", ErrInvalidString
			}
			if isNumber {
				return "", ErrInvalidString
			}
			isNumber = true

			digit, err := strconv.Atoi(string(val))
			if err != nil {
				return "", err
			}
			if digit == 0 {
				buf := response.String()[:response.Len()-1]
				response.Reset()
				response.WriteString(buf)
				continue
			}
			response.WriteString(strings.Repeat(string(str[index-1]), digit-1))
		} else {
			response.WriteRune(val)
			isNumber = false
			isSafe = false
		}
	}

	return response.String(), nil
}
