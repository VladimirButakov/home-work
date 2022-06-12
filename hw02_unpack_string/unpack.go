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

	for key, value := range runesArray {

		if unicode.IsDigit(value) && key == 0 {
			return "", ErrInvalidString
		}

		if value == '\\' && slash == false {
			slash = true
			continue
		}

		if unicode.IsDigit(value) && (key == len(runesArray)-1) && slash == false {
			continue
		}

		if ((!unicode.IsDigit(value)) && value != '\\') && slash == true {
			return "", ErrInvalidString
		}

		if !(key == len(runesArray)-1) {
			nextValue = runesArray[key+1]

			if unicode.IsDigit(value) && unicode.IsDigit(nextValue) && slash == false {
				return "", ErrInvalidString
			}

			if unicode.IsDigit(value) && slash == false {
				continue
			}

			slash = false
			if unicode.IsDigit(nextValue) {
				countRepeat, _ = strconv.Atoi(string(nextValue))

				if countRepeat == 0 {
					continue
				}

				stringBuilder.WriteString(strings.Repeat(string(value), countRepeat))
				continue
			}
		}

		stringBuilder.WriteString(string(value))
		slash = false
	}

	return stringBuilder.String(), nil
}
