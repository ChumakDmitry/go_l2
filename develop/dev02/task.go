package main

import (
	"fmt"
)

func checkNumber(num rune) bool {
	if num >= '0' && num <= '9' {
		return true
	}
	return false
}

func checkDoubleNumber(runes []rune) bool {
	result := false
	for i := range runes {
		if checkNumber(runes[i]) {
			if result {
				return true
			}
			result = true
			continue
		}
		result = false
	}
	return false
}

func ParseStr(str string) (string, error) {
	if len(str) < 1 {
		return "", fmt.Errorf("empty string")
	}

	runes := []rune(str)

	if checkNumber(runes[0]) || checkDoubleNumber(runes) {
		return "", fmt.Errorf("uncorrect string")
	}

	result := make([]rune, 0, len(str))

	for i, value := range runes {
		if checkNumber(value) {
			for j := 0; j < int(value-'0')-1; j++ {
				result = append(result, runes[i-1])
			}
			continue
		}
		result = append(result, runes[i])
	}

	return string(result), nil
}

func main() {
	str := ""

	result, err := ParseStr(str)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
