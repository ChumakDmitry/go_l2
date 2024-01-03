package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Args struct {
	k     int
	n     bool
	nk    int
	r     bool
	u     bool
	files []string
}

func splitByWords(lines []string) [][]string {
	words := make([][]string, len(lines))
	for ind := range lines {
		words[ind] = strings.Fields(lines[ind])
	}

	return words
}

func checkNumber(str []string, column int) int {
	res, err := strconv.Atoi(str[column])
	if err != nil {
		fmt.Println(err)
	}
	return res
}

func checkCondition(lines []string, condition func(str []string) bool) []string {
	words := splitByWords(lines)
	min := 0

	for i := range words {
		if condition(words[i]) {
			lines[i], lines[min] = lines[min], lines[i]
			min++
		}
	}
	fmt.Println(lines[:min])
	return lines[:min]
}

func columnSort(lines []string, column int) {
	index := column - 1
	temp := checkCondition(lines, func(str []string) bool {
		return len(str) >= index
	})

	sort.Slice(temp, func(i, j int) bool {
		temp := splitByWords(temp)
		return temp[i][index] < temp[j][index]
	})
}

func numberSort(lines []string, column int) {
	index := column - 1

	temp := checkCondition(lines, func(str []string) bool {
		return len(str) >= index && checkNumber(str, index) != 0
	})

	sort.Slice(temp, func(i, j int) bool {
		temp := splitByWords(temp)
		return checkNumber(temp[i], index) < checkNumber(temp[j], index)
	})
}

func reverse(words []string) []string {
	left, right := 0, len(words)-1
	for left < right {
		words[left], words[right] = words[right], words[left]
		left++
		right--
	}

	return words
}

func unique(words []string) []string {
	uniqueWords := make(map[string]bool, len(words))

	for _, val := range words {
		uniqueWords[val] = true
	}

	result := make([]string, 0, len(uniqueWords))

	for val := range uniqueWords {
		result = append(result, val)
	}

	return result
}

func sortStr(args *Args, words []string) []string {
	if args.u {
		words = unique(words)
	}

	switch {
	case args.k > 0 && !args.n:
		columnSort(words, args.k)
		break
	case args.n:
		if args.k < 1 {
			fmt.Println("missing column number")
			return nil
		}
		numberSort(words, args.k)
		break
	default:
		sort.Strings(words)
	}

	if args.r {
		words = reverse(words)
	}

	return words
}

func getFlags() (*Args, error) {
	k := flag.Int("k", 0, "define sort column")
	n := flag.Bool("n", false, "sort by number value")
	r := flag.Bool("r", false, "reverse sort")
	u := flag.Bool("u", false, "only unique string")

	flag.Parse()

	args := &Args{
		k: *k,
		n: *n,
		r: *r,
		u: *u,
	}

	if args.k < 0 {
		return nil, errors.New("counter can't be negative")
	}

	args.files = append(args.files, flag.Args()...)

	return args, nil
}

func getStr(args *Args) ([]string, error) {
	var str []string

	if len(args.files) < 1 {
		return nil, errors.New("missing file name")
	}

	for _, val := range args.files {
		file, err := os.Open(val)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			str = append(str, scanner.Text())
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}
	}

	return str, nil
}

func Sort() (string, error) {
	args, err := getFlags()
	if err != nil {
		return "", err
	}

	str, err := getStr(args)
	if err != nil {
		return "", err
	}

	sortedStr := sortStr(args, str)
	return strings.Join(sortedStr, "\n"), nil
}

func main() {
	sorted, err := Sort()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(sorted)
}
