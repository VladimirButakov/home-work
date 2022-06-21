package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	var stringBuilder strings.Builder
	runesArray := []rune(str)
	var countRepeat int
	var nextValue rune
	var slash bool
	var lastElem bool

	for key, value := range runesArray {
		lastElem = key == len(runesArray)-1

		if (!unicode.IsDigit(value) && value != '\\' && slash) ||
			(unicode.IsDigit(value) && key == 0) {
			return "", ErrInvalidString
		}

		if value == '\\' && !slash {
			slash = true
			continue
		}

		if unicode.IsDigit(value) && lastElem && !slash {
			continue
		}

		if !lastElem {
			nextValue = runesArray[key+1]

			if unicode.IsDigit(value) && unicode.IsDigit(nextValue) && !slash {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(value) && !slash {
				continue
			}

			slash = false
			if unicode.IsDigit(nextValue) {
				countRepeat, _ = strconv.Atoi(string(nextValue))
				stringBuilder.WriteString(strings.Repeat(string(value), countRepeat))
				continue
			}
		}

		stringBuilder.WriteString(string(value))
		slash = false
	}

	return stringBuilder.String(), nil
}
