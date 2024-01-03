package main

import (
	"bufio"
	"errors"
	"fmt"
	ps2 "github.com/mitchellh/go-ps"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

// лёгкий вариант

//func shell(prompt string, in io.Reader, out io.Writer) error {
//	scanner := bufio.NewScanner(in)
//
//	for {
//		fmt.Print(prompt)
//		scanned := scanner.Scan()
//		if !scanned {
//			continue
//		}
//
//		input := scanner.Text()
//
//		if input == "quit" {
//			return nil
//		}
//
//		parts := strings.Split(input, " ")
//
//		if len(parts) < 1 {
//			continue
//		}
//
//		command := strings.TrimSpace(parts[0])
//		args := parts[1:]
//
//		runner := exec.Command(command, args...)
//
//		result, err := runner.CombinedOutput()
//		if err != nil {
//			return err
//		}
//
//		out.Write(result)
//	}
//}

func cd(args []string) string {
	if len(args) > 1 {
		return "too many items\n"
	}

	if len(args) < 1 {
		return ""
	}

	if args[0] == ".." {
		path, _ := os.Getwd()
		tmp := strings.Split(path, `\`)
		path = strings.Join(tmp[:len(tmp)-1], `\`)

		err := os.Chdir(path)
		if err != nil {
			return err.Error()
		}

		return ""
	}

	err := os.Chdir(args[0])
	if err != nil {
		return err.Error()
	}

	return ""
}

func ps() string {
	processes, err := ps2.Processes()
	if err != nil {
		return err.Error()
	}

	output := strings.Builder{}
	output.WriteString("PID\tPPID\tFileName\n---\t----\t--------\n")

	for _, process := range processes {
		output.WriteString(fmt.Sprintf("%d\t%d\t%s\n", process.Pid(),
			process.PPid(), process.Executable()))
	}

	return output.String()
}

func netcat(args []string) string {
	prot := "tcp"

	fmt.Println(args)

	if len(args) < 2 {
		return "you must specify HOST and PORT\n"
	}

	var (
		host string
		port string
	)

	if len(args) > 3 && args[0] == "-u" {
		prot = "udp"
		host = args[1]
		port = args[2]
	} else {
		host = args[0]
		port = args[1]
	}

	address := fmt.Sprintf("%s:%s", host, port)

	conn, err := net.Dial(prot, address)
	if err != nil {
		return err.Error()
	}

	sig := make(chan os.Signal, 1)
	errs := make(chan error, 1)
	signal.Notify(sig, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print(">> ")
			text, err := reader.ReadString('\n')
			if err != nil {
				errs <- err
				return
			}
			fmt.Fprintf(conn, text+"\n")

			if strings.TrimSpace(text) == "STOP" {
				sig <- syscall.SIGQUIT
				return
			}

			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				errs <- err
				return
			}
			fmt.Print("->: " + message)
		}
	}()

	select {
	case err := <-errs:
		return err.Error()
	case s := <-sig:
		return fmt.Sprintf("stoped by signal: %v\n", s)
	}
}

func kill(args []string) string {
	var (
		pid        int
		signalName string
		sig        syscall.Signal
		err        error
	)

	switch len(args) {
	case 1:
		pid, err = strconv.Atoi(args[0])
		if err != nil {
			return "invalid pid. pid must be a number"
		}
		sig = syscall.SIGTERM
	case 3:
		pid, err = strconv.Atoi(args[2])
		if err != nil {
			return "invalid pid. pid must be a number"
		}
		signalName = args[1]

		switch signalName {
		case "SIGINT":
			sig = syscall.SIGINT
		case "SIGTERM":
			sig = syscall.SIGTERM
		case "SIGQUIT":
			sig = syscall.SIGQUIT
		case "SIGKILL":
			sig = syscall.SIGKILL
		case "SIGHUP":
			sig = syscall.SIGHUP
		default:
			return "invalid input sigspec\n"
		}
	default:
		return "invalid input. ex.: kill [-s sigspec] pid\n"
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return err.Error()
	}

	err = process.Signal(sig)
	if err != nil {
		return err.Error()
	}

	return fmt.Sprintf("kill: %d was %s\n", pid, sig.String())
}

func runCommand(input string) string {
	if len(input) < 1 {
		return ""
	}

	command := strings.Fields(input)

	args := command[1:]

	var output string

	switch command[0] {
	case "cd":
		output = cd(args)
	case "pwd":
		output, _ = os.Getwd()
		output = "\nPath\n" + "----\n" + output + "\n"
	case "echo":
		output = strings.Join(args, " ")
		output += "\n"
	case "kill":
		output = kill(args)
	case "ps":
		output = ps()
	case "netcat":
		output = netcat(args)
	default:
		runner := exec.Command(command[0], args...)

		result, err := runner.CombinedOutput()
		if err != nil {
			return err.Error()
		}

		output = string(result)
	}

	return output
}

func readInput() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	fmt.Fprintf(os.Stdout, "My Shell: %s>", currentDir)
	scanner := bufio.NewScanner(os.Stdin)
	scanned := scanner.Scan()
	if !scanned {
		return "", nil
	}

	input := scanner.Text()
	if input == "" {
		return "", errors.New("empty string")
	}

	return input, nil
}

func shell() error {
	for {
		input, err := readInput()
		if err != nil {
			if err.Error() == "empty string" {
				continue
			}
			return err
		}

		if input == "quit" {
			break
		}

		output := runCommand(input)

		fmt.Println(output)
	}

	return nil
}

func main() {
	err := shell()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
