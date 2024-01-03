package main

//Реализовать утилиту фильтрации по аналогии с консольной утилитой (man grep — смотрим описание и основные параметры).
//
//
//Реализовать поддержку утилитой следующих ключей:
//-A - "after" печатать +N строк после совпадения
//-B - "before" печатать +N строк до совпадения
//-C - "context" (A+B) печатать ±N строк вокруг совпадения
//-c - "count" (количество строк)
//-i - "ignore-case" (игнорировать регистр)
//-v - "invert" (вместо совпадения, исключать)
//-F - "fixed", точное совпадение со строкой, не паттерн
//-n - "line num", напечатать номер строки

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Args struct {
	A int
	B int
	C int
	c bool
	i bool
	v bool
	F bool
	n bool

	pattern string
	files   []string
}

func getFlags() (*Args, error) {
	A := flag.Int("A", 0, "print +N string after matches")
	B := flag.Int("B", 0, "print +N string before matches")
	C := flag.Int("C", 0, "print +-N string around matches")
	c := flag.Bool("c", false, "count string")
	i := flag.Bool("i", false, "ignore case")
	v := flag.Bool("v", false, "invert (exclude)")
	F := flag.Bool("F", false, "exact match to the string")
	n := flag.Bool("n", false, "print line num")

	flag.Parse()

	args := &Args{
		A: *A,
		B: *B,
		C: *C,
		c: *c,
		i: *i,
		v: *v,
		F: *F,
		n: *n,
	}

	if args.A < 0 || args.B < 0 || args.C < 0 {
		return nil, errors.New("counter can't be negative")
	}

	if args.A > 0 && args.C > 0 || args.B > 0 && args.C > 0 {
		return nil, errors.New("too many params, use one of these: A or (and) B or only C")
	}

	if len(flag.Args()) == 0 {
		return nil, errors.New("missing pattern and file name")
	}

	pattern := flag.Args()[0]

	if args.i {
		args.pattern = strings.ToLower(pattern)
	} else {
		args.pattern = pattern
	}

	args.files = append(args.files, flag.Args()[1:]...)

	return args, nil
}

func findContainsIndex(fileStr map[string][]string, args *Args) map[string][]bool {
	index := make(map[string][]bool, len(fileStr))

	for k, file := range fileStr {
		for _, str := range file {
			if args.i {
				strings.ToLower(str)
			}

			if args.F {
				if reflect.DeepEqual(str, args.pattern) {
					index[k] = append(index[k], true)
					continue
				}
			} else if strings.Contains(str, args.pattern) {
				index[k] = append(index[k], true)
				continue
			}

			index[k] = append(index[k], false)
		}
	}

	return index
}

func addElemNumber(fileStr map[string][]string) map[string][]string {
	for _, file := range fileStr {
		for i, str := range file {
			numberStr := strings.Builder{}
			numberStr.Grow(len(str) + 3)

			number := fmt.Sprintf("%d: ", i)
			numberStr.WriteString(number)
			numberStr.WriteString(str)

			file[i] = numberStr.String()
		}
	}

	return fileStr
}

func addElemAfter(index map[string][]bool, offset int) map[string][]bool {
	newIndex := make(map[string][]bool, len(index))
	for i, str := range index {
		newIndex[i] = str
	}

	for filename, data := range index {
		for i, strBool := range data {
			if strBool {
				for j := i; j <= i+offset && j < len(data); j++ {
					newIndex[filename][j] = true
				}
			}
		}
	}

	return newIndex
}

func addElemBefore(index map[string][]bool, offset int) map[string][]bool {
	newIndex := make(map[string][]bool, len(index))
	for i, str := range index {
		newIndex[i] = str
	}

	for _, file := range newIndex {
		for i, strBool := range file {
			if strBool {
				for j := i; j >= i-offset && j > 0; j-- {
					file[j] = true
				}
			}
		}
	}

	return newIndex
}

func invertIndex(index map[string][]bool) map[string][]bool {
	newIndex := make(map[string][]bool, len(index))
	for key, str := range index {
		newSlice := make([]bool, 0, len(str))
		for _, val := range str {
			newSlice = append(newSlice, !val)
		}

		newIndex[key] = newSlice
	}

	return newIndex
}

func resultContains(fileStr map[string][]string, index map[string][]bool) map[string][]string {
	result := make(map[string][]string, len(fileStr))

	for filename, data := range fileStr {
		resultSlice := make([]string, 0, len(data))
		for i, str := range data {
			if index[filename][i] {
				resultSlice = append(resultSlice, str)
			}
		}

		if len(resultSlice) == 0 {
			continue
		}
		result[filename] = resultSlice
	}

	return result
}

func getCountContains(index map[string][]bool) map[string][]string {
	count := 0

	for _, data := range index {
		for _, val := range data {
			if val {
				count++
			}
		}
	}

	result := make(map[string][]string)
	result["count"] = append(result["count"], strconv.Itoa(count))
	return result
}

func findContainsStr(fileStr map[string][]string, index map[string][]bool, args *Args) map[string][]string {
	if args.c {
		result := getCountContains(index)
		return result
	}

	if args.n {
		fileStr = addElemNumber(fileStr)
	}

	if args.C > 0 {
		args.A = args.C
		args.B = args.C
	}

	if args.A > 0 {
		index = addElemAfter(index, args.A)
	}

	if args.B > 0 {
		index = addElemBefore(index, args.B)
	}

	if args.v {
		index = invertIndex(index)
	}

	result := resultContains(fileStr, index)

	return result
}

func grepStr(args *Args, fileStr map[string][]string) map[string][]string {
	index := findContainsIndex(fileStr, args)
	result := findContainsStr(fileStr, index, args)

	return result
}

func getStr(args *Args) (map[string][]string, error) {
	fileStr := make(map[string][]string)

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
			fileStr[val] = append(fileStr[val], scanner.Text())
		}

		err = scanner.Err()
		if err != nil {
			return nil, err
		}
	}

	return fileStr, nil
}

func grep() (map[string][]string, error) {
	args, err := getFlags()
	if err != nil {
		return nil, err
	}

	fileStr, err := getStr(args)
	if err != nil {
		return nil, err
	}

	grepedStr := grepStr(args, fileStr)
	return grepedStr, nil
}

func main() {
	greped, err := grep()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if len(greped) == 0 {
		fmt.Println("no mismatch")
	}

	for file, val := range greped {
		if file == "count" {
			fmt.Printf("count = %s", val[0])
			break
		}

		fmt.Printf("File: %s\n", file)
		for _, str := range val {
			fmt.Printf("\tstr: %s\n", str)
		}
	}
}
