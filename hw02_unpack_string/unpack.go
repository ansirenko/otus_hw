package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

// ErrInvalidString is Unpack error that happened when string is incorrect.
var ErrInvalidString = errors.New("invalid string")

func checkForCorrect(str string) error {
	isDigit := false
	isNumber := false
	for i := 0; i < len(str); i++ {
		if str[i] == '\\' {
			isNumber = false
			i++
			continue
		}
		isDigit = unicode.IsDigit(rune(str[i]))
		if isDigit && isNumber {
			return ErrInvalidString
		}
		isNumber = isDigit
		if isDigit && i == 0 {
			return ErrInvalidString
		}
	}
	return nil
}

func isRepeat(val rune) (count int, err error) {
	if !unicode.IsDigit(val) {
		return -1, nil
	}
	digit, err := strconv.Atoi(string(val))
	if err != nil {
		return 0, err
	}
	return digit, nil
}

// Unpack is function for unpack string.
func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	if err := checkForCorrect(str); err != nil {
		return "", err
	}
	var resp strings.Builder
	input := []rune(str)
	for i := 0; i < len(input); i++ {
		if str[i] == '\\' {
			i++
			resp.WriteRune(input[i])
			continue
		}
		count, err := isRepeat(input[i])
		if err != nil {
			return "", err
		}
		switch count {
		case 0:
			{
				buf := resp.String()[:resp.Len()-1]
				resp.Reset()
				resp.WriteString(buf)
				continue
			}
		case -1:
			{
				resp.WriteRune(input[i])
			}
		default:
			resp.WriteString(strings.Repeat(string(input[i-1]), count-1))
		}
	}
	return resp.String(), nil
}
