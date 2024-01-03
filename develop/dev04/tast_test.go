package main

import (
	"strings"
	"testing"
)

func TestCreateAnagramMap(t *testing.T) {
	dictionary := []string{
		"Пятак",
		"Пятак",
		"пятка",
		"Тяпка",
		"слиток",
		"слиток",
		"столик",
		"листок",
		"Топот",
		"Потоп",
		"силач",
		"числа",
		"Каприз",
		"приказ",
		"шишка",
	}

	result := map[string][]string{
		"пятак":  {"пятак", "пятка", "тяпка"},
		"слиток": {"листок", "слиток", "столик"},
		"силач":  {"силач", "числа"},
		"каприз": {"каприз", "приказ"},
	}

	funcResult := createAnagramMap(dictionary)

	for key, val := range funcResult {
		anagrams, ok := result[key]
		if !ok {
			t.Errorf("excess key: %s", key)
		}
		joinedResultAnagrams := strings.Join(anagrams, " ")
		joinedFuncAnagrams := strings.Join(val, " ")
		if joinedFuncAnagrams != joinedResultAnagrams {
			t.Errorf("wrong anagrams:\nShould: %s\nGot: %s\n", joinedResultAnagrams, joinedFuncAnagrams)
		}
	}
}
