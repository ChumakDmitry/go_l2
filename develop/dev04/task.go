package main

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
)

func unique(words []string) []string {
	uniqueWords := make(map[string]bool, len(words))

	for _, val := range words {
		uniqueWords[val] = true
	}

	words = nil

	for val := range uniqueWords {
		words = append(words, val)
	}
	return words
}

func getLetterStat(str string) map[string]int {
	stat := make(map[string]int, len(str))
	letters := strings.Split(str, "")
	for _, val := range letters {
		stat[val] += 1
	}

	return stat
}

func createAnagramMap(words []string) map[string][]string {
	if len(words) < 1 {
		fmt.Println("Dictionary empty")
		return nil
	}

	result := make(map[string][]string)
	var flag bool

	for _, val := range words {
		//DICTIONARY:
		flag = false
		for word := range result {
			if len(word) == len(val) {
				if reflect.DeepEqual(getLetterStat(strings.ToLower(val)), getLetterStat(word)) {
					result[word] = append(result[word], strings.ToLower(val))
					//goto DICTIONARY
					flag = true
					break
				}
			}
		}
		if !flag {
			result[strings.ToLower(val)] = append(result[strings.ToLower(val)], strings.ToLower(val))
		}
	}

	for key, val := range result {
		uniqueWords := unique(val)

		sort.Strings(uniqueWords)

		if len(uniqueWords) < 2 {
			delete(result, key)
		} else {
			result[key] = uniqueWords
		}
	}
	return result
}

func main() {
	dictionary := []string{
		"пятак",
		"Тяпка",
		"пятка",
		"слиток",
		"листок",
		"столик",
		"арфа",
		"фара",
		"шишка",
		"шишка",
	}

	result := createAnagramMap(dictionary)

	for k, v := range result {
		fmt.Printf("key: %s\nanagarm: %v\n\n", k, v)
	}
}
