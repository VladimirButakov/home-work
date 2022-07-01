package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type element struct {
	Key   string
	Value int
}

func Top10(str string) []string {
	arrayStrings := strings.Fields(str)
	sortMap := make(map[string]int)
	result := []string{}
	sortSlice := []element{}

	for _, value := range arrayStrings {
		sortMap[value]++
	}

	for key, value := range sortMap {
		sortSlice = append(sortSlice, element{key, value})
	}

	sort.Slice(sortSlice, func(i, j int) bool {
		iWord, jWord := sortSlice[i], sortSlice[j]
		if iWord.Value == jWord.Value {
			return iWord.Key < jWord.Key
		}

		return iWord.Value > jWord.Value
	})

	for key, value := range sortSlice {
		if key == 10 {
			break
		}

		result = append(result, value.Key)
	}

	return result
}
