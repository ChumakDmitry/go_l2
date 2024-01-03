package main

import (
	"testing"
)

type Data struct {
	name   string
	input  string
	args   Args
	result string
}

func TestCut(t *testing.T) {
	data := []Data{
		{
			name:  "#1 Out first column",
			input: "Hello world and another",
			args: Args{
				f: []string{"1-4"},
				d: " ",
				s: false,
			},
			result: "Hello world and another",
		},
		{
			name:  "#2 change delimiter",
			input: "Hello world; i say; hello",
			args: Args{
				f: []string{"1"},
				d: ";",
				s: false,
			},
			result: "Hello world",
		},
		{
			name:  "#3 check only-delimiter string",
			input: "HelloWorld\nHello-world\nHello world",
			args: Args{
				f: []string{"1"},
				d: " ",
				s: true,
			},
			result: "HelloWorld\nHello-world\nHello",
		},
	}

	for _, testCase := range data {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := cut(testCase.input, &testCase.args)
			if err != nil {
				t.Error(err)
			}
			if result != testCase.result {
				t.Errorf("result function != real result\n %s != %s", result, testCase.result)
			}
		})
	}
}
