package main

import (
	"os"
	"os/exec"
	"strings"
	"testing"
)

type data struct {
	name  string
	files []string
	flags []string
}

func TestSort(t *testing.T) {
	t.Run("test without flags", func(t *testing.T) {
		cases := []data{
			{
				name:  "file text.txt",
				files: []string{"text.txt"},
			},
			{
				name:  "file text1.txt",
				files: []string{"text1.txt"},
			},
		}

		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				command := []string{"run", "task.go"}
				command = append(command, testCase.files...)
				funcResult, err := exec.Command("go", command...).CombinedOutput()
				if err != nil {
					t.Error(err)
					os.Exit(2)
				}

				realRes, err := exec.Command("sort", testCase.files...).CombinedOutput()
				if err != nil {
					t.Error(err)
					os.Exit(2)
				}

				realResult := strings.ReplaceAll(string(realRes), "\r", "")

				for i, val := range realResult {
					if string(val) != string(funcResult[i]) {
						t.Errorf("Real result %s != %s function result", string(val), string(funcResult[i]))
						os.Exit(1)
					}
				}
			})
		}
	})

	t.Run("test with flag", func(t *testing.T) {
		cases := []data{
			{
				name:  "file text.txt -u",
				files: []string{"text.txt"},
				flags: []string{"-u"},
			},
			{
				name:  "file text1.txt -k3",
				files: []string{"text1.txt"},
				flags: []string{"-k3"},
			},
			{
				name:  "file text.txt -u",
				files: []string{"text.txt"},
				flags: []string{"-u"},
			},
			{
				name:  "file text1.txt -k2 -n",
				files: []string{"text1.txt"},
				flags: []string{"-k=2", "-n"},
			},
		}

		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				command := []string{"run", "task.go"}
				command = append(command, testCase.flags...)
				command = append(command, testCase.files...)
				funcResult, err := exec.Command("go", command...).CombinedOutput()
				if err != nil {
					t.Error(err)
					os.Exit(2)
				}

				command = append(testCase.flags, testCase.files...)
				realRes, err := exec.Command("sort", command...).CombinedOutput()
				if err != nil {
					t.Error(err)
					os.Exit(2)
				}

				realResult := strings.ReplaceAll(string(realRes), "\r", "")

				for i, val := range realResult {
					if string(val) != string(funcResult[i]) {
						t.Errorf("Real result %s != %s function result", string(val), string(funcResult[i]))
						os.Exit(1)
					}
				}
			})
		}
	})
}
