package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Args struct {
	f []string
	d string
	s bool
}

func getFlags() (*Args, error) {
	f := flag.String("f", "", "Display specific fields")
	d := flag.String("d", "\t", "Use specific delimiter (default TAB)")
	s := flag.Bool("s", false, "Only-delimited string output")

	flag.Parse()

	if len(*f) < 1 {
		return nil, errors.New("invalid input f, e.g.: 1, 2, or 1-8")
	}

	fields := strings.Split(*f, ",")

	args := &Args{
		f: fields,
		d: *d,
		s: *s,
	}

	return args, nil
}

func getStr() (string, error) {
	var str strings.Builder
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		str.WriteString(scanner.Text())
	}

	result := str.String()
	fmt.Println(result)

	if result == "" {
		return "", errors.New("empty string")
	}

	return result, nil
}

func cut(inputStr string, args *Args) (string, error) {
	var result strings.Builder
	fields := make(map[int]bool)

	delimiter := "\t"
	if args.d != "" {
		if len(args.d) > 1 {
			return "", errors.New("delimiter must be a single character")
		}
		delimiter = args.d
	}

	if len(args.f) > 0 {
		for _, field := range args.f {
			fieldRange := strings.Split(field, "-")
			if len(fieldRange) == 2 {
				startRange, err := strconv.Atoi(fieldRange[0])
				if err != nil {
					return "", errors.New(fmt.Sprintf("invalid field value: %v", fieldRange[0]))
				}

				endRange, err := strconv.Atoi(fieldRange[1])
				if err != nil {
					return "", errors.New(fmt.Sprintf("invalid field value: %v", fieldRange[1]))
				}

				if endRange < startRange {
					return "", errors.New("incorrect range selected")
				}

				if startRange < 1 {
					return "", errors.New("the beginning cannot be less than 1")
				}

				for i := startRange; i <= endRange; i++ {
					fields[i] = true
				}
			} else {
				numField, err := strconv.Atoi(field)
				if err != nil {
					return "", errors.New(fmt.Sprintf("invalid field value: %v", field))
				}

				if numField < 1 {
					return "", errors.New("the beginning cannot be less than 1")
				}

				fields[numField] = true
			}
		}
	} else {
		return "", errors.New("you must specify a list fields")
	}

	str := strings.Split(inputStr, delimiter)

	if len(str) == 1 {
		return "", nil
	}

	needDelimiter := false

	for i, part := range str {
		_, ok := fields[i+1]
		if ok {
			if needDelimiter {
				result.WriteString(delimiter + part)
			} else {
				result.WriteString(part)
				needDelimiter = true
			}
		}
	}

	return result.String(), nil
}

func Start() (string, error) {
	args, err := getFlags()
	if err != nil {
		return "", err
	}

	inputStr, err := getStr()
	if err != nil {
		return "", err
	}

	result, err := cut(inputStr, args)
	if err != nil {
		return "", err
	}

	return result, nil
}

func main() {
	result, err := Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println(result)
}
